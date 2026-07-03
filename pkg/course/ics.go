package course

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
)

// PeriodTimes 节次作息时间表（下沙校区常用作息，如有出入请自行调整）
var PeriodTimes = map[int][2]string{
	1:  {"08:05", "08:50"},
	2:  {"08:55", "09:40"},
	3:  {"09:55", "10:40"},
	4:  {"10:45", "11:30"},
	5:  {"11:35", "12:20"},
	6:  {"13:30", "14:15"},
	7:  {"14:20", "15:05"},
	8:  {"15:15", "16:00"},
	9:  {"16:05", "16:50"},
	10: {"18:30", "19:15"},
	11: {"19:20", "20:05"},
	12: {"20:10", "20:55"},
}

var weekdayCN = []string{"", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日"}

// icsEscape 转义ICS文本字段
func icsEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, ";", `\;`)
	s = strings.ReplaceAll(s, ",", `\,`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}

// foldLine 按RFC5545折叠长行（不超过约74字节，按utf8边界断行）
func foldLine(line string) string {
	const limit = 74
	if len(line) <= limit {
		return line
	}
	var b strings.Builder
	cur := 0
	for _, r := range line {
		rl := len(string(r))
		if cur+rl > limit {
			b.WriteString("\r\n ")
			cur = 1
		}
		b.WriteString(string(r))
		cur += rl
	}
	return b.String()
}

func icsLine(b *strings.Builder, line string) {
	b.WriteString(foldLine(line))
	b.WriteString("\r\n")
}

// sessionEvent 一次课的日历事件
type sessionEvent struct {
	uid      string
	start    time.Time
	end      time.Time
	summary  string
	location string
	desc     string
}

// semesterMonday 解析配置中的第1周星期一日期与时区
func semesterMonday(cfg *config.Config) (time.Time, *time.Location, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	m, err := time.ParseInLocation("2006-01-02", cfg.SemesterStartDate, loc)
	if err != nil {
		return time.Time{}, loc, errors.New("semester_start_date格式错误，应为YYYY-MM-DD(第1周星期一): " + cfg.SemesterStartDate)
	}
	if m.Weekday() != time.Monday {
		log.Info("Notice！: semester_start_date不是星期一，请确认其为第1周星期一的日期")
	}
	return m, loc, nil
}

// buildEvents 由配置中的选课课程生成全部单次上课事件
func buildEvents(cfg *config.Config, idx CourseIndex) ([]sessionEvent, error) {
	weekOneMonday, loc, err := semesterMonday(cfg)
	if err != nil {
		return nil, err
	}

	var events []sessionEvent
	seenCode := make(map[string]string) // 课程号 -> 已生成日历的教学班

	for _, k := range cfg.Course.Keys() {
		v, _ := cfg.Course.Get(k)
		if v != "1" {
			continue // 只为选课的课程生成日历
		}
		code := CourseCodeOfJxbmc(k)
		if first, dup := seenCode[code]; dup {
			log.Info("生成ICS: 跳过同一课程的备选教学班 ", k, " (已使用 ", first, ")")
			continue
		}
		seenCode[code] = k

		info, ok := idx[k]
		if !ok {
			log.Error("生成ICS: 未找到课程信息，跳过 ", k)
			continue
		}

		// 一个上课时间段可能对应多个地点（以;分隔，与时间段一一对应）
		locParts := strings.Split(info.Jxdd, ";")
		slots := ParseSksj(info.Sksj)
		for si, slot := range slots {
			slotLoc := info.Jxdd
			if len(locParts) > si && len(slots) > 1 {
				slotLoc = locParts[si]
			}
			p1t, ok1 := PeriodTimes[slot.P1]
			p2t, ok2 := PeriodTimes[slot.P2]
			if !ok1 || !ok2 {
				log.Error("生成ICS: 未知节次，跳过 ", k, " 第", slot.P1, "-", slot.P2, "节")
				continue
			}
			for w := 1; w <= 30; w++ {
				if !slot.Weeks[w] {
					continue
				}
				day := weekOneMonday.AddDate(0, 0, (w-1)*7+slot.Day-1)
				start, _ := time.ParseInLocation("2006-01-02 15:04", day.Format("2006-01-02")+" "+p1t[0], loc)
				end, _ := time.ParseInLocation("2006-01-02 15:04", day.Format("2006-01-02")+" "+p2t[1], loc)

				desc := fmt.Sprintf("教学班: %s\n教师: %s\n第%d周 %s 第%d-%d节",
					k, info.Jzgxx, w, weekdayCN[slot.Day], slot.P1, slot.P2)

				events = append(events, sessionEvent{
					uid:      fmt.Sprintf("%s-w%d-d%d-p%d@hdu-killcourse", code, w, slot.Day, slot.P1),
					start:    start,
					end:      end,
					summary:  info.Kcmc,
					location: slotLoc,
					desc:     desc,
				})
			}
		}
	}
	return events, nil
}

