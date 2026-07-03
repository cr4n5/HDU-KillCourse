package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/pkg/course"
	"github.com/cr4n5/HDU-KillCourse/pkg/login"
	"github.com/cr4n5/HDU-KillCourse/pkg/version"
	"github.com/cr4n5/HDU-KillCourse/pkg/web"
	"github.com/cr4n5/HDU-KillCourse/vars"
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

	// 检查版本更新
	go version.VersionUpdate()

	// 启动web服务器编辑配置
	go web.StartWebServer()
	log.Info("按下Enter继续，或在web界面编辑配置文件...")
	fmt.Scanln()

	// 读取配置文件
	log.Info("开始读取配置文件...")
	cfg, err := config.InitCfg()
	if err != nil {
		log.Error("读取配置文件失败: ", err)
		return
	}
	log.Info("读取配置文件成功")
	log.Info(log.ErrorColor("Notice！: 读取配置文件后，如果配置文件有修改，请重启程序"))

	// 登录
	log.Info("开始登录...")
	c, err := login.Login(cfg)
	if err != nil {
		log.Error("登录失败...")
		return
	}
	log.Info("登录成功...")

	// 获取课程信息
	log.Info("开始获取课程信息...")
	courses, err := course.GetCourse(c, cfg)
	if err != nil {
		log.Error("获取课程信息失败: ", err)
		return
	}
	log.Info("获取课程信息成功...")

	// 课程自动排序：退课先于与其冲突的选课执行（course_sort_enabled置1开启，默认保持手动排序）
	if cfg.CourseSortEnabled == "1" {
		log.Info("开始课程自动排序...")
		if err := course.AutoSortCourses(cfg, courses); err != nil {
			log.Error("课程自动排序失败: ", err)
		}
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	// 捕获终止信号
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 选退课
	if cfg.WaitCourse.Enabled == "1" {
		go course.WaitCourse(cancelCtx, channel, c, cfg, courses)
	} else {
		go course.KillCourse(cancelCtx, channel, c, cfg, courses)
	}

	select {
	case <-stopChan:
		log.Info("收到终止信号，正在退出...")
		cancel()
	case <-channel:
		log.Info("此程序已完成，正在退出...")
		cancel()
		// 抢课完成后依据最新课表生成ICS日历：反映实际选课结果，
		// 自动排除未选上的备选教学班，并包含教务处预置分配的课程。
		if err := course.ExportTimetableICS(c, cfg); err != nil {
			log.Error("生成课表ICS失败: ", err)
		}
	}
}
