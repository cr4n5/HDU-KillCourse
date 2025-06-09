package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
)

func GetCourse(c *client.Client, cfg *config.Config) (*client.GetCourseResp, error) {
	// 先从本地course.json读取课程信息
	courses, err := ReadCourse(cfg)
	if err == nil {
		return courses, nil
	}

	// 本地课程信息读取失败，从服务器获取课程信息
	log.Error("本地课程信息读取失败，正在从服务器获取课程信息...")
	log.Info("Notice！: 等待时间可能较长，请耐心等待...")
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

	// 保存课程信息到本地
	err = SaveCourse(courseResp)
	if err != nil {
		return nil, err
	}

	return courseResp, nil
}

// VarifyCourse 验证课程信息
func VarifyCourse(course *client.GetCourseResp, cfg *config.Config) error {
	// 检查课程信息是否为空
	if course == nil || len(course.Items) == 0 {
		return errors.New("course.json课程信息为空，请检查配置文件或网络连接")
	}

	// 检查学年学期是否匹配
	xn := cfg.Time.XueNian
	xq := cfg.Time.XueQi
	// 只验证第一个课程的学年学期
	if len(course.Items) > 0 {
		Jxbmc := course.Items[0].Jxbmc
		if len(Jxbmc) < 12 {
			return errors.New("课程信息格式错误，Jxbmc长度不足")
		}
		if Jxbmc[1:5] != xn || Jxbmc[11:12] != xq {
			return errors.New("course.json课程信息学年学期与配置不匹配")
		}
	}

	return nil
}

// ReadCourse 读取课程信息
func ReadCourse(cfg *config.Config) (*client.GetCourseResp, error) {
	// 读取课程信息
	bytes, err := os.ReadFile("course.json")
	if err != nil {
		return nil, err
	}

	// 解析课程信息
	var course client.GetCourseResp
	if err := json.Unmarshal(bytes, &course); err != nil {
		return nil, err
	}

	// 验证课程信息
	if err := VarifyCourse(&course, cfg); err != nil {
		log.Error("课程信息验证失败: ", err)
		return nil, err
	}

	return &course, nil
}

// SaveCourse 保存课程信息
func SaveCourse(course *client.GetCourseResp) error {
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
