package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/vars"
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

func (c *Client) Get(url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if !vars.NoDebugUrl[url] {
		// 保存请求响应日志
		log.Debug(fmt.Sprintf("Request URL: %s [GET]\nRequest Headers: %v\nResponseCode: %d\nResponse: %s",
			req.URL.String(),
			req.Header,
			resp.StatusCode,
			string(result)))
	}

	return result, resp.StatusCode, nil
}

func (c *Client) Post(url string, formData string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if !vars.NoDebugUrl[url] {
		// 保存请求响应日志
		log.Debug(fmt.Sprintf("Request URL: %s [POST]\nRequest Headers: %v\nRequest Body: %s\nResponseCode: %d\nResponse: %s",
			req.URL.String(),
			req.Header,
			formData,
			resp.StatusCode,
			string(result)))
	}

	return result, resp.StatusCode, nil
}

// SaveCookies 保存cookies
func (c *Client) SaveCookies(cfg *config.Config) error {
	urlStr := "https://newjw.hdu.edu.cn/jwglxt"
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	cookies := c.client.Jar.Cookies(parsedURL)
	for _, cookie := range cookies {
		if cookie.Name == "JSESSIONID" {
			cfg.Cookies.JSESSIONID = cookie.Value
		}
		if cookie.Name == "route" {
			cfg.Cookies.Route = cookie.Value
		}
	}
	// 保存配置文件
	err = config.SaveConfig(cfg)
	if err != nil {
		return err
	}

	return nil
}

// LoadCookies 加载cookies
func (c *Client) LoadCookies(cfg *config.Config) error {
	urlStr := "https://newjw.hdu.edu.cn/jwglxt"
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	cookies := []*http.Cookie{
		{
			Name:  "JSESSIONID",
			Value: cfg.Cookies.JSESSIONID,
		},
		{
			Name:  "route",
			Value: cfg.Cookies.Route,
		},
	}
	c.client.Jar.SetCookies(parsedURL, cookies)

	return nil
}
