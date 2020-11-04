/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-11-04 18:37:19
# File Name: http.go
# Description:
####################################################################### */

package http

import (
	"github.com/ant-libs-go/util"
	"github.com/golang/protobuf/proto"
)

type Codec string

const (
	CODEC_NOSET  Codec = ""
	CODEC_JSON   Codec = "json"
	CODEC_PB     Codec = "pb"
	CODEC_GOB    Codec = "gob"
	CODEC_THRIFT Codec = "thrift"
)

func Encode(codec Codec, inp interface{}) (r interface{}, err error) {
	switch codec {
	case CODEC_JSON:
		r, err = util.JsonEncode(inp)
	case CODEC_PB:
		r, err = util.PbEncode(inp.(proto.Message))
	case CODEC_GOB:
		r, err = util.GobEncode(inp)
	case CODEC_THRIFT:
		r, err = util.ThriftEncode(inp)
	}
	return
}

func Decode(codec Codec, b []byte, inp interface{}) (err error) {
	switch codec {
	case CODEC_JSON:
		err = util.JsonDecode(b, inp)
	case CODEC_PB:
		err = util.PbDecode(b, inp.(proto.Message))
	case CODEC_GOB:
		err = util.GobDecode(b, inp)
	case CODEC_THRIFT:
		err = util.ThriftDecode(b, inp)
	}
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
