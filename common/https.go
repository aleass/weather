package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HttpMethod string

const (
	PostType HttpMethod = "POST"
	GetType  HttpMethod = "GET"
)

func (m HttpMethod) toString() string {
	return string(m)
}

// method 请求方式
// url 链接
// body 请求体
// header 头信息
func HttpRequest(method HttpMethod, url string, body interface{}, header [][2]string) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(body)
	req, err := http.NewRequest(method.toString(), url, buff)
	if err != nil {
		return nil, err
	}
	for _, v := range header {
		req.Header.Set(v[0], v[1])
	}
	req.Header.Set("cookie", "Ares_SessionId=0eqec5blqamfekhsgzhblrx1b4c344a0e040b13")
	req.Header.Set("content-type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	_ = response.Body.Close()
	return raw, nil
}
