package main

import (
	"errors"
	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/util"
	"strings"
)

// NewjwLogin newjw登录
func NewjwLogin(c *client.Client, cfg *config.Config) error {
	log.Info("获取csrftoken...")
	// 获取csrftoken
	csrftoken, err := c.GetCsrftoken()
	if err != nil {
		return err
	}

	log.Info("获取公钥...")
	// 获取公钥
	publicKey, err := c.GetPublicKey()
	if err != nil {
		return err
	}

	// 加密密码
	encryptedPassword, err := util.RsaEncrypt(publicKey.Modules, cfg.NewjwLogin.Password)
	if err != nil {
		return err
	}
	cfg.NewjwLogin.Password = encryptedPassword

	log.Info("正在登录...")
	// 登录
	loginReq := &client.LoginReq{
		Csrftoken: csrftoken,
		Username:  cfg.NewjwLogin.Username,
		Password:  cfg.NewjwLogin.Password,
	}
	result, err := c.NewjwLoginPost(loginReq)
	if err != nil {
		return err
	}

	// 判断是否登录成功
	if strings.Contains(result, "用户名或密码不正确，请重新输入") {
		return errors.New("用户名或密码不正确！")
	}

	return nil
}

// CasLogin cas登录
func CasLogin(c *client.Client, cfg *config.Config) error {
	log.Info("获取cas登录配置...")
	// 获取cas登录配置
	execution, croypto, err := c.GetCasLoginConfig()
	if err != nil {
		return err
	}

	// 加密密码
	encryptedPassword, err := util.DesEncrypt(croypto, cfg.CasLogin.Password)
	if err != nil {
		return err
	}
	cfg.CasLogin.Password = encryptedPassword

	log.Info("正在cas登录...")
	// cas登录
	casLoginReq := &client.CasLoginReq{
		Username:    cfg.CasLogin.Username,
		Type:        "UsernamePassword",
		EventID:     "submit",
		Geolocation: "",
		Execution:   execution,
		CaptchaCode: "",
		Croypto:     croypto,
		Password:    cfg.CasLogin.Password,
	}
	result, err := c.CasLoginPost(casLoginReq)
	if err != nil {
		return err
	}
	// 判断是否登录成功
	if strings.Contains(result, "用户名密码登录") {
		return errors.New("用户名或密码不正确！")
	}

	// 通过cas登录newjw
	log.Info("正在通过cas登录newjw...")
	result, err = c.CasLoginNewjw()
	if err != nil {
		return err
	}
	// 判断是否登录成功
	if !strings.Contains(result, "杭州电子科技大学本科教学管理服务平台") {
		return errors.New("未知错误, cas登录newjw失败, 请重新尝试登录")
	}

	return nil
}

// Login 根据Level优先级登录
func Login(cfg *config.Config) (*client.Client, error) {
	// 创建一个新的客户端
	c := client.NewClient(cfg)

	// 根据Level优先级登录
	if cfg.CasLogin.Level < cfg.NewjwLogin.Level {
		// cas登录
		log.Info("正在通过cas登录...")
		err := CasLogin(c, cfg)
		if err != nil {
			log.Error("cas登录失败: ", err)
			// newjw登录
			log.Info("正在通过newjw登录...")
			// 重置client
			c = client.NewClient(cfg)
			err := NewjwLogin(c, cfg)
			if err != nil {
				log.Error("newjw登录失败: ", err)
				return nil, err
			}
		}
	} else {
		// newjw登录
		log.Info("正在通过newjw登录...")
		err := NewjwLogin(c, cfg)
		if err != nil {
			log.Error("newjw登录失败: ", err)
			// cas登录
			log.Info("正在通过cas登录...")
			// 重置client
			c = client.NewClient(cfg)
			err := CasLogin(c, cfg)
			if err != nil {
				log.Error("cas登录失败: ", err)
				return nil, err
			}
		}
	}

	return c, nil
}
