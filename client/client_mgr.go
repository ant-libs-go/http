/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-11-04 11:34:51
# File Name: client_mgr.go
# Description:
####################################################################### */

package client

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
	ht "github.com/ant-libs-go/http"
)

var (
	once  sync.Once
	lock  sync.RWMutex
	pools map[string]*RestClientPool
)

func init() {
	pools = map[string]*RestClientPool{}
}

type restConfig struct {
	Rest *struct {
		Cfgs map[string]*Cfg `toml:"client"`
	} `toml:"rest"`
}

type Cfg struct {
	// basic
	Url               string            `toml:"url"`
	Codec             ht.Codec          `toml:"codec"`
	Method            REST_METHOD       `toml:"method"`
	Headers           map[string]string `toml:"headers"`
	FailRetry         int               `toml:"fail_retry"`
	FailRetryInterval time.Duration     `toml:"fail_retry_interval"`

	// dial
	DialTimeout          time.Duration `toml:"dial_timeout"` // 连接超时时间
	DialDisableKeepAlive bool          `toml:"dial_disable_keep_alive"`

	// pool
	PoolMaxIdle        int           `toml:"pool_max_idle"`          // 所有host最大空闲连接数
	PoolMaxIdlePerHost int           `toml:"pool_max_idle_per_host"` // 每个host最大空闲连接数
	PoolIdleTimeout    time.Duration `toml:"pool_idle_time"`         // 闲置连接的过期时间
}

func Call(name string, params interface{}, body interface{}, resp interface{}) (r *http.Response, err error) {
	var cli *RestClientPool
	cli, err = SafePool(name)
	if err == nil {
		r, err = cli.Call(params, body, resp)
	}
	return
}

func Pool(name string) (r *RestClientPool) {
	var err error
	if r, err = getPool(name); err != nil {
		panic(err)
	}
	return
}

func SafePool(name string) (r *RestClientPool, err error) {
	return getPool(name)
}

func getPool(name string) (r *RestClientPool, err error) {
	lock.RLock()
	r = pools[name]
	lock.RUnlock()
	if r == nil {
		r, err = addPool(name)
	}
	return
}

func addPool(name string) (r *RestClientPool, err error) {
	var cfg *Cfg
	if cfg, err = LoadCfg(name); err != nil {
		return
	}
	r = NewRestClientPool(cfg)

	lock.Lock()
	pools[name] = r
	lock.Unlock()
	return
}

func LoadCfg(name string) (r *Cfg, err error) {
	var cfgs map[string]*Cfg
	if cfgs, err = loadCfgs(); err != nil {
		return
	}
	if r = cfgs[name]; r == nil {
		err = fmt.Errorf("rest#%s not configed", name)
		return
	}
	return
}

func loadCfgs() (r map[string]*Cfg, err error) {
	r = map[string]*Cfg{}

	cfg := &restConfig{}
	once.Do(func() {
		_, err = config.Load(cfg, options.WithOnChangeFn(func(cfg interface{}) {
			lock.Lock()
			defer lock.Unlock()
			pools = map[string]*RestClientPool{}
		}))
	})

	cfg = config.Get(cfg).(*restConfig)
	if err == nil && (cfg.Rest == nil || cfg.Rest.Cfgs == nil || len(cfg.Rest.Cfgs) == 0) {
		err = fmt.Errorf("not configed")
	}
	if err != nil {
		err = fmt.Errorf("rest load cfgs error, %s", err)
		return
	}
	r = cfg.Rest.Cfgs
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
