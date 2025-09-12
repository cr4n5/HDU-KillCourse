package login

import (
	"errors"
	"strings"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/util"
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

	log.Info("正在登录...")
	// 登录
	loginReq := &client.LoginReq{
		Csrftoken: csrftoken,
		Username:  cfg.NewjwLogin.Username,
		Password:  encryptedPassword,
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

// CasPassWordLogin cas使用用户名密码登录
func CasPassWordLogin(c *client.Client, cfg *config.Config) error {
	log.Info("获取cas登录配置...")
	// 获取cas登录配置
	execution, croypto, err := c.GetCasLoginConfig()
	if err != nil {
		return err
	}

	// 加密密码
	encryptedPassword, err := util.AesEncrypt(croypto, cfg.CasLogin.Password)
	if err != nil {
		return err
	}

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
		Password:    encryptedPassword,
	}
	result, err := c.CasLoginPost(casLoginReq)
	if err != nil {
		return err
	}
	// 判断是否登录成功
	if strings.Contains(result, "统一身份认证") {
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

// CasQrLogin cas使用钉钉扫码登录
func CasQrLogin(c *client.Client, cfg *config.Config) error {
	log.Info("获取cas登录配置...")
	// 获取cas登录配置
	execution, _, err := c.GetCasLoginConfig()
	if err != nil {
		return err
	}

	log.Info("正在获取QrLoginId...")
	// 获取QrLoginId
	qrLoginIdResp, err := c.GetQrLoginId()
	if err != nil {
		return err
	}

	log.Info("正在获取二维码...")
	log.Info("请使用" + log.ErrorColor("钉钉") + "扫码登录...")
	var qrLoginStatus *client.QrLoginStatusResp
	for {
		// 获取二维码
		qrCodeBytes, err := c.GetQrCode(qrLoginIdResp.Data)
		if err != nil {
			return err
		}
		// 解析二维码
		qrCode, err := util.QrCodeDecode(qrCodeBytes)
		if err != nil {
			return err
		}
		// 打印二维码
		err = util.QrCodePrint(qrCode)
		if err != nil {
			return err
		}

		// 检查登录状态
		qrLoginStatus, err = c.GetQrLoginStatus(qrLoginIdResp.Data)
		if err != nil {
			return err
		}
		if qrLoginStatus.Code == 200 {
			break
		}

		// 二维码过期
		util.ClearQrCode()
		log.Error("二维码已过期, 请重新扫码登录")
	}

	log.Info("正在cas登录...")
	// cas登录
	casLoginReq := &client.CasLoginReq{
		Username:    qrLoginStatus.Data,
		Type:        "dingDingQr",
		EventID:     "submit",
		Geolocation: "",
		Execution:   execution,
	}
	result, err := c.CasLoginPost(casLoginReq)
	if err != nil {
		return err
	}
	// 判断是否登录成功
	if strings.Contains(result, "用户名密码登录") {
		return errors.New("cas使用钉钉扫码登录失败！")
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

	// 使用保存的cookies登录
	if cfg.Cookies.JSESSIONID != "" && cfg.Cookies.Route != "" && cfg.Cookies.Enabled == "1" {
		log.Info("正在使用保存的cookies登录...")
		err := c.LoadCookies(cfg)
		if err != nil {
			return nil, err
		}

		// 检查cookies是否有效
		log.Info("正在检查cookies是否有效...")
		err = c.GetClientBodyConfig()
		if err != nil && err.Error() == "可能登录过期" {
			log.Error("cookies已过期, 重新登录...")
			// 重置client
			c = client.NewClient(cfg)
		} else {
			log.Info("cookies应该maybe有效")
			return c, nil
		}
	}

	// 根据Level优先级登录
	if cfg.CasLogin.Level < cfg.NewjwLogin.Level {
		// cas登录
		log.Info("正在通过cas登录...")
		var err error
		if cfg.CasLogin.DingDingQrLoginEnabled == "1" {
			err = CasQrLogin(c, cfg)
		} else {
			err = CasPassWordLogin(c, cfg)
		}
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
			var err error
			if cfg.CasLogin.DingDingQrLoginEnabled == "1" {
				err = CasQrLogin(c, cfg)
			} else {
				err = CasPassWordLogin(c, cfg)
			}
			if err != nil {
				log.Error("cas登录失败: ", err)
				return nil, err
			}
		}
	}

	// 保存cookies
	err := c.SaveCookies(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}