// renderICS 渲染ICS文件内容
func renderICS(events []sessionEvent, calName string) string {
	var b strings.Builder
	icsLine(&b, "BEGIN:VCALENDAR")
	icsLine(&b, "VERSION:2.0")
	icsLine(&b, "PRODID:-//HDU-KillCourse//Timetable//CN")
	icsLine(&b, "CALSCALE:GREGORIAN")
	icsLine(&b, "METHOD:PUBLISH")
	icsLine(&b, "X-WR-CALNAME:"+icsEscape(calName))
	icsLine(&b, "X-WR-TIMEZONE:Asia/Shanghai")
	dtstamp := time.Now().UTC().Format("20060102T150405Z")
	for _, e := range events {
		icsLine(&b, "BEGIN:VEVENT")
		icsLine(&b, "UID:"+e.uid)
		icsLine(&b, "DTSTAMP:"+dtstamp)
		// 采用UTC时间，中国无夏令时，Apple/Google日历均兼容
		icsLine(&b, "DTSTART:"+e.start.UTC().Format("20060102T150405Z"))
		icsLine(&b, "DTEND:"+e.end.UTC().Format("20060102T150405Z"))
		icsLine(&b, "SUMMARY:"+icsEscape(e.summary))
		if e.location != "" {
			icsLine(&b, "LOCATION:"+icsEscape(e.location))
		}
		icsLine(&b, "DESCRIPTION:"+icsEscape(e.desc))
		icsLine(&b, "END:VEVENT")
	}
	icsLine(&b, "END:VCALENDAR")
	return b.String()
}

// XnXqName 学年学期展示名，如 2026-2027-1
func XnXqName(cfg *config.Config) string {
	var y int
	if _, err := fmt.Sscanf(cfg.Time.XueNian, "%d", &y); err == nil {
		return fmt.Sprintf("%d-%d-%s", y, y+1, cfg.Time.XueQi)
	}
	return fmt.Sprintf("%s-%s", cfg.Time.XueNian, cfg.Time.XueQi)
}

// GenerateICS 生成ICS日历内容
func GenerateICS(cfg *config.Config, idx CourseIndex) (string, error) {
	events, err := buildEvents(cfg, idx)
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "", errors.New("没有可生成日历的课程(仅为值为1的选课课程生成)")
	}
	return renderICS(events, "课程表 "+XnXqName(cfg)), nil
}

// ExportIcsFiles 生成ICS日历文件
func ExportIcsFiles(cfg *config.Config) error {
	if cfg.SemesterStartDate == "" {
		return errors.New("未设置semester_start_date(第1周星期一日期)，跳过ICS生成")
	}
	idx, err := LoadCourseIndex(cfg)
	if err != nil {
		return err
	}
	content, err := GenerateICS(cfg, idx)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("timetable_%s.ics", XnXqName(cfg))
	if err := os.WriteFile(fileName, []byte(content), 0666); err != nil {
		return err
	}
	log.Info("ICS日历已生成: ", fileName)
	return nil
}
