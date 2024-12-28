package main

import (
	"HDU-KillCourse/client"
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 用于结束程序
	defer func() {
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}()
	ctx := context.Background()

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
	c := client.NewClient(cfg)
	err = CasLogin(c, cfg)
	if err != nil {
		log.Error("登录失败: ", err)
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
	go KillCourse(cancelCtx, c, cfg, courses)

	select {
	case <-stopChan:
		log.Info("收到终止信号，正在退出...")
		cancel()
	case <-channel:
		log.Info("退选课程已完成，正在退出...")
		cancel()
	}
}
