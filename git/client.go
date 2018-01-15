package git

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Client interface {
	Request(method, url string) (string, error)
	GetContentFile(branch, file string) (string, error)
	GetRelease(release string) (string, string, error)
	GetBranch(branch string) (string, string, error)
	DownloadFile(file, url string) error
	Termination()
}

type download struct {
	cx     context.Context
	cancel context.CancelFunc
}

func (d *download) downloadFile(header, file, url string) error {
	// Create the file
	dir := filepath.Dir(file)
	os.MkdirAll(dir, 0755)

	// Get the data
	d.cx, d.cancel = context.WithCancel(context.Background())
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(d.cx)
	if header != "" {
		arr := strings.Split(header, ":")
		req.Header.Set(arr[0], strings.TrimSpace(arr[1]))
	}

	done := make(chan bool)

	var err error
	var resp *http.Response
	go func() {
		resp, err = http.DefaultClient.Do(req)
		done <- true
	}()

	select {
	case <-done:
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			// Writer the body to file
			out, err := os.Create(file)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				os.Remove(file)
				return err
			}
		} else {
			data, _ := ioutil.ReadAll(resp.Body)
			return errors.New(string(data))
		}
	case <-d.cx.Done():
		//canceled
		return d.cx.Err()
	}
	return nil
}
