package course

import (
	"context"
	"errors"
	"time"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/pkg/login"
	"github.com/cr4n5/HDU-KillCourse/util"
)

// StartWaitCourse 开始蹲课
func StartWaitCourse(ctx context.Context, c *client.Client, cfg *config.Config, courses *client.GetCourseResp, CourseName string, waitCourseChannel chan string) {
	defer func() {
		waitCourseChannel <- "完成"
	}()

	firstRun := true

	for {
		if firstRun {
			firstRun = false
		} else {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(cfg.WaitCourse.Interval) * time.Second):
			}
		}

		// 检验是否有余量
		isOk, err := GetIsCourseOk(c, cfg, courses, CourseName)
		if err != nil {
			log.Error(CourseName+"查询失败: ", err)
			if err.Error() == "可能登录过期" {
				waitCourseChannel <- "登录过期"
				return
			}
			SendEmail(cfg, "查询失败", CourseName+"查询失败,将会继续蹲课: "+err.Error())
			continue
		}

		if isOk {
			// 选课
			err := HandleCourse(c, cfg, courses, CourseName, "1")
			if err != nil {
				log.Error(CourseName+"选课失败: ", err)
				if err.Error() == "可能登录过期" {
					waitCourseChannel <- "登录过期"
					return
				}
				SendEmail(cfg, "蹲选课失败", CourseName+"选课失败,将会继续蹲课: "+err.Error())
				continue
			}

			log.Info(CourseName + "蹲选课成功")
			// 发送邮件
			SendEmail(cfg, "蹲选课成功", CourseName+"选课成功？(,请自行查看确认")

			// 将此CourseName设置为0
			cfg.Course.Set(CourseName, "0")

			return
		}
	}
}

// SendEmail 发送邮件
func SendEmail(cfg *config.Config, subject string, body string) {
	if cfg.SmtpEmail.Enabled == "1" {
		err := util.SendEmail(cfg.SmtpEmail.Host, cfg.SmtpEmail.Username, cfg.SmtpEmail.Password, cfg.SmtpEmail.To, subject, body)
		if err != nil {
			log.Error("发送邮件失败: ", err)
		}
	}
}

// GetIsCourseOk 检验是否有余量
func GetIsCourseOk(c *client.Client, cfg *config.Config, course *client.GetCourseResp, CourseName string) (bool, error) {
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
				return false, errors.New("课程类型错误")
			}

			// 设置Xqm
			Xqm := cfg.Time.XueQi
			if Xqm == "1" {
				Xqm = "3"
			} else if Xqm == "2" {
				Xqm = "12"
			} else {
				return false, errors.New("学期格式错误")
			}

			// 获取课程是否有余量
			req := &client.SearchCourseReq{
				Xkxnm:      cfg.Time.XueNian,
				Xkxqm:      Xqm,
				Kklxdm:     Kklxdm,
				Kspage:     "1",
				Jspage:     "10",
				Yllist:     "1",
				Filterlist: CourseName,
				NjdmIDXs:   c.NjdmIDXs,
				ZyhIDXs:    c.ZyhIDXs,
			}

			// 发送请求
			result, err := c.SearchCourse(req)
			if err != nil {
				return false, err
			}

			// 检验是否有余量, TmpList是否长度为0
			if len(result.TmpList) == 0 {
				log.Info(CourseName + "" + v.Kcmc + ": " + "课程无余量")
				return false, nil
			} else {
				log.Info(CourseName + "" + v.Kcmc + ": " + "课程有余量")
				return true, nil
			}

		}
	}
	return false, errors.New(CourseName + "课程不存在")
}

// WaitCourse 蹲课
func WaitCourse(ctx context.Context, channel chan string, c *client.Client, cfg *config.Config, course *client.GetCourseResp) {
	defer func() {
		channel <- "完成"
	}()

	log.Info("开始蹲课...")

	// 关闭Cookies
	cfg.Cookies.Enabled = "0"

	// 获取选课配置
	err := ReadClientBodyConfig(c)
	if err != nil {
		err = c.GetClientBodyConfig()
		if err != nil {
			log.Error("获取选课配置失败: ", err)
			return
		}
	}
	log.Info("选课配置获取成功")

	for {
		waitCourseChannel := make(chan string)
		numWaitCourse := 0
		waitCourseCtx, cancel := context.WithCancel(ctx)

		// 蹲课
		for _, k := range cfg.Course.Keys() {
			v, _ := cfg.Course.Get(k)
			// 开始蹲课
			if v == "1" {
				numWaitCourse++
				go StartWaitCourse(waitCourseCtx, c, cfg, course, k, waitCourseChannel)
			}
		}

		// 等待蹲课结束
	outerLoop:
		for {
			select {
			case <-ctx.Done():
				cancel()
				return
			case message := <-waitCourseChannel:
				if message == "登录过期" {
					log.Error("登录过期")
					cancel()
					SendEmail(cfg, "登录过期", "登录过期,自动重新登录")
					close(waitCourseChannel)

					log.Info("重新登录...")
					// 重新登录
					c, err = login.Login(cfg)
					if err != nil {
						log.Error("登录失败...")
						SendEmail(cfg, "登录失败", "登录失败,程序停止,请检查")
						return
					}
					break outerLoop
				}

				numWaitCourse--
				if numWaitCourse == 0 {
					cancel()
					return
				}
			}
		}
	}
}
