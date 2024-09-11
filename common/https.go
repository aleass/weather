package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

var AmtpApiCount, HeFengApiCount int

type HttpMethod string

const (
	PostType    HttpMethod = "POST"
	GetType                = "GET"
	MapType                = "Map"
	WeatherType            = "Weather"
	OtherType              = "Other"
)

func (m HttpMethod) toString() string {
	return string(m)
}

func proxy() *http.Client {
	// 设置代理
	proxyURL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return nil
	}

	// 创建 HTTP 客户端并设置代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	return client
}

// method 请求方式
// url 链接
// body 请求体
// header 头信息
func HttpRequest(types, method HttpMethod, url string, body interface{}, header [][2]string, isProxy bool, unmarsh any) ([]byte, error) {
	switch types {
	case MapType:
		AmtpApiCount++
	case WeatherType:
		HeFengApiCount++
	}
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(body)
	req, err := http.NewRequest(method.toString(), url, buff)
	if err != nil {
		return nil, err
	}
	for _, v := range header {
		req.Header.Set(v[0], v[1])
	}
	req.Header.Set("content-type", "application/json; charset=UTF-8")

	client := &http.Client{}
	if isProxy {
		client = proxy()
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	_ = response.Body.Close()
	if unmarsh == nil {
		return raw, nil
	}

	json.Unmarshal(raw, unmarsh)
	return raw, nil
}

func UploadFile(url string, form [][2]string, file io.Reader, fileFrom string) error {
	// 创建一个新的缓冲区
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文本字段
	for _, str := range form {
		err := writer.WriteField(str[0], str[1])
		if err != nil {
			return err
		}
	}

	// 创建一个文件字段
	part, err := writer.CreateFormFile(fileFrom, "temp")
	if err != nil {
		return err
	}

	// 将文件内容写入到表单字段
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// 结束表单
	err = writer.Close()
	if err != nil {
		return err
	}

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}

	// 设置请求的 Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := proxy()
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
