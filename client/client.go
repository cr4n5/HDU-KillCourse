package client

import (
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type ClientBodyConfig struct {
	XkkzId map[string]string
	Ccdm   string
	BhId   string
	JgId   string
	Xsbj   string
	Xz     string
	Mzm    string
	Xslbdm string
	Xbm    string
	ZyfxId string
	XqhId  string
}

// Client 客户端结构体
type Client struct {
	client           *http.Client
	ClientBodyConfig *ClientBodyConfig
}

// NewClient 创建一个新的客户端
func NewClient(cfg *config.Config) *Client {
	// 创建一个cookie jar
	jar, _ := cookiejar.New(nil)
	return &Client{
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 保存请求响应日志
	log.Debug(fmt.Sprintf("Request URL: %s [GET]\nRequest Headers: %v\nResponse: %s",
		req.URL.String(),
		req.Header,
		string(result)))

	return result, nil
}

func (c *Client) Post(url string, formData string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 保存请求响应日志
	log.Debug(fmt.Sprintf("Request URL: %s [POST]\nRequest Headers: %v\nRequest Body: %s\nResponse: %s",
		req.URL.String(),
		req.Header,
		formData,
		string(result)))

	return result, nil
}
