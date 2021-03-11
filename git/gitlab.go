package git

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gimke/riff/api"
	"io/ioutil"
	"net/http"
	"net/url"
)

var _ Client = &Gitlab{}

type Gitlab struct {
	Token      string
	Repository string
	d          *download
}

func (g *Gitlab) getUrl() string {
	u, _ := url.Parse(g.Repository)
	return u.Scheme + "://" + u.Host + "/api/v4/projects/" + url.PathEscape(u.Path[1:]) + "/repository"
}

func GitlabClient(token, repo string) Client {
	return &Gitlab{Token: token, Repository: repo}
}

func (g *Gitlab) Request(method, url string) (string, error) {
	req, _ := http.NewRequest(method, url, nil)
	if g.Token != "" {
		req.Header.Set("PRIVATE-TOKEN", g.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			//success
			if err == nil {
				return string(data), nil
			} else {
				return "", err
			}
		} else {
			return "", errors.New(string(data))
		}
	}
}

func (g *Gitlab) GetContentFile(branch, file string) (string, error) {
	u := g.getUrl()
	fmt.Println(u)
	u += "/files/" + url.PathEscape(file) + "?ref=" + url.PathEscape(branch)
	data, err := g.Request("GET", u)
	if err != nil {
		return "", err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		return "", err
	}
	decode, err := base64.StdEncoding.DecodeString(jsonData["content"].(string))
	content := string(decode)
	return content, nil
}

func (g *Gitlab) GetRelease(release string) (string, string, error) {
	//latest or tag name
	u := g.getUrl()
	//tag := g.Version
	if release == "latest" {
		u += "/tags"
	} else {
		u += "/tags/" + release
	}
	data, err := g.Request("GET", u)
	if err != nil {
		return "", "", err
	}
	if release == "latest" {
		var jsonData []map[string]interface{}
		err = json.Unmarshal([]byte(data), &jsonData)
		if err != nil {
			return "", "", err
		}
		if len(jsonData) > 0 {
			version := jsonData[0]["name"].(string)
			sha := jsonData[0]["commit"].(map[string]interface{})["id"].(string)
			zipball := g.getUrl() + "/archive.zip?sha=" + sha
			return version, zipball, nil
		} else {
			return "", "", errors.New("not found")
		}

	} else {
		var jsonData map[string]interface{}
		err = json.Unmarshal([]byte(data), &jsonData)
		if err != nil {
			return "", "", err
		}
		version := jsonData["name"].(string)
		sha := jsonData["commit"].(map[string]interface{})["id"].(string)
		zipball := g.getUrl() + "/archive.zip?sha=" + sha

		return version, zipball, nil
	}

}

func (g *Gitlab) GetBranch(branch string) (string, string, error) {
	u := g.getUrl()
	//branche := g.Version
	u += "/branches/" + branch

	data, err := g.Request("GET", u)
	if err != nil {
		return "", "", err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		return "", "", err
	}
	version := jsonData["commit"].(map[string]interface{})["id"].(string)
	asset := g.getUrl() + "/archive.zip?sha=" + version

	return version, asset, nil
}

func (g *Gitlab) DownloadFile(file, url string, progress api.Progress) error {
	header := "PRIVATE-TOKEN: " + g.Token
	g.d = &download{}
	return g.d.downloadFile(header, file, url, progress)
}

func (g *Gitlab) Termination() {
	//Termination download
	if g.d != nil && g.d.cancel != nil {
		g.d.cancel()
	}
}
