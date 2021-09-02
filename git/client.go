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
	"time"

	"github.com/teatak/riff/api"
)

type WriteCounter struct {
	Total    int32
	Current  int32
	Progress api.Progress
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += int32(n)
	wc.Progress(wc.Current, wc.Total)
	return n, nil
}

type Client interface {
	Request(method, url string) (string, error)
	GetContentFile(branch, file string) (string, error)
	GetRelease(release string) (string, string, error)
	GetTag(tag string) (string, string, error)
	GetBranch(branch string) (string, string, error)
	DownloadFile(file, url string, progress api.Progress) error
	Termination()
}

type download struct {
	cx     context.Context
	cancel context.CancelFunc
}

func (d *download) downloadFile(header, file, url string, progress api.Progress) error {
	// Create the file
	dir := filepath.Dir(file)
	os.MkdirAll(dir, 0755)

	// Get the data 5 mins time out
	d.cx, d.cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer d.cancel()
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
			total := resp.ContentLength
			counter := &WriteCounter{Total: int32(total), Current: 0, Progress: progress}
			_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
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
