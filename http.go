/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-11-04 18:37:19
# File Name: http.go
# Description:
####################################################################### */

package http

import (
	"fmt"

	"github.com/ant-libs-go/util"
	"google.golang.org/protobuf/proto"
)

type Codec string

const (
	CODEC_NOSET  Codec = ""
	CODEC_JSON   Codec = "json"
	CODEC_PB     Codec = "pb"
	CODEC_GOB    Codec = "gob"
	CODEC_THRIFT Codec = "thrift"
)

func Encode(codec Codec, inp interface{}) (r []byte, err error) {
	switch codec {
	case CODEC_JSON:
		r, err = util.JsonEncode(inp)
	case CODEC_PB:
		r, err = util.PbEncode(inp.(proto.Message))
	case CODEC_GOB:
		r, err = util.GobEncode(inp)
	case CODEC_THRIFT:
		r, err = util.ThriftEncode(inp)
	default:
		err = fmt.Errorf("codec#%s not support")
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
	default:
		err = fmt.Errorf("codec#%s not support")
	}
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
