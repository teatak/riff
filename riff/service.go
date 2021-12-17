package riff

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/teatak/riff/api"
	"github.com/teatak/riff/common"
	"github.com/teatak/riff/git"
)

const branch = "branch"
const release = "release"
const tag = "tag"
const content = "content"
const latest = "latest"

func getVersionType(version string) (t string, v string) {
	arr := strings.Split(version, ":")
	if len(arr) > 1 {
		t = arr[0]
		v = strings.Join(arr[1:], ":")
	}
	return t, v
}

type Services map[string]*Service

type Service struct {
	Version       uint64
	State         api.StateType //Current state
	StateChange   time.Time     //Time last state change happened
	Progress      *Progress     //update percent
	Config        string        //config file
	StartTime     time.Time     //start time
	StatusContent string        //status content
	*ServiceConfig
}

type Progress struct {
	Current    int32
	Total      int32
	InProgress bool
}

type ServiceConfig struct {
	Name       string   `yaml:"name,omitempty"`
	Ip         string   `yaml:"ip,omitempty"`
	Port       int      `yaml:"port,omitempty"`
	Env        []string `yaml:"env,omitempty"`
	Dir        string   `yaml:"dir,omitempty"`
	Command    []string `yaml:"command,omitempty"`
	StatusPage string   `yaml:"status_page,omitempty"`
	PidFile    string   `yaml:"pid_file,omitempty"`
	StdOutFile string   `yaml:"std_out_file,omitempty"`
	StdErrFile string   `yaml:"std_err_file,omitempty"`
	Grace      bool     `yaml:"grace,omitempty"`
	RunAtLoad  bool     `yaml:"run_at_load,omitempty"`
	KeepAlive  bool     `yaml:"keep_alive,omitempty"`
	*Deploy    `yaml:"deploy,omitempty"`
}

type Deploy struct {
	Provider    string `yaml:"provider,omitempty"`
	Token       string `yaml:"token,omitempty"`
	Repository  string `yaml:"repository,omitempty"`
	Version     string `yaml:"version,omitempty"`
	Payload     string `yaml:"payload,omitempty"`
	ServicePath string `yaml:"service_path,omitempty"`
}

func (s *Server) initServices() {
	s.Self.LoadServices()
	s.Shutter()
}

