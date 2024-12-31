package main

import (
	"HDU-KillCourse/client"
	"HDU-KillCourse/config"
	"HDU-KillCourse/log"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"
)

// killcourse完成
var channel = make(chan string)

// GetDoJxbId 获取doJxbId
func GetDoJxbId(c *client.Client, KchId string, JxbId string, Kklxdm string, NjdmId string, XueNian string, Xqm string) (string, error) {
	// 检查c.ClientBodyConfig是否为nil，测试用
	if c.ClientBodyConfig == nil {
		return "", errors.New("ClientBodyConfig未初始化")
	}

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

	if result.Flag == "1" {
		log.Info("选课成功")
	} else if result.Flag == "0" {
		log.Error("选课失败: ", result.Msg)
	} else {
		log.Error("选课失败: 人数可能已满", result)
	}

	return nil
}

// CancelCourse 退课
func CancelCourse(c *client.Client, JxbIds string, KchId string, XueNian string, Xqm string) error {
	// 设置请求参数
	req := &client.CancelCourseReq{
		JxbIDs: JxbIds,
		KchID:  KchId,
		Xkxnm:  XueNian,
		Xkxqm:  Xqm,
	}

	// 发送请求
	result, err := c.CancelCourse(req)
	if err != nil {
		return err
	}

	if result == "\"1\"" {
		log.Info("退课成功(可能？)")
	} else {
		log.Error("退课失败：", result)
	}

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
			NjdmId := "20" + c.ClientBodyConfig.BhId[0:2]

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
				err = CancelCourse(c, doJxbId, v.KchID, cfg.Time.XueNian, Xqm)
				if err != nil {
					return err
				}
			}

			return nil
		}
	}
	return errors.New(CourseName + "课程不存在")
}

// KillCourse 选退课
func KillCourse(ctx context.Context, c *client.Client, cfg *config.Config, course *config.Course) {
	// 计算需要等待的时间
	// 时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Error("初始化时间地区失败，正在使用手动定义的时区信息 :", err)
		loc = time.FixedZone("CST", 8*3600)
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", cfg.StartTime, loc)
	if err != nil {
		log.Error("时间格式错误: ", err)
		return
	}
	log.Info("选课开始时间: ", t)
	waitTime := t.Unix() - time.Now().Unix()

	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Duration(waitTime) * time.Second):
		log.Info("时间已到，开始处理课程...")
		//go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 获取选课配置
				err = ReadClientBodyConfig(c)
				if err != nil {
					err = c.GetClientBodyConfig()
					if err != nil {
						log.Error("获取选课配置失败: ", err)
						continue
					}
				}
				// 保存选课配置
				if cfg.ClientBodyConfigEnabled == "1" {
					err = SaveClientBodyConfig(c)
					if err != nil {
						log.Error("保存选课配置失败: ", err)
						return
					}
				}
				log.Info("选课配置获取成功")
				// 选退课
				for k, v := range cfg.Course {
					// 处理课程
					log.Info("正在处理课程: ", k)
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
		//}()
	}
}

// SaveClientBodyConfig 保存选课配置
func SaveClientBodyConfig(c *client.Client) error {
	// 将c.ClientBodyConfig保存到文件CLientBodyConfig.json
	clientBodyConfig := c.ClientBodyConfig
	bytes, err := json.Marshal(clientBodyConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile("ClientBodyConfig.json", bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

// ReadClientBodyConfig 读取选课配置
func ReadClientBodyConfig(c *client.Client) error {
	// 读取文件CLientBodyConfig.json到c.ClientBodyConfig
	bytes, err := os.ReadFile("ClientBodyConfig.json")
	if err != nil {
		return err
	}

	var clientBodyConfig client.ClientBodyConfig
	err = json.Unmarshal(bytes, &clientBodyConfig)
	if err != nil {
		return err
	}

	c.ClientBodyConfig = &clientBodyConfig

	return nil
}
