/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-30 22:22:32
# File Name: server.go
# Description:
####################################################################### */

package server

import (
	"net/http"
	"time"
)

func NewRestServer(cfg *Cfg) (r *http.Server, err error) {
	r = &http.Server{Addr: cfg.DialAddr}
	if cfg.DialReadTimeout > 0 {
		r.ReadTimeout = cfg.DialReadTimeout * time.Millisecond
	}
	if cfg.DialWriteTimeout > 0 {
		r.ReadTimeout = cfg.DialWriteTimeout * time.Millisecond
	}
	if cfg.DialIdleTimeout > 0 {
		r.ReadTimeout = cfg.DialIdleTimeout * time.Millisecond
	}

	go func() {
		if err = r.ListenAndServe(); err != nil && err == http.ErrServerClosed {
			err = nil
		}
	}()
	time.Sleep(time.Second)
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
