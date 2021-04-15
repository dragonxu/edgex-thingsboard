package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const contentTypeJSON = "application/json"

type HTTPError struct {
	Status int
	Body   []byte
}

func NewHTTPError(code int, body []byte) error {
	return &HTTPError{
		Status: code,
		Body:   body,
	}
}

func (e HTTPError) HTTPStatus() int {
	return e.Status
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("code=%d, body=%s", e.Status, e.Body)
}

func GetHTTPStatus(err error) int {
	if httpErr, ok := err.(interface {
		HTTPStatus() int
	}); ok {
		return httpErr.HTTPStatus()
	}

	if err == nil {
		return http.StatusOK
	} else {
		return http.StatusInternalServerError
	}
}

type RawMessage []byte

// 发起HTTP JSON请求
//
// method 请求方法
// url 请求地址
// timeout 超时时间
// payload 请求数据。会序列化为JSON字节数组（除非 payload 为 RawMessage 类型）
// result 响应数据。会对响应数据进行JSON反序列化（除非 result 为 RawMessage 类型），并将结果存入result
func RequestJSON(method string, url string, timeout time.Duration, requestBody interface{}, responseBody interface{}) error {
	var (
		req     *http.Request
		rep     *http.Response
		err     error
		binBody *bytes.Reader
	)
	if requestBody != nil {
		var buf []byte
		switch requestBody.(type) {
		case RawMessage:
			buf = responseBody.(RawMessage)
		case *RawMessage:
			if m := responseBody.(*RawMessage); m != nil {
				buf = *m
			}
		default:
			buf, _ = json.Marshal(requestBody)
		}
		binBody = bytes.NewReader(buf)
	} else {
		binBody = bytes.NewReader(make([]byte, 0))
	}
	req, err = http.NewRequest(method, url, binBody)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentTypeJSON)
	req.Header.Add("Accept", contentTypeJSON)
	client := http.DefaultClient
	client.Timeout = timeout
	rep, err = client.Do(req)
	if err != nil {
		return err
	}
	defer rep.Body.Close()

	status := rep.StatusCode
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return NewHTTPError(status, body)
	}
	if body == nil || responseBody == nil {
		return nil
	}

	switch responseBody.(type) {
	case RawMessage:
		return errors.New("responseBody not settable")
	case *RawMessage:
		rv := reflect.ValueOf(responseBody).Elem()
		rv.Set(reflect.ValueOf(RawMessage(body)))
	default:
		contentType := rep.Header.Get("Content-Type")
		if !strings.Contains(contentType, contentTypeJSON) {
			return nil // 忽略非JSON格式的响应数据
		}

		if err = json.Unmarshal(body, responseBody); err != nil {
			return err
		}
	}
	return nil
}
