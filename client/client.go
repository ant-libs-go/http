/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-11-04 11:34:32
# File Name: client.go
# Description:
####################################################################### */

package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	ht "github.com/ant-libs-go/http"
	"github.com/ant-libs-go/util"
)

type REST_METHOD string

const (
	REST_METHOD_GET        REST_METHOD = "get"
	REST_METHOD_POST       REST_METHOD = "post"
	REST_METHOD_JSON_POST  REST_METHOD = "json_post"
	REST_METHOD_FILE_POST  REST_METHOD = "file_post"
	REST_METHOD_BYTES_POST REST_METHOD = "bytes_post"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36"
)

type RestClientPool struct {
	lock   sync.RWMutex
	client *http.Client
	cfg    *Cfg
}

func NewRestClientPool(cfg *Cfg) *RestClientPool {
	o := &RestClientPool{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.DialTimeout * time.Millisecond,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   cfg.DialTimeout * time.Millisecond, // 连接超时时间
					KeepAlive: 30 * time.Second,                   // 连接保持超时时间
				}).DialContext,
				DisableKeepAlives:   cfg.DialDisableKeepAlive,               // 是否开启长连接
				MaxIdleConns:        cfg.PoolMaxIdle,                        // 所有host最大空闲连接数
				MaxIdleConnsPerHost: cfg.PoolMaxIdlePerHost,                 // 每个host最大空闲连接数
				IdleConnTimeout:     cfg.PoolIdleTimeout * time.Millisecond, // 闲置连接的过期时间
			}}}
	return o
}

func (this *RestClientPool) SetHeader(key string, value string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.cfg.Headers == nil {
		this.cfg.Headers = make(map[string]string)
	}
	this.cfg.Headers[key] = value
}

// 注意：当resp不指定值时，需要手动进行response.Body.Close()
func (this *RestClientPool) Call(params interface{}, body interface{}, resp interface{}) (r *http.Response, err error) {
	var retry int
	for {
		if r, err = this.call(params, body); err == nil {
			break
		}
		retry += 1
		if retry >= this.cfg.FailRetry {
			err = errors.New(fmt.Sprintf("request failed, %v, retry completed", err))
			return
		}
		time.Sleep(this.cfg.FailRetryInterval * time.Millisecond)
	}

	if resp == nil {
		return
	}
	defer r.Body.Close()

	var b []byte
	if b, err = ioutil.ReadAll(r.Body); err == nil && len(b) > 0 {
		err = ht.Decode(this.cfg.Codec, b, resp)
	}
	return
}

func (this *RestClientPool) call(params interface{}, body interface{}) (r *http.Response, err error) {
	var req *http.Request
	req, err = this.buildRequest(params, body)
	if err == nil {
		r, err = this.client.Do(req)
	}
	if err == nil {
		if r.StatusCode != 200 && r.StatusCode != 204 {
			err = errors.New(fmt.Sprintf("http status is not equal 200/204, is: %d", r.StatusCode))
			return
		}
	}
	return
}

func (this *RestClientPool) buildRequest(params interface{}, body interface{}) (r *http.Request, err error) {
	var reader io.Reader
	if reader, err = this.buildBody(body); err != nil {
		return
	}

	r, err = http.NewRequest(this.buildMethod(), this.buildUrl(params), reader)
	if err == nil {
		for k, v := range this.buildHeaders() {
			r.Header.Set(k, v)
		}
		for k, v := range this.cfg.Headers {
			r.Header.Set(k, v)
		}
	}
	return
}

func (this *RestClientPool) buildMethod() (r string) {
	switch this.cfg.Method {
	case REST_METHOD_GET:
		r = http.MethodGet
	case REST_METHOD_POST, REST_METHOD_JSON_POST, REST_METHOD_FILE_POST, REST_METHOD_BYTES_POST:
		r = http.MethodPost
	default:
		r = http.MethodGet
	}
	return
}

func (this *RestClientPool) buildHeaders() (r map[string]string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	r = map[string]string{"User-Agent": DefaultUserAgent}
	switch this.cfg.Method {
	case REST_METHOD_GET, REST_METHOD_POST:
		r["Content-Type"] = "application/x-www-form-urlencoded"
	case REST_METHOD_JSON_POST:
		r["Content-Type"] = "application/json"
	case REST_METHOD_BYTES_POST:
		r["Content-Type"] = "application/octet-stream"
	}
	return
}

func (this *RestClientPool) buildUrl(inp interface{}) string {
	query := ""
	switch v := inp.(type) {
	case map[string]string:
		query = util.MapToQueryStr(v)
	case string:
		query = v
	}

	sep := "?"
	if strings.Index(this.cfg.Url, "?") > -1 {
		sep = "&"
	}
	return strings.Join([]string{this.cfg.Url, query}, sep)
}

func (this *RestClientPool) buildBody(inp interface{}) (r io.Reader, err error) {
	inp, err = ht.Encode(this.cfg.Codec, inp)
	if err != nil {
		return
	}

	switch v := inp.(type) {
	case map[string]string:
		r = strings.NewReader(util.MapToQueryStr(v))
	case string:
		r = strings.NewReader(v)
	case []byte:
		r = bytes.NewBuffer(v)
	case io.Reader:
		r = v
	}
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
