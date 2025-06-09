package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/iancoleman/orderedmap"
)

// Config 配置文件结构体
type Config struct {
	CasLogin                `json:"cas_login"`
	NewjwLogin              `json:"newjw_login"`
	Cookies                 `json:"cookies"`
	Time                    `json:"time"`
	Course                  *orderedmap.OrderedMap `json:"course"`
	WaitCourse              `json:"wait_course"`
	SmtpEmail               `json:"smtp_email"`
	StartTime               string `json:"start_time"`
	ClientBodyConfigEnabled string `json:"ClientBodyConfigEnabled,omitempty"`
	CrossGradeEnabled       string `json:"CrossGradeEnabled,omitempty"`
}

// CasLogin CAS 登录配置
type CasLogin struct {
	Username               string `json:"username"`
	Password               string `json:"password"`
	DingDingQrLoginEnabled string `json:"dingDingQrLoginEnabled"`
	Level                  string `json:"level"`
}

// NewjwLogin 新教务登录配置
type NewjwLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Level    string `json:"level"`
}

// Cookies Cookie 配置
type Cookies struct {
	JSESSIONID string `json:"JSESSIONID"`
	Route      string `json:"route"`
	Enabled    string `json:"enabled"`
}

// Time 学年学期配置
type Time struct {
	XueNian string `json:"XueNian"`
	XueQi   string `json:"XueQi"`
}

// WaitCourse 蹲课配置
type WaitCourse struct {
	Interval int    `json:"interval"`
	Enabled  string `json:"enabled"`
}

// SmtpEmail SMTP 邮件配置
type SmtpEmail struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	To       string `json:"to"`
	Enabled  string `json:"enabled"`
}

func InitCfg() (*Config, error) {
	// 读取配置文件
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	// 解析配置文件
	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	// 验证配置文件
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// 默认配置文件
var DefaultConfig = `{
    "cas_login": {
        "username": "2201xxxx",
        "password": "xxxxxxxx",
        "dingDingQrLoginEnabled": "0",
        "level": "0"
    },
    "newjw_login": {
        "username": "2201xxxx",
        "password": "xxxxxxxx",
        "level": "1"
    },
    "cookies": {
        "JSESSIONID": "",
        "route": "",
        "enabled": "1"
    },
    "time": {
        "XueNian": "2024",
        "XueQi": "1"
    },
    "course": {
        "(2024-2025-1)-C2092011-01" : "1"
    },
    "wait_course": {
        "interval": 60,
        "enabled": "0"
    },
    "smtp_email": {
        "host": "smtp.qq.com",
        "username": "...@qq.com",
        "password": "xxxxxxxx",
        "to": "...@qq.com",
        "enabled": "0"
    },
    "start_time": "2024-07-25 12:00:00"
}`

// LoadConfig 加载配置文件  用于在线编辑配置
func LoadConfig() (*Config, error) {
	// 读取配置文件
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Error("读取配置文件失败: ", err)
		log.Info("使用默认配置文件")
		// 如果读取失败，则使用默认配置文件
		bytes = []byte(DefaultConfig)
	}
	// 解析配置文件
	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		log.Error("解析配置文件失败: ", err)
		return nil, err
	}
	return &cfg, nil
}

// Validate 验证配置文件
func (cfg *Config) Validate() error {
	if (cfg.CasLogin.Username == "" || cfg.CasLogin.Password == "") && (cfg.NewjwLogin.Username == "" || cfg.NewjwLogin.Password == "") {
		return errors.New("用户名或密码为空")
	}
	if cfg.Time.XueNian == "" || cfg.Time.XueQi == "" {
		return errors.New("学年或学期为空")
	}
	if cfg.Course == nil {
		return errors.New("课程为空")
	}
	if cfg.WaitCourse.Interval == 0 || cfg.WaitCourse.Enabled == "" {
		return errors.New("WaitCourse为空")
	}
	if cfg.SmtpEmail.Enabled == "1" {
		if cfg.SmtpEmail.Host == "" || cfg.SmtpEmail.Username == "" || cfg.SmtpEmail.Password == "" || cfg.SmtpEmail.To == "" {
			return errors.New("SmtpEmail为空")
		}
	}
	if cfg.StartTime == "" {
		return errors.New("StartTime为空")
	}

	// 打印配置文件
	// 空行
	log.Info("")

	log.Info(log.InfoColor("CasLogin:"))
	log.Info("  Username: ", cfg.CasLogin.Username)
	log.Info("  Password: ", cfg.CasLogin.Password)
	log.Info("  DingDingQrLoginEnabled: ", cfg.CasLogin.DingDingQrLoginEnabled)
	log.Info("  Level: ", cfg.CasLogin.Level)
	log.Info(log.InfoColor("NewjwLogin:"))
	log.Info("  Username: ", cfg.NewjwLogin.Username)
	log.Info("  Password: ", cfg.NewjwLogin.Password)
	log.Info("  Level: ", cfg.NewjwLogin.Level)
	log.Info(log.InfoColor("XueNian: "), cfg.Time.XueNian)
	log.Info(log.InfoColor("XueQi: "), cfg.Time.XueQi)
	log.Info(log.InfoColor("WaitCourse:"))
	log.Info("  Interval: ", cfg.WaitCourse.Interval)
	log.Info("  Enabled: ", cfg.WaitCourse.Enabled)
	log.Info(log.InfoColor("SmtpEmail:"))
	if cfg.SmtpEmail.Enabled == "1" {
		log.Info("  Host: ", cfg.SmtpEmail.Host)
		log.Info("  Username: ", cfg.SmtpEmail.Username)
		log.Info("  Password: ", cfg.SmtpEmail.Password)
		log.Info("  To: ", cfg.SmtpEmail.To)
	} else {
		log.Info("  SmtpEmailEnabled: ", cfg.SmtpEmail.Enabled)
	}
	log.Info(log.InfoColor("StartTime: "), cfg.StartTime)
	log.Info(log.InfoColor("Course:"))
	for _, k := range cfg.Course.Keys() {
		v, _ := cfg.Course.Get(k)
		log.Info(k, ": ", v)
	}

	// 空行
	log.Info("")

	return nil
}

// SaveConfig 保存配置文件
func SaveConfig(cfg *Config) error {
	// 转换为json
	bytes, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	// 保存配置文件
	if err := os.WriteFile("config.json", bytes, 0666); err != nil {
		return err
	}

	return nil
}
