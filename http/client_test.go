package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := NewClient(WithTimeout(time.Second * 20))
	url := fmt.Sprintf("http://www.baidu.com")
	req, err := c.NewRequest(http.MethodGet, WithURL(url))
	if err != nil {
		t.Error(err)
		return
	}

	var data interface{}
	resp, err := c.Do(req, WithResponseBodyData(data))
	if err != nil || !resp.IsOK() {
		t.Log(fmt.Errorf("http do fail or response status err %v, url %v, status %v", err, req, resp))
	}

	t.Log(data)
	t.Log(resp)
}

func TestClient2(t *testing.T) {
	c := NewClient(WithTimeout(time.Second * 20))

	header := http.Header{}
	header.Add("token", "sssss")

	req, err := c.NewRequest(http.MethodGet,
		WithUrlBuild("http", "www.baidu.com", ""),
		WithHeader(header))
	if err != nil {
		t.Error(err)
		return
	}

	var data interface{}
	resp, err := c.Do(req, WithResponseBodyData(data))
	if err != nil || !resp.IsOK() {
		t.Log(fmt.Errorf("http do fail or response status err %v, url %v, status %v", err, req.URL, resp.Status))
	}

	t.Log(data)
	t.Log(resp)
}

func TestClient3(t *testing.T) {
	c := NewClient(WithTimeout(time.Second * 20))
	url := fmt.Sprintf("http://www.baidu.com")
	req, err := c.NewRequest(http.MethodPost, WithURL(url))
	if err != nil {
		t.Error(err)
		return
	}

	var data interface{}
	resp, err := c.Do(req, WithResponseBodyData(data))
	if err != nil || !resp.IsOK() {
		t.Log(fmt.Errorf("http do fail or response status err %v, url %v, status %v", err, req, resp))
	}

	t.Log(data)
	t.Log(resp)
}
