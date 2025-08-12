package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/iancoleman/orderedmap"
)

//go:embed main.html
var htmlFile embed.FS

var cfg *config.Config
var webCfg WebConfig

// 将 OrderedMap 转化为二维数组列表
func orderedMapToArrayList(omap *orderedmap.OrderedMap) [][]string {
	var result [][]string
	for _, key := range omap.Keys() {
		value, _ := omap.Get(key)
		result = append(result, []string{key, value.(string)})
	}
	return result
}

// 将二维数组列表转化为 OrderedMap
func arrayListToOrderedMap(arrayList [][]string) *orderedmap.OrderedMap {
	omap := orderedmap.New()
	for _, item := range arrayList {
		if len(item) == 2 {
			omap.Set(item[0], item[1])
		} else {
			log.Error("二维数组列表中的项不是长度为2的数组，跳过: ", item)
		}
	}
	return omap
}

type WebConfig struct {
	CasLogin   config.CasLogin   `json:"cas_login"`
	NewjwLogin config.NewjwLogin `json:"newjw_login"`
	Cookies    config.Cookies    `json:"cookies"`
	Time       config.Time       `json:"time"`
	Course     [][]string        `json:"course"` // 使用二维数组列表来存储课程信息
	WaitCourse config.WaitCourse `json:"wait_course"`
	SmtpEmail  config.SmtpEmail  `json:"smtp_email"`
	StartTime  string            `json:"start_time"`
}

func StartWebServer() {
	// 初始端口号
	port := 6688
	// 检查端口是否被占用
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			// 如果端口未被占用，关闭监听器并继续启动服务器
			listener.Close()
			break
		}
		port++
	}

	// 设置路由
	// 提供 HTML 页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := htmlFile.ReadFile("main.html")
		if err != nil {
			http.Error(w, "无法加载页面", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	// 获取配置
	http.HandleFunc("/getConfig", func(w http.ResponseWriter, r *http.Request) {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			http.Error(w, "无法读取配置文件: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// 将 OrderedMap 转化为二维数组列表
		courseArray := orderedMapToArrayList(cfg.Course)
		// 创建 WebConfig 对象
		webCfg = WebConfig{
			CasLogin:   cfg.CasLogin,
			NewjwLogin: cfg.NewjwLogin,
			Cookies:    cfg.Cookies,
			Time:       cfg.Time,
			Course:     courseArray,
			WaitCourse: cfg.WaitCourse,
			SmtpEmail:  cfg.SmtpEmail,
			StartTime:  cfg.StartTime,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webCfg)
	})

	// 保存配置
	http.HandleFunc("/saveConfig", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&webCfg)
		if err != nil {
			http.Error(w, "无法解析配置: "+err.Error(), http.StatusBadRequest)
			return
		}
		// 将二维数组列表转化为 OrderedMap
		courseMap := arrayListToOrderedMap(webCfg.Course)
		// 更新 cfg
		cfg.CasLogin = webCfg.CasLogin
		cfg.NewjwLogin = webCfg.NewjwLogin
		cfg.Cookies = webCfg.Cookies
		cfg.Time = webCfg.Time
		cfg.Course = courseMap
		cfg.WaitCourse = webCfg.WaitCourse
		cfg.SmtpEmail = webCfg.SmtpEmail
		cfg.StartTime = webCfg.StartTime
		// 验证配置
		err = cfg.Validate()
		if err != nil {
			http.Error(w, "配置验证失败: "+err.Error(), http.StatusBadRequest)
			return
		}
		// 保存配置
		err = config.SaveConfig(cfg)
		if err != nil {
			http.Error(w, "无法保存配置文件: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("配置保存成功"))
	})

	log.Info("访问该地址编辑配置: http://localhost:" + fmt.Sprintf("%d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Error("Web服务器启动失败: ", err)
		return
	}
}
