package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
)

//go:embed main.html
var htmlFile embed.FS

func StartWebServer() error {
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
		cfg, err := config.LoadConfig()
		if err != nil {
			http.Error(w, "无法读取配置文件: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg)
	})

	// 保存配置
	http.HandleFunc("/saveConfig", func(w http.ResponseWriter, r *http.Request) {
		var cfg config.Config
		err := json.NewDecoder(r.Body).Decode(&cfg)
		if err != nil {
			http.Error(w, "无法解析配置: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = config.SaveConfig(&cfg)
		if err != nil {
			http.Error(w, "无法保存配置文件: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("配置保存成功"))
	})

	log.Info("访问该地址编辑配置: http://localhost:" + fmt.Sprintf("%d", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
