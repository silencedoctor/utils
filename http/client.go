package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client http 调用通用客户端
type Client struct {
	HTTPClient http.Client
}

// NewClient 创建 client
func NewClient(opts ...Options) *Client {
	opt := defaultOptions
	opt.ExecuteOptions(opts)

	c := &Client{}

	// Proxy
	if opt.socks5 != nil {
		ts := &http.Transport{
			Dial: opt.socks5.Dial,
		}
		c.HTTPClient.Transport = ts
	} else if len(opt.proxy) > 0 {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(opt.proxy)
		}

		ts := &http.Transport{
			Proxy: proxy,
		}
		c.HTTPClient.Transport = ts
	}

	// timeout
	c.HTTPClient.Timeout = opt.timeout

	return c
}

// NewRequest 新建请求
func (c *Client) NewRequest(method string, opts ...RequestOptions) (*http.Request, error) {
	opt := defaultRequestOptions
	opt.ExecuteOptions(opts)

	b, err := json.Marshal(opt.body)
	if err != nil {
		return nil, fmt.Errorf("marshal body err %v", err)
	}
	body := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, opt.url, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest err %v", err)
	}

	req.Header = opt.header

	// http Content-Type
	req.Header.Set("Content-Type", opt.contentType)

	return req, nil
}

// Do 执行请求
// resp.Body 已统一关闭, 调用者不需要再关闭
// 解析 body 失败将会返回 errorutils.SyntaxError
func (c *Client) Do(req *http.Request, opts ...DoOptions) (*Response, error) {
	opt := defaultDoOptions
	opt.ExecuteOptions(opts)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request err %v", err)
	}
	defer resp.Body.Close()

	// 必须全部读完, 否则会关闭连接无法使用长连接方式
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body err %v", err)
	}

	if opt.responseReader != nil {
		*opt.responseReader = bytes.NewBuffer(body)
	} else if opt.response != nil {
		*opt.response = body
	} else if opt.responseData != nil {
		err = json.Unmarshal(body, opt.responseData)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal err %v body %s opt.responseData %v url %v", err, body, opt.responseData, req.URL)
		}
	}

	return (*Response)(resp), nil
}
