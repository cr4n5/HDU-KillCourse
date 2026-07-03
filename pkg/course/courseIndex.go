package course

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/xuri/excelize/v2"
)

// CourseInfo 课程基本信息（用于排序、网页显示课程名、生成ICS日历）
type CourseInfo struct {
	Kcmc   string `json:"kcmc"`   // 课程名称
	Sksj   string `json:"sksj"`   // 上课时间
	Jxdd   string `json:"jxdd"`   // 上课地点
	Jzgxx  string `json:"jzgxx"`  // 教职工信息
	Kklxmc string `json:"kklxmc"` // 开课类型(主修课程/通识选修课/体育分项/特殊课程)
}

// CourseIndex 教学班名称 -> 课程信息
type CourseIndex map[string]CourseInfo

// jxbmcCodeRe 从教学班名称中提取课程号，如 (2026-2027-1)-A0600910-02 -> A0600910
var jxbmcCodeRe = regexp.MustCompile(`\)-(.+)-\d+$`)

// CourseCodeOfJxbmc 提取教学班名称中的课程号，失败时返回原字符串
func CourseCodeOfJxbmc(jxbmc string) string {
	m := jxbmcCodeRe.FindStringSubmatch(jxbmc)
	if m == nil {
		return jxbmc
	}
	return m[1]
}

// indexFromResp 从course.json的课程信息构建索引
func indexFromResp(courses *client.GetCourseResp) CourseIndex {
	idx := make(CourseIndex)
	if courses == nil {
		return idx
	}
	for _, item := range courses.Items {
		if _, ok := idx[item.Jxbmc]; !ok {
			idx[item.Jxbmc] = CourseInfo{Kcmc: item.Kcmc, Sksj: item.Sksj, Kklxmc: item.Kklxmc}
		}
	}
	return idx
}

// indexFromExcel 从导出的任务落实情况Excel构建索引（含地点、教师信息）
func indexFromExcel(cfg *config.Config) (CourseIndex, error) {
	xueNian := cfg.Time.XueNian
	if len(xueNian) != 4 {
		return nil, errors.New("学年格式错误")
	}
	// 与getCourse.go导出文件名保持一致，如 2026-2027_1_任务落实情况课程导出.xlsx
	var next int
	if _, err := fmt.Sscanf(xueNian, "%d", &next); err != nil {
		return nil, errors.New("学年格式错误")
	}
	fileName := fmt.Sprintf("%s-%d_%s_任务落实情况课程导出.xlsx", xueNian, next+1, cfg.Time.XueQi)

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("课程信息")
	if err != nil {
		return nil, err
	}

	idx := make(CourseIndex)
	for i, row := range rows {
		// 跳过表头；列: A教学班名称 C课程名称 G上课时间 H上课地点 K教职工信息
		if i == 0 || len(row) < 7 || row[0] == "" {
			continue
		}
		info := CourseInfo{Kcmc: row[2], Sksj: row[6]}
		if len(row) > 7 {
			info.Jxdd = row[7]
		}
		if len(row) > 10 {
			info.Jzgxx = row[10]
		}
		if len(row) > 22 {
			info.Kklxmc = row[22] // 开课类型
		}
		if _, ok := idx[row[0]]; !ok {
			idx[row[0]] = info
		}
	}
	return idx, nil
}

// LoadCourseIndex 加载课程索引：优先course.json，再用Excel补充地点教师等信息
func LoadCourseIndex(cfg *config.Config) (CourseIndex, error) {
	idx := make(CourseIndex)

	if courses, err := ReadCourse(cfg); err == nil {
		idx = indexFromResp(courses)
	}

	if excelIdx, err := indexFromExcel(cfg); err == nil {
		for k, v := range excelIdx {
			idx[k] = v // Excel信息更全，覆盖course.json
		}
	} else if len(idx) == 0 {
		log.Error("读取任务落实情况Excel失败: ", err)
	}

	if len(idx) == 0 {
		return nil, errors.New("无本地课程信息，请先运行一次程序获取课程信息(course.json或任务落实情况Excel)")
	}
	return idx, nil
}
