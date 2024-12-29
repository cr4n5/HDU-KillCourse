package config

import (
	"HDU-KillCourse/log"
	"encoding/json"
	"errors"
	"os"
)

// Config 配置文件结构体
type Config struct {
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
	Time struct {
		XueNian string `json:"XueNian"`
		XueQi   string `json:"XueQi"`
	} `json:"time"`
	Course                  map[string]string `json:"course"`
	StartTime               string            `json:"start_time"`
	ClientBodyConfigEnabled string            `json:"ClientBodyConfigEnabled"`
}

// Course 课程信息结构体
type Course struct {
	Items []struct {
		Jxbmc  string `json:"jxbmc"`
		KchID  string `json:"kch_id"`
		JxbID  string `json:"jxb_id"`
		Jxbzc  string `json:"jxbzc"`
		Kklxmc string `json:"kklxmc"`
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
	if cfg.Login.Username == "" || cfg.Login.Password == "" {
		return errors.New("用户名或密码为空")
	}
	if cfg.Time.XueNian == "" || cfg.Time.XueQi == "" {
		return errors.New("学年或学期为空")
	}
	if len(cfg.Course) == 0 {
		return errors.New("课程为空")
	}
	if cfg.StartTime == "" {
		return errors.New("StartTime为空")
	}

	// 打印配置文件
	log.Info("Username: ", cfg.Login.Username)
	log.Info("Password: ", cfg.Login.Password)
	log.Info("XueNian: ", cfg.Time.XueNian)
	log.Info("XueQi: ", cfg.Time.XueQi)
	log.Info("StartTime: ", cfg.StartTime)
	log.Info("Course:")
	for k, v := range cfg.Course {
		log.Info(k, ": ", v)
	}

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
