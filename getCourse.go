package main

import (
	"HDU-KillCourse/client"
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func GetCourse(c *client.Client, cfg *config.Config) (*config.Course, error) {
	// 先从本地course.json读取课程信息
	courses, err := config.ReadCourse()
	if err == nil {
		return courses, nil
	}

	// 本地课程信息读取失败，从服务器获取课程信息
	log.Error("本地课程信息读取失败，正在从服务器获取课程信息...")
	// 初始化请求
	XueNian := cfg.Time.XueNian
	intXueNian, err := strconv.Atoi(XueNian)
	if err != nil {
		return nil, errors.New("学年格式错误")
	}
	xnmc := fmt.Sprintf("%s-%s", XueNian, strconv.Itoa(intXueNian+1))
	xqmc := cfg.Time.XueQi
	var xqm string
	if xqmc == "1" {
		xqm = "3"
	} else if xqmc == "2" {
		xqm = "12"
	} else {
		return nil, errors.New("学期格式错误")
	}
	req := &client.GetCourseReq{
		Cxfs:        "1",
		Zymc:        "全部",
		Xnmc:        xnmc,
		Xqmc:        xqmc,
		Kkxymc:      "全部",
		Jgmc:        "全部",
		Ywtk:        "0",
		Skfs:        "0",
		Xnm:         XueNian,
		Xqm:         xqm,
		Search:      "false",
		Nd:          fmt.Sprintf("%d", time.Now().Unix()),
		ShowCount:   "9999",
		CurrentPage: "1",
		SortOrder:   "asc",
		Time:        "0",
	}

	// 获取课程信息
	courseResp, err := c.GetCourse(req)
	if err != nil {
		return nil, err
	}
	// 转换为config.Course
	courses = &config.Course{
		Items: courseResp.Items,
	}

	// 保存课程信息到本地
	err = config.SaveCourse(courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}
