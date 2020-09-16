package api

import (
	"io"
	"net/http"
	"os/user"
	"path/filepath"
	"time"

	"github.com/jdxcode/netrc"
)

var (
	// BaseURL is the root of the Section API
	BaseURL = "https://aperture.section.io/api/v1"
)

func getBasicAuth() (u, p string, err error) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
	if err != nil {
		return u, p, err
	}
	u = n.Machine("aperture.section.io").Get("login")
	p = n.Machine("aperture.section.io").Get("password")
	return u, p, err
}

func request(method string, url string, body io.Reader) (resp *http.Response, err error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return resp, err
	}

	username, password, err := getBasicAuth()
	if err != nil {
		return resp, err
	}
	req.SetBasicAuth(username, password)

	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, err
}
