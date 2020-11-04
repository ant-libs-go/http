/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-11-04 20:52:53
# File Name: http/http_test.go
# Description:
####################################################################### */

package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/config/parser"
	//. "github.com/smartystreets/goconvey/convey"
	//rds "github.com/gomodule/redigo/redis"
)

var globalCfg *config.Config

func TestMain(m *testing.M) {
	config.New(parser.NewTomlParser(),
		options.WithCfgSource("./test.toml"),
		options.WithCheckInterval(1))
	os.Exit(m.Run())
}

type Resp struct {
	A string
	B string
}

func TestBasic(t *testing.T) {
	r, err := Call("ping",
		map[string]string{
			"a": "aaa",
			"b": "bbb"},
		map[string]string{
			"c": "ccc",
			"d": "ddd",
		}, nil)
	defer r.Body.Close()
	if err == nil {
		b, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(b))
	}
	fmt.Println("---------")

	data := &Resp{}
	r, err = Call("ping",
		map[string]string{
			"a": "aaa",
			"b": "bbb"},
		map[string]string{
			"c": "ccc",
			"d": "ddd",
		}, data)
	if err == nil {
		b, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(b))
	}
	fmt.Println(data)

	/*
		Convey("TestEncode", t, func() {
			Convey("TestDecode err should return nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("TestDecode result should equal rawStr", func() {
				So(decStr, ShouldEqual, rawStr)
			})
		})
	*/
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
