package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/vars"
	"github.com/cr4n5/HDU-KillCourse/web"
)

// 程序结束信号
var channel = make(chan string)

func main() {
	// 用于结束程序
	defer func() {
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}()
	ctx := context.Background()

	vars.ShowPortal()

	// 5秒倒计时，用户可以选择是否进入 Web 编辑配置模式
	log.Info("按下 Enter 跳过进入 Web 编辑配置模式，5秒后自动进入...")
	enterChan := make(chan bool)
	go func() {
		fmt.Scanln()
		enterChan <- true
	}()

	select {
	case <-enterChan:
		log.Info("跳过 Web 编辑配置模式...")
	case <-time.After(5 * time.Second):
		log.Info("进入 Web 编辑配置模式...")
		err := web.StartWebServer()
		if err != nil {
			log.Error("Web服务器启动失败: ", err)
			return
		}
		return
	}

	// 读取配置文件
	log.Info("开始读取配置文件...")
	cfg, err := config.InitCfg()
	if err != nil {
		log.Error("读取配置文件失败: ", err)
		return
	}
	log.Info("读取配置文件成功")

	// 登录
	log.Info("开始登录...")
	c, err := Login(cfg)
	if err != nil {
		log.Error("登录失败...")
		return
	}
	log.Info("登录成功...")

	// 获取课程信息
	log.Info("开始获取课程信息...")
	courses, err := GetCourse(c, cfg)
	if err != nil {
		log.Error("获取课程信息失败: ", err)
		return
	}
	log.Info("获取课程信息成功...")
	log.Info(log.ErrorColor("Notice！: 在下学期选课开始前，请删除course.json文件，获取最新课程信息"))

	cancelCtx, cancel := context.WithCancel(ctx)
	// 捕获终止信号
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 选退课
	if cfg.WaitCourse.Enabled == "1" {
		go WaitCourse(cancelCtx, c, cfg, courses)
	} else {
		go KillCourse(cancelCtx, c, cfg, courses)
	}

	select {
	case <-stopChan:
		log.Info("收到终止信号，正在退出...")
		cancel()
	case <-channel:
		log.Info("此程序已完成，正在退出...")
		cancel()
	}
}
