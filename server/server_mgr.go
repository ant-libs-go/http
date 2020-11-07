/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-30 22:01:02
# File Name: server_mgr.go
# Description:
####################################################################### */

package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/safe_stop"
	"github.com/smallnest/rest/server"
)

var (
	once    sync.Once
	lock    sync.RWMutex
	servers map[string]*server.Server
)

func init() {
	servers = map[string]*server.Server{}
}

type restConfig struct {
	Rest *struct {
		Cfgs map[string]*Cfg `toml:"server"`
	} `toml:"rest"`
}

type Cfg struct {
	// dial
	DialAddr         string        `toml:"addr"`
	DialReadTimeout  time.Duration `toml:"read_timeout"`
	DialWriteTimeout time.Duration `toml:"write_timeout"`
	DialIdleTimeout  time.Duration `toml:"idle_timeout"`
}

func StartDefaultServer(rcvr interface{}) (err error) {
	return StartServer("default", rcvr)
}

func StopDefaultServer() (err error) {
	return StopServer("default")
}

func DefaultServer() (r *server.Server) {
	return Server("default")
}

func StartServer(name string, rcvr interface{}) (err error) {
	safe_stop.Lock(1)
	var srv *server.Server
	if srv, err = SafeServer(name); err == nil {
		srv.Handler = rcvr
	}
	return
}

func StopServer(name string) (err error) {
	defer safe_stop.Unlock()
	var srv *server.Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if srv, err = SafeServer(name); err == nil {
		err = srv.Shutdown(ctx)
	}
	return
}

func Server(name string) (r *server.Server) {
	var err error
	if r, err = getServer(name); err != nil {
		panic(err)
	}
	return
}

func SafeServer(name string) (r *server.Server, err error) {
	return getServer(name)
}

func getServer(name string) (r *server.Server, err error) {
	lock.RLock()
	r = servers[name]
	lock.RUnlock()
	if r == nil {
		r, err = addServer(name)
	}
	return
}

func addServer(name string) (r *server.Server, err error) {
	var cfg *Cfg
	if cfg, err = loadCfg(name); err != nil {
		return
	}
	if r, err = NewRestServer(cfg); err != nil {
		return
	}

	lock.Lock()
	servers[name] = r
	lock.Unlock()
	return
}

func loadCfg(name string) (r *Cfg, err error) {
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
		_, err = config.Load(cfg)
	})

	config.Get(cfg)
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
