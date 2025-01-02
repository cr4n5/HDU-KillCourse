package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// GetCsrftoken 获取csrftoken
func (c *Client) GetCsrftoken() (string, error) {
	login_url := "https://newjw.hdu.edu.cn/jwglxt/xtgl/login_slogin.html"

	// 获取csrftoken
	result, err := c.Get(login_url)
	if err != nil {
		return "", err
	}
	// 解析csrftoken
	doc, err := htmlquery.Parse(strings.NewReader(string(result)))
	if err != nil {
		return "", err
	}
	node := htmlquery.FindOne(doc, `//input[@name="csrftoken"]/@value`)
	var csrftoken string
	if node != nil {
		csrftoken = htmlquery.InnerText(node)
	} else {
		return "", errors.New("获取csrftoken失败")
	}

	return csrftoken, nil
}

// GetCasLoginConfig 获取cas登录配置
func (c *Client) GetCasLoginConfig() (execution string, croypto string, err error) {
	result, err := c.Get("https://sso.hdu.edu.cn/login")
	if err != nil {
		return "", "", err
	}

	// 解析cas登录配置
	doc, err := htmlquery.Parse(strings.NewReader(string(result)))
	if err != nil {
		return "", "", err
	}
	node := htmlquery.FindOne(doc, `//*[@id="login-page-flowkey"]/text()`)
	if node != nil {
		execution = htmlquery.InnerText(node)
	} else {
		return "", "", errors.New("获取cas登录配置失败")
	}
	node = htmlquery.FindOne(doc, `//*[@id="login-croypto"]/text()`)
	if node != nil {
		croypto = htmlquery.InnerText(node)
	} else {
		return "", "", errors.New("获取cas登录配置失败")
	}
	if execution == "" || croypto == "" {
		return "", "", errors.New("获取cas登录配置失败")
	}

	return execution, croypto, nil
}

