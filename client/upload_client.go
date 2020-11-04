/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-02-27 16:03:53
# File Name: upload_client.go
# Description:
####################################################################### */

package client

/*
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"
)

func UploadRequest(url string, file string, body map[string]string, resp interface{}, client *HttpClient, opts ...RequestOption) (err error) {
	buf := &bytes.Buffer{}
	bufWriter := multipart.NewWriter(buf)

	fileWriter, _ := bufWriter.CreateFormFile("file", fmt.Sprintf("%d%s", time.Now().Unix(), path.Ext(file)))
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(fileWriter, f)

	for k, v := range body {
		_ = bufWriter.WriteField(k, v)
	}
	bufWriter.Close()

	opts = append(opts, WithHeaders(map[string]string{"Content-Type": bufWriter.FormDataContentType()}))
	opts = append(opts, WithBody(buf.String()))
	httpReq := NewRequest(url, HTTP_METHOD_FILE_POST, opts...)

	if client == nil {
		client = DefaultHttpClient.clone()
	}
	if err = client.SetRequest(httpReq).Do(); err != nil {
		return
	}
	httpResp := client.GetResponse()
	if httpResp.Status != 200 {
		err = errors.New(fmt.Sprintf("http status is not equal 200, is: %d", httpResp.Status))
		return
	}
	err = json.Unmarshal(httpResp.Body, resp)
	return
}
*/

// vim: set noexpandtab ts=4 sts=4 sw=4 :
