package http

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go.net/proxy"
)

var (
	defaultOptions = options{
		timeout: 60 * time.Second,
	}
	defaultRequestOptions = requestOptions{
		header:      http.Header{},
		contentType: ApplicationJSON,
		url:         "",
		body:        nil,
	}
	defaultDoOptions = doOptions{}
)

// Options ...
type Options func(o *options)

type options struct {
	timeout time.Duration

	socks5 proxy.Dialer
	proxy  string
}

func (o *options) ExecuteOptions(opt []Options) {
	for _, fn := range opt {
		fn(o)
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Options {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithProxy 配置 http 代理地址
//
// Note: 若 WithProxy 与 WithSOCKS5 同时使用, WithSOCKS5 将会优先使用
func WithProxy(proxy string) Options {
	return func(o *options) {
		o.proxy = proxy
	}
}

// WithSOCKS5 配置 socks 代理
func WithSOCKS5(dialer proxy.Dialer) Options {
	return func(o *options) {
		o.socks5 = dialer
	}
}

// RequestOptions ...
type RequestOptions func(o *requestOptions)

type requestOptions struct {
	header http.Header

	contentType string // 暂时不支持依据contentType做不同的解析
	url         string

	body interface{}
}

func (o *requestOptions) ExecuteOptions(opt []RequestOptions) {
	for _, fn := range opt {
		fn(o)
	}
}

// WithURL 配置请求的完整url
func WithURL(url string) RequestOptions {
	return func(o *requestOptions) {
		o.url = url
	}
}

// WithUrlBuild 组装url
// path 必须要以 '/' 开头, 否则某些情况下会报400错误
func WithUrlBuild(scheme string, host string, path string) RequestOptions {
	return func(o *requestOptions) {
		u := url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   path,
		}
		o.url = u.String()
	}
}

// WithContentType 设置 请求 Content-Type
func WithContentType(contentType string) RequestOptions {
	return func(o *requestOptions) {
		o.contentType = contentType
	}
}

// WithHeader 设置 请求头
func WithHeader(header http.Header) RequestOptions {
	return func(o *requestOptions) {
		o.header = header
	}
}

// WithBody 配置请求 body
func WithBody(body interface{}) RequestOptions {
	return func(o *requestOptions) {
		o.body = body
	}
}

// DoOptions ...
type DoOptions func(o *doOptions)

type doOptions struct {
	responseData   interface{}
	responseReader *io.Reader
	response       *[]byte
}

func (o *doOptions) ExecuteOptions(opt []DoOptions) {
	for _, fn := range opt {
		fn(o)
	}
}

// WithResponseBodyData 配置响应消息体数据
// data 将会根据响应消息的 Content-Type 反序列化
func WithResponseBodyData(data interface{}) DoOptions {
	return func(o *doOptions) {
		o.responseData = data
	}
}

// WithResponseBodyReader 配置响应消息体数据
func WithResponseBodyReader(data *io.Reader) DoOptions {
	return func(o *doOptions) {
		o.responseReader = data
	}
}

// WithResponseBody 配置响应消息体数据
func WithResponseBody(data *[]byte) DoOptions {
	return func(o *doOptions) {
		o.response = data
	}
}