// CasLoginPost cas登录请求
func (c *Client) CasLoginPost(req *CasLoginReq) (string, error) {
	login_url := "https://sso.hdu.edu.cn/login"
	// 登录
	formData := req.ToFormData()
	result, err := c.Post(login_url, formData.Encode())
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// CasLoginNewjw cas登录newjw
func (c *Client) CasLoginNewjw() (string, error) {
	new_jw := "https://newjw.hdu.edu.cn/sso/driot4login"
	// 通过cas登录newjw
	result, err := c.Get(new_jw)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GetPublicKey 获取公钥
func (c *Client) GetPublicKey() (*GetPublicKeyResp, error) {
	result, err := c.Get(fmt.Sprintf("https://newjw.hdu.edu.cn/jwglxt/xtgl/login_getPublicKey.html?time=%d", time.Now().Unix()))
	if err != nil {
		return nil, err
	}
	// 解析公钥
	var PublicKeyResp GetPublicKeyResp
	err = json.Unmarshal(result, &PublicKeyResp)
	if err != nil {
		return nil, err
	}

	return &PublicKeyResp, nil
}

// NewjwLoginPost Newjw登录请求
func (c *Client) NewjwLoginPost(req *LoginReq) (string, error) {
	login_url := "https://newjw.hdu.edu.cn/jwglxt/xtgl/login_slogin.html"
	// 登录
	formData := req.ToFormData()
	result, err := c.Post(login_url, formData.Encode())
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GetCourse 获取课程
func (c *Client) GetCourse(req *GetCourseReq) (*GetCourseResp, error) {
	course_url := "https://newjw.hdu.edu.cn/jwglxt/rwlscx/rwlscx_cxRwlsIndex.html?doType=query&gnmkdm=N1548"
	// 获取课程
	formData := req.ToFormData()
	result, err := c.Post(course_url, formData.Encode())
	if err != nil {
		return nil, err
	}
	// 检验是否可以获取课程
	if strings.Contains(string(result), "无功能权限") {
		return nil, errors.New("任务落实查询并未开放")
	}
	// 解析课程
	var courseResp GetCourseResp
	err = json.Unmarshal(result, &courseResp)
	if err != nil {
		return nil, err
	}

	return &courseResp, nil
}

// GetClientBodyConfig 获取选课配置
func (c *Client) GetClientBodyConfig() error {
	// 解析选课配置函数
	AnalysisConfig := func(doc *html.Node) error {
		// 解析XkkzId
		AnalysisXkkzId := func(node []*html.Node) error {
			pattern := regexp.MustCompile(`queryCourse\(this,'(\d+)'`)
			pattern1 := regexp.MustCompile(`queryCourse\(this,'(?:[^']*)','(\w+)'`)
			for _, n := range node {
				match := pattern.FindStringSubmatch(htmlquery.InnerText(n))
				match1 := pattern1.FindStringSubmatch(htmlquery.InnerText(n))
				if match == nil || match1 == nil {
					return errors.New("XkkzId解析失败") // 待测试
				}
				c.ClientBodyConfig.XkkzId[match[1]] = match1[1]
			}
			return nil
		}

		c.ClientBodyConfig = &ClientBodyConfig{
			XkkzId: make(map[string]string),
		}
		// 如果未找到对应的xpath，报错
		if node := htmlquery.FindOne(doc, `//input[@name="ccdm"]/@value`); node != nil {
			c.ClientBodyConfig.Ccdm = htmlquery.InnerText(node)
		} else {
			return errors.New("ccdm获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="bh_id"]/@value`); node != nil {
			c.ClientBodyConfig.BhId = htmlquery.InnerText(node)
		} else {
			return errors.New("bh_id获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="jg_id_1"]/@value`); node != nil {
			c.ClientBodyConfig.JgId = htmlquery.InnerText(node)
		} else {
			return errors.New("jg_id获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="xsbj"]/@value`); node != nil {
			c.ClientBodyConfig.Xsbj = htmlquery.InnerText(node)
		} else {
			return errors.New("xsbj获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="xz"]/@value`); node != nil {
			c.ClientBodyConfig.Xz = htmlquery.InnerText(node)
		} else {
			return errors.New("xz获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="mzm"]/@value`); node != nil {
			c.ClientBodyConfig.Mzm = htmlquery.InnerText(node)
		} else {
			return errors.New("mzm获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="xslbdm"]/@value`); node != nil {
			c.ClientBodyConfig.Xslbdm = htmlquery.InnerText(node)
		} else {
			return errors.New("xslbdm获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="xbm"]/@value`); node != nil {
			c.ClientBodyConfig.Xbm = htmlquery.InnerText(node)
		} else {
			return errors.New("xbm获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="zyfx_id"]/@value`); node != nil {
			c.ClientBodyConfig.ZyfxId = htmlquery.InnerText(node)
		} else {
			return errors.New("zyfx_id获取失败")
		}
		if node := htmlquery.FindOne(doc, `//input[@name="xqh_id"]/@value`); node != nil {
			c.ClientBodyConfig.XqhId = htmlquery.InnerText(node)
		} else {
			return errors.New("xqh_id获取失败")
		}

		// 获取选课控制ID
		if node := htmlquery.Find(doc, `//a[@role="tab"]/@onclick`); node != nil {
			err := AnalysisXkkzId(node)
			if err != nil {
				return err
			}
		} else {
			return errors.New("XkkzID获取失败")
		}

		return nil
	}

	url := "https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzb_cxZzxkYzbIndex.html?gnmkdm=N253512&layout=default"
	// 获取选课配置
	result, err := c.Get(url)
	if err != nil {
		return err
	}
	// 检验是否可以获取选课配置
	if strings.Contains(string(result), "对不起，当前不属于选课阶段") {
		return errors.New("当前不属于选课阶段")
	}
	// 解析选课配置
	doc, err := htmlquery.Parse(strings.NewReader(string(result)))
	if err != nil {
		return err
	}
	err = AnalysisConfig(doc)
	if err != nil {
		return err
	}

	return nil
}

// GetDoJxbId 获取do_jxb_id
func (c *Client) GetDoJxbId(req *GetDoJxbIdReq) ([]GetDoJxbIdResp, error) {
	url := "https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzbjk_cxJxbWithKchZzxkYzb.html?gnmkdm=N253512"
	// 获取do_jxb_id
	formData := req.ToFormData()
	result, err := c.Post(url, formData.Encode())
	if err != nil {
		return nil, err
	}
	// 解析do_jxb_id
	var doJxbIdResp []GetDoJxbIdResp
	err = json.Unmarshal(result, &doJxbIdResp)
	if err != nil {
		return nil, err
	}

	return doJxbIdResp, nil
}

// SelectCourse 选课
func (c *Client) SelectCourse(req *SelectCourseReq) (*SelectCourseResq, error) {
	url := "https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzbjk_xkBcZyZzxkYzb.html?gnmkdm=N253512"
	// 选课
	formData := req.ToFormData()
	result, err := c.Post(url, formData.Encode())
	if err != nil {
		return nil, err
	}
	// 解析选课结果
	var selectCourseResq SelectCourseResq
	err = json.Unmarshal(result, &selectCourseResq)
	if err != nil {
		return nil, err
	}

	return &selectCourseResq, nil
}

// CancelCourse 退课
func (c *Client) CancelCourse(req *CancelCourseReq) (string, error) {
	url := "https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzb_tuikBcZzxkYzb.html?gnmkdm=N253512"
	// 退课
	formData := req.ToFormData()
	result, err := c.Post(url, formData.Encode())
	if err != nil {
		return "", err
	}

	return string(result), nil
}
