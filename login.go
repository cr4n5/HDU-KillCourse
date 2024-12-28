package main

import (
	"HDU-KillCourse/client"
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"HDU-KillCourse/util"
	"errors"
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
	encryptedPassword, err := util.RsaEncrypt(publicKey.Modules, cfg.Login.Password)
	if err != nil {
		return err
	}
	cfg.Login.Password = encryptedPassword

	log.Info("正在登录...")
	// 登录
	loginReq := &client.LoginReq{
		Csrftoken: csrftoken,
		Username:  cfg.Login.Username,
		Password:  cfg.Login.Password,
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
	encryptedPassword, err := util.DesEncrypt(croypto, cfg.Login.Password)
	if err != nil {
		return err
	}
	cfg.Login.Password = encryptedPassword

	log.Info("正在cas登录...")
	// cas登录
	casLoginReq := &client.CasLoginReq{
		Username:    cfg.Login.Username,
		Type:        "UsernamePassword",
		EventID:     "submit",
		Geolocation: "",
		Execution:   execution,
		CaptchaCode: "",
		Croypto:     croypto,
		Password:    cfg.Login.Password,
	}
	result, err := c.CasLoginPost(casLoginReq)
	if err != nil {
		return err
	}
	// 判断是否登录成功
	if strings.Contains(result, "用户名密码登录") {
		return errors.New("cas登录失败, 请检查用户名和密码是否正确")
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
