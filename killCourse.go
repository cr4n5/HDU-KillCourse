package main

import (
	"HDU-KillCourse/client"
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"context"
	"errors"
	"time"
)

// killcourse完成
var channel = make(chan string)

// GetDoJxbId 获取doJxbId
func GetDoJxbId(c *client.Client, KchId string, JxbId string, Kklxdm string, NjdmId string, XueNian string, Xqm string) (string, error) {
	// 设置请求参数
	req := &client.GetDoJxbIdReq{
		BklxID: "0",
		NjdmID: NjdmId,
		Xkxnm:  XueNian,
		Xkxqm:  Xqm,
		Kklxdm: Kklxdm,
		KchID:  KchId,
		XkkzID: c.ClientBodyConfig.XkkzId[Kklxdm],
		Xsbj:   c.ClientBodyConfig.Xsbj,
		Ccdm:   c.ClientBodyConfig.Ccdm,
		Xz:     c.ClientBodyConfig.Xz,
		Mzm:    c.ClientBodyConfig.Mzm,
		Xslbdm: c.ClientBodyConfig.Xslbdm,
		Xbm:    c.ClientBodyConfig.Xbm,
		BhID:   c.ClientBodyConfig.BhId,
		ZyfxID: c.ClientBodyConfig.ZyfxId,
		JgID:   c.ClientBodyConfig.JgId,
		XqhID:  c.ClientBodyConfig.XqhId,
	}

	// 发送请求
	resp, err := c.GetDoJxbId(req)
	if err != nil {
		return "", err
	}

	// 解析doJxbId
	for _, v := range resp {
		if v.JxbID == JxbId {
			return v.DoJxbID, nil
		}
	}

	return "", errors.New("doJxbId不存在")

}

// SelectCourse 选课
func SelectCourse(c *client.Client, JxbIds string, KchId string, Kklxdm string, Jxbzc string) error {
	// 设置请求参数
	req := &client.SelectCourseReq{
		JxbIDs: JxbIds,
		KchID:  KchId,
		Qz:     "0",
	}

	// 若为主修课程
	if Kklxdm == "01" {
		req.NjdmID = "20" + Jxbzc[0:2]
		req.ZyhID = Jxbzc[2:6]
	}

	// 发送请求
	result, err := c.SelectCourse(req)
	if err != nil {
		return err
	}

	log.Info(result) // 打印选课结果，待完善
	return nil
}

// CancelCourse 退课
func CancelCourse(c *client.Client, JxbIds string, KchId string) error {
	// 设置请求参数
	req := &client.CancelCourseReq{
		JxbIDs: JxbIds,
		KchID:  KchId,
	}

	// 发送请求
	result, err := c.CancelCourse(req)
	if err != nil {
		return err
	}

	log.Info(result) // 打印退课结果，待完善
	return nil
}

// HandleCourse 处理课程
func HandleCourse(c *client.Client, cfg *config.Config, course *config.Course, CourseName string, SelectFlag string) error {
	for _, v := range course.Items {
		if v.Jxbmc == CourseName {
			// 更改Kklxdm
			Kklxdm := v.Kklxmc
			if Kklxdm == "主修课程" {
				Kklxdm = "01"
			} else if Kklxdm == "通识选修课" {
				Kklxdm = "10"
			} else if Kklxdm == "体育分项" {
				Kklxdm = "05"
			} else if Kklxdm == "特殊课程" {
				Kklxdm = "09"
			} else {
				return errors.New("课程类型错误")
			}

			// 设置NjdmId
			NjdmId := "20" + cfg.Login.Username[0:2]

			// 设置Xqm
			Xqm := cfg.Time.XueQi
			if Xqm == "1" {
				Xqm = "3"
			} else if Xqm == "2" {
				Xqm = "12"
			} else {
				return errors.New("学期格式错误")
			}

			// 获取doJxbId
			doJxbId, err := GetDoJxbId(c, v.KchID, v.JxbID, Kklxdm, NjdmId, cfg.Time.XueNian, Xqm)
			if err != nil {
				return err
			}

			// 选课
			if SelectFlag == "1" {
				err = SelectCourse(c, doJxbId, v.KchID, Kklxdm, v.Jxbzc)
				if err != nil {
					return err
				}
			} else {
				// 退课
				err = CancelCourse(c, doJxbId, v.KchID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// KillCourse 选退课
func KillCourse(ctx context.Context, c *client.Client, cfg *config.Config, course *config.Course) {
	// 计算需要等待的时间
	t, err := time.Parse("2006-01-02 15:04:05", cfg.StartTime)
	if err != nil {
		log.Error("时间格式错误: ", err)
		return
	}
	log.Info("选课开始时间: ", t)
	waitTime := t.Unix() - time.Now().Unix()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(waitTime) * time.Second):
			// 获取选课配置
			err = c.GetClientBodyConfig()
			if err != nil {
				log.Error("获取选课配置失败: ", err)
				continue
			}
			// 选退课
			for k, v := range cfg.Course {
				// 处理课程
				err = HandleCourse(c, cfg, course, k, v)
				if err != nil {
					log.Error("处理课程失败: ", err)
					continue
				}
			}
			// 完成
			channel <- "完成"
			return
		}
	}
}
