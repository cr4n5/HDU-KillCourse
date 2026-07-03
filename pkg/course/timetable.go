package course

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
)

// xqmCode 由学期(1/2)得到教务系统学期码
func xqmCode(xueQi string) (string, error) {
	switch xueQi {
	case "1":
		return "3", nil
	case "2":
		return "12", nil
	default:
		return "", errors.New("学期格式错误")
	}
}

// parseWeekday 解析星期：优先数字xqj，其次名称xqjmc(如"星期二")
func parseWeekday(xqj, xqjmc string) int {
	if n, err := strconv.Atoi(strings.TrimSpace(xqj)); err == nil && n >= 1 && n <= 7 {
		return n
	}
	xqjmc = strings.TrimPrefix(strings.TrimSpace(xqjmc), "星期")
	if len(xqjmc) > 0 {
		if d, ok := dayMap[string([]rune(xqjmc)[0])]; ok {
			return d
		}
	}
	return 0
}

// parsePeriods 解析节次，如"3-5"、"3"
func parsePeriods(jcs string) (int, int) {
	jcs = strings.TrimSpace(jcs)
	if jcs == "" {
		return 0, 0
	}
	if i := strings.Index(jcs, "-"); i >= 0 {
		p1, _ := strconv.Atoi(strings.TrimSpace(jcs[:i]))
		p2, _ := strconv.Atoi(strings.TrimSpace(jcs[i+1:]))
		return p1, p2
	}
	p, _ := strconv.Atoi(jcs)
	return p, p
}

// scheduleItemsToEvents 将个人课表项转换为单次上课事件
func scheduleItemsToEvents(items []client.ScheduleItem, cfg *config.Config) ([]sessionEvent, error) {
	weekOneMonday, loc, err := semesterMonday(cfg)
	if err != nil {
		return nil, err
	}

	var events []sessionEvent
	for _, it := range items {
		day := parseWeekday(it.Xqj, it.Xqjmc)
		if day == 0 {
			log.Error("生成课表ICS: 无法解析星期，跳过 ", it.Kcmc, " ", it.Xqjmc, it.Xqj)
			continue
		}
		jcs := it.Jcs
		if jcs == "" {
			jcs = it.Jc
		}
		p1, p2 := parsePeriods(jcs)
		p1t, ok1 := PeriodTimes[p1]
		p2t, ok2 := PeriodTimes[p2]
		if !ok1 || !ok2 {
			log.Error("生成课表ICS: 未知节次，跳过 ", it.Kcmc, " ", jcs)
			continue
		}
		weeks := parseWeeks(it.Zcd)
		if len(weeks) == 0 {
			log.Error("生成课表ICS: 无法解析周次，跳过 ", it.Kcmc, " ", it.Zcd)
			continue
		}
		code := it.Kch
		if code == "" {
			code = CourseCodeOfJxbmc(it.Jxbmc)
		}
		for w := 1; w <= 30; w++ {
			if !weeks[w] {
				continue
			}
			d := weekOneMonday.AddDate(0, 0, (w-1)*7+day-1)
			start, _ := time.ParseInLocation("2006-01-02 15:04", d.Format("2006-01-02")+" "+p1t[0], loc)
			end, _ := time.ParseInLocation("2006-01-02 15:04", d.Format("2006-01-02")+" "+p2t[1], loc)
			desc := fmt.Sprintf("教师: %s\n第%d周 %s 第%d-%d节", it.Xm, w, weekdayCN[day], p1, p2)
			if it.Jxbmc != "" {
				desc = "教学班: " + it.Jxbmc + "\n" + desc
			}
			events = append(events, sessionEvent{
				uid:      fmt.Sprintf("%s-w%d-d%d-p%d@hdu-killcourse-kb", code, w, day, p1),
				start:    start,
				end:      end,
				summary:  it.Kcmc,
				location: it.Cdmc,
				desc:     desc,
			})
		}
	}
	return events, nil
}

// ExportPersonalTimetableICS 拉取个人完整课表(含教务处预置分配的课程)并生成ICS日历
func ExportPersonalTimetableICS(c *client.Client, cfg *config.Config) error {
	if cfg.SemesterStartDate == "" {
		return errors.New("未设置semester_start_date(第1周星期一日期)，跳过个人课表ICS生成")
	}
	xqm, err := xqmCode(cfg.Time.XueQi)
	if err != nil {
		return err
	}

	schedule, err := c.GetPersonalSchedule(cfg.Time.XueNian, xqm)
	if err != nil {
		return err
	}
	items := append([]client.ScheduleItem{}, schedule.KbList...)
	items = append(items, schedule.SjkList...)
	if len(items) == 0 {
		return errors.New("个人课表为空(可能学年学期无课或登录态失效)")
	}

	events, err := scheduleItemsToEvents(items, cfg)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return errors.New("个人课表未生成任何日历事件")
	}

	content := renderICS(events, "个人课表 "+XnXqName(cfg))
	fileName := fmt.Sprintf("timetable_%s_full.ics", XnXqName(cfg))
	if err := os.WriteFile(fileName, []byte(content), 0666); err != nil {
		return err
	}
	log.Info("个人完整课表ICS已生成: ", fileName, " (", len(events), " 个上课事件)")
	return nil
}
