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
	CasLogin struct {
		Username               string `json:"username"`
		Password               string `json:"password"`
		DingDingQrLoginEnabled string `json:"dingDingQrLoginEnabled"`
		Level                  string `json:"level"`
	} `json:"cas_login"`
	NewjwLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Level    string `json:"level"`
	} `json:"newjw_login"`
	Cookies struct {
		JSESSIONID string `json:"JSESSIONID"`
		Route      string `json:"route"`
		Enabled    string `json:"enabled"`
	} `json:"cookies"`
	Time struct {
		XueNian string `json:"XueNian"`
		XueQi   string `json:"XueQi"`
	} `json:"time"`
	Course     *orderedmap.OrderedMap `json:"course"`
	WaitCourse struct {
		Interval int    `json:"interval"`
		Enabled  string `json:"enabled"`
	} `json:"wait_course"`
	SmtpEmail struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		To       string `json:"to"`
		Enabled  string `json:"enabled"`
	} `json:"smtp_email"`
	StartTime               string `json:"start_time"`
	ClientBodyConfigEnabled string `json:"ClientBodyConfigEnabled,omitempty"`
	DontTouchForDebug       string `json:"DontTouchForDebug,omitempty"`
}

// Course 课程信息结构体
type Course struct {
	Items []struct {
		Jxbmc  string `json:"jxbmc"`
		KchID  string `json:"kch_id"`
		JxbID  string `json:"jxb_id"`
		Jxbzc  string `json:"jxbzc"`
		Kklxmc string `json:"kklxmc"`
		Kcmc   string `json:"kcmc"` // 课程名称
		Sksj   string `json:"sksj"` // 上课时间
	} `json:"items"`
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

// ReadCourse 读取课程信息
func ReadCourse() (*Course, error) {
	// 读取课程信息
	bytes, err := os.ReadFile("course.json")
	if err != nil {
		return nil, err
	}

	// 解析课程信息
	var course Course
	if err := json.Unmarshal(bytes, &course); err != nil {
		return nil, err
	}

	return &course, nil
}

// SaveCourse 保存课程信息
func SaveCourse(course *Course) error {
	// 转换为json
	bytes, err := json.Marshal(course)
	if err != nil {
		return err
	}

	// 保存课程信息
	if err := os.WriteFile("course.json", bytes, 0666); err != nil {
		return err
	}

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
