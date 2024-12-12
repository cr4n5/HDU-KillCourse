package main

import (
	"HDU-KillCourse/client"
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"HDU-KillCourse/util"
	"errors"
	"strings"
)

// Login 登录
func Login(c *client.Client, cfg *config.Config) error {
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
	result, err := c.LoginPost(loginReq)
	if err != nil {
		return err
	}

	// 判断是否登录成功
	if strings.Contains(result, "用户名或密码不正确，请重新输入") {
		return errors.New("用户名或密码不正确！")
	}

	return nil
}