func (s *Service) rewriteConfig() {
	ipNet, _, _ := net.ParseCIDR(server.config.IP)
	//get this server ip
	ip := ipNet.String()
	name := s.Name
	port := s.Port
	s.Ip = ip

	replaceValue := func(value string) string {
		value = strings.ReplaceAll(value, "${name}", name)
		value = strings.ReplaceAll(value, "${ip}", ip)
		value = strings.ReplaceAll(value, "${port}", strconv.Itoa(port))
		return value
	}

	replaceValues := func(values []string) []string {
		v := []string{}
		for _, item := range values {
			item = replaceValue(item)
			v = append(v, item)
		}
		return v
	}

	s.Env = replaceValues(s.Env)
	s.Command = replaceValues(s.Command)
	s.StatusPage = replaceValue(s.StatusPage)
	s.PidFile = replaceValue(s.PidFile)

	if s.Deploy != nil {
		s.Deploy.ServicePath = replaceValue(s.Deploy.ServicePath)
		s.Deploy.Repository = replaceValue(s.Deploy.Repository)
	}
}
func (s *Server) handleServices() {
	go func() {
		for {
			select {
			case <-s.ShutdownCh:
				return
			default:
			}
			for _, service := range s.Self.Services {
				//first run it
				service.checkState()
			}
			preSnap := s.Self.SnapShot
			s.Self.Shutter()
			nowSnap := s.Self.SnapShot
			if preSnap != nowSnap { //if presnap != nowsnap then add version and create server snapshort
				s.Self.VersionInc()
				s.Self.Shutter()
				s.Shutter()
				s.watch.Dispatch(WatchParam{
					Name:      s.Self.Name,
					WatchType: NodeChanged,
				})
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			select {
			case <-s.ShutdownCh:
				return
			default:
			}
			for _, service := range s.Self.Services {
				//first run it
				service.keepAlive()
				service.update()
			}
			preSnap := s.Self.SnapShot
			s.Self.Shutter()
			nowSnap := s.Self.SnapShot
			if preSnap != nowSnap {
				s.Shutter()
				s.watch.Dispatch(WatchParam{
					Name:      s.Self.Name,
					WatchType: NodeChanged,
				})
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func (s *Service) checkState() {
	if s.Command != nil && len(s.Command) > 0 {
		//if have command onlycheck statuspage
		if s.StatusPage != "" {
			status := 0
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			req, err := http.NewRequest("GET", s.StatusPage, nil)
			if err == nil {
				res, err := http.DefaultClient.Do(req.WithContext(ctx))
				if err == nil {
					status = res.StatusCode
					if status == 200 {
						body, _ := ioutil.ReadAll(res.Body)
						defer res.Body.Close()
						s.StatusContent = string(body)
					}
				}
			}
		}
		if pid := s.GetPid(); pid == 0 {
			if s.State != api.StateDead {
				s.State = api.StateDead
				server.watch.Dispatch(WatchParam{
					Name:      s.Name,
					WatchType: ServiceChanged,
				})
			}
		} else {
			if s.State != api.StateAlive {
				s.State = api.StateAlive
				server.watch.Dispatch(WatchParam{
					Name:      s.Name,
					WatchType: ServiceChanged,
				})
			}
		}
	} else {
		//only use statuspage to check status
		if s.StatusPage != "" {
			status := 0
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			req, err := http.NewRequest("GET", s.StatusPage, nil)
			if err == nil {
				res, err := http.DefaultClient.Do(req.WithContext(ctx))
				if err == nil {
					status = res.StatusCode
					if status == 200 {
						body, _ := ioutil.ReadAll(res.Body)
						defer res.Body.Close()
						s.StatusContent = string(body)
					}
				}
				if status == 200 {
					if s.State != api.StateAlive {
						server.watch.Dispatch(WatchParam{
							Name:      s.Name,
							WatchType: ServiceChanged,
						})
					}
					s.State = api.StateAlive
				} else {
					if s.State != api.StateDead {
						server.watch.Dispatch(WatchParam{
							Name:      s.Name,
							WatchType: ServiceChanged,
						})
					}
					s.State = api.StateDead
				}
			}
		}
	}
	//if s.StatusPage != "" {
	//	status := 0
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	//	defer cancel()
	//	req, _ := http.NewRequest("GET", s.StatusPage, nil)
	//	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	//	if err == nil {
	//		status = res.StatusCode
	//		if status == 200 {
	//			body, _ := ioutil.ReadAll(res.Body)
	//			defer res.Body.Close()
	//			s.StatusContent = string(body)
	//		}
	//	}
	//	if status == 200 {
	//		if s.State != api.StateAlive {
	//			server.watch.Dispatch(WatchParam{
	//				Name:      s.Name,
	//				WatchType: ServiceChanged,
	//			})
	//		}
	//		s.State = api.StateAlive
	//	} else {
	//		if s.State != api.StateDead {
	//			server.watch.Dispatch(WatchParam{
	//				Name:      s.Name,
	//				WatchType: ServiceChanged,
	//			})
	//		}
	//		s.State = api.StateDead
	//	}
	//} else {
	//	if pid := s.GetPid(); pid == 0 {
	//		if s.State != api.StateDead {
	//			server.watch.Dispatch(WatchParam{
	//				Name:      s.Name,
	//				WatchType: ServiceChanged,
	//			})
	//		}
	//		s.State = api.StateDead
	//	} else {
	//		if s.State != api.StateAlive {
	//			server.watch.Dispatch(WatchParam{
	//				Name:      s.Name,
	//				WatchType: ServiceChanged,
	//			})
	//		}
	//		s.State = api.StateAlive
	//	}
	//}
}

func (s *Service) keepAlive() {
	if pid := s.GetPid(); pid == 0 && s.KeepAlive {
		err := s.Start()
		if err != nil {
			server.Logger.Printf(errorServicePrefix+"%s running error: %v", s.Name, err)
		}
	}
}

func (s *Service) update() {
	defer func() {
		if err := recover(); err != nil {
			server.Logger.Printf(errorServicePrefix+"%s update error: %v", s.Name, err)
		}
	}()
	//if pid := s.GetPid(); pid != 0 || !s.IsExist() {
	deploy := s.Deploy
	if deploy != nil && deploy.Provider != "" {
		var client git.Client
		switch strings.ToLower(deploy.Provider) {
		case "github":
			client = git.GithubClient(deploy.Token, deploy.Repository)
			break
		case "gitlab":
			client = git.GitlabClient(deploy.Token, deploy.Repository)
		}
		if client != nil {
			s.processGit(client)
		}
	}
	//}
}

func (s *Service) IsExist() bool {
	command := s.resoveCommand()
	if _, err := exec.LookPath(command); err == nil {
		return true
	}
	return false
}

func (s *Service) processGit(client git.Client) {
	//get content from remote git
	var (
		preVersion string
		version    string
		asset      string
		doPayload  = true
		err        error
	)

	defer func() {
		if doPayload {
			payloadUrl := s.Deploy.Payload
			if payloadUrl != "" {
				//Payload callback
				data := url.Values{}
				hostName, _ := os.Hostname()
				jsons := map[string]interface{}{
					"hostName": hostName,
					"name":     s.Name,
				}
				if err != nil {
					jsons["event"] = "update"
					jsons["status"] = "failed"
					jsons["error"] = err.Error()
				} else {
					jsons["event"] = "update"
					jsons["status"] = "success"
					jsons["preVersion"] = preVersion
					jsons["version"] = version
				}
				jsonb, _ := json.Marshal(jsons)
				data.Add("event", "update")
				data.Add("payload", string(jsonb))
				resp, err := http.PostForm(payloadUrl, data)
				if err != nil {
					server.Logger.Printf(errorServicePrefix+"%s payload: error: %v", s.Name, err)
				} else {
					resultData, _ := ioutil.ReadAll(resp.Body)
					if resp.StatusCode == 200 {
						server.Logger.Printf(infoServicePrefix+"%s payload: success: %s", s.Name, string(resultData))
					} else {
						server.Logger.Printf(errorServicePrefix+"%s payload: error: %s", s.Name, string(resultData))
					}
				}
			}
		}
	}()
	config := s
	t, v := getVersionType(config.Deploy.Version)
	switch t {
	case release:
		version, asset, err = client.GetRelease(v)
	case tag:
		version, asset, err = client.GetTag(v)
	case branch:
		if v != "" {
			version, asset, err = client.GetBranch(v)
		} else {
			version, asset, err = client.GetBranch(config.Deploy.Version)
		}
		break
	case latest:
		version, asset, err = client.GetRelease(config.Deploy.Version)
		break
	case content:
		arr := strings.Split(v, ":")
		version, err = client.GetContentFile(arr[0], strings.Join(arr[1:], ":"))
		version = strings.TrimSpace(version)
		version = strings.Trim(version, "\n")
		version = strings.Trim(version, "\r")

		if err != nil {
			server.Logger.Printf(errorServicePrefix+"%s get file error: %v", s.Name, err)
		}
		version, asset, err = client.GetRelease(version)
		break
	}
	if err != nil {
		server.Logger.Printf(errorServicePrefix+"%s find version error: %v", s.Name, err)
		return
	}
	//check local version
	preVersion = s.GetVersion()
	if preVersion == version {
		//server.Logger.Printf(infoServicePrefix+"%s preVersion=newVersion: %s", s.Name, version)
		doPayload = false
		return
	} else {
		server.Logger.Printf(infoServicePrefix+"%s find new version: %s", s.Name, version)
	}

	//download zip file and unzip
	//add dir
	if config.Deploy.ServicePath == "" {
		server.Logger.Printf(errorServicePrefix+"update %s error: no service_path in deploy", s.Name)
		return
	}
	dir := config.Deploy.ServicePath
	file := common.BinDir + "/update/" + s.Name + "/" + version + ".zip"

	//Termination download when shouldQuit close
	var quitLoop = make(chan bool)
	go func() {
		for {
			select {
			case <-quitLoop:
				return
			case <-server.ShutdownCh:
				client.Termination()
				server.Logger.Printf(infoServicePrefix+"update %s termination download", s.Name)
				return
			}
		}
	}()

	now := time.Now()
	progress := func(current, total int32) {
		s.Progress.Current = current
		s.Progress.Total = total
		s.Progress.InProgress = true

		//dispatch every 1 sec
		go func() {
			if time.Now().Sub(now).Seconds() > 1 {
				now = time.Now()
				server.watch.Dispatch(WatchParam{
					Name:      s.Name,
					WatchType: ServiceChanged,
				})
			}
		}()
	}

	err = client.DownloadFile(file, asset, progress)
	close(quitLoop)
	s.Progress.InProgress = false

	//update status
	server.watch.Dispatch(WatchParam{
		Name:      s.Name,
		WatchType: ServiceChanged,
	})

	if err != nil {
		server.Logger.Printf(errorServicePrefix+"update %s download error: %v", s.Name, err)
		return
	}
	err = common.Unzip(file, dir, true)
	if err != nil {
		server.Logger.Printf(errorServicePrefix+"update %s unzip error: %v", s.Name, err)
		return
	}
	s.SetVersion(version)
	err = s.Restart()
	if err != nil {
		server.Logger.Printf(errorServicePrefix+"restart %s error: %v", s.Name, err)

	} else {
		server.Logger.Printf(infoServicePrefix+"update %s success pre version: %s new version: %s", s.Name, preVersion, version)
	}
}

func (s *Service) GetVersion() string {
	versionPath := common.BinDir + "/run/" + s.Name + ".ver"
	content, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return ""
	}
	return string(content)
}
func (s *Service) SetVersion(version string) {
	versionPath := common.BinDir + "/run/" + s.Name + ".ver"
	data := []byte(version)
	os.MkdirAll(filepath.Dir(versionPath), 0755)
	ioutil.WriteFile(versionPath, data, 0666)
}

func (s *Service) Address() string {
	return net.JoinHostPort(s.Ip, strconv.Itoa(int(s.Port)))
}

func (s *Services) Keys() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *s {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *Service) runAtLoad() {
	if pid := s.GetPid(); pid == 0 && s.RunAtLoad {
		err := s.Start()
		if err != nil {
			server.Logger.Printf(errorServicePrefix+"%s running error: %v", s.Name, err)
		} else {
			server.Logger.Printf(infoServicePrefix+"%s running success", s.Name)
		}
	}
}

func (s *Service) SetPid(pid int) {
	pidString := []byte(strconv.Itoa(pid))
	os.MkdirAll(filepath.Dir(s.pidFile()), 0755)
	ioutil.WriteFile(s.pidFile(), pidString, 0666)
}

func (s *Service) RemovePid() error {
	return os.Remove(s.pidFile())
}

func (s *Service) GetPid() int {
	content, err := ioutil.ReadFile(s.pidFile())
	if err != nil {
		return 0
	} else {
		pid, _ := strconv.Atoi(strings.Trim(string(content), "\n"))
		if _, find := s.processExist(pid); find {
			f, err := os.Open(s.pidFile())
			if err == nil {
				fi, err := f.Stat()
				if err == nil {
					s.StartTime = fi.ModTime()
				}
			}
			return pid
		} else {
			return 0
		}
	}
}

func (s *Service) resovePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		if strings.HasPrefix(path, "."+string(os.PathSeparator)) {
			return common.BinDir + path[1:]
		} else {
			return common.BinDir + "/" + path
		}
	}
}

func (s *Service) resoveCommand() string {
	path := s.Command[0]
	if filepath.IsAbs(path) {
		return path
	} else {
		if strings.HasPrefix(path, "."+string(os.PathSeparator)) {
			return common.BinDir + path[1:]
		} else {
			return path
		}
	}
}

func (s *Service) pidFile() string {
	if s != nil && s.PidFile != "" {
		pid := s.resovePath(s.PidFile)
		return pid
	} else {
		return common.BinDir + "/run/" + s.Name + ".pid"
	}
}

func (s *Service) processExist(pid int) (*os.Process, bool) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			return nil, false
		}
	}
	return process, true
}
