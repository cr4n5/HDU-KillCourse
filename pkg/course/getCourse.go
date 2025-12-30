package course

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
	"github.com/xuri/excelize/v2"
)

func GetCourse(c *client.Client, cfg *config.Config) (*client.GetCourseResp, error) {
	// 获取info
	err := c.GetStuInfo()
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 先从本地course.json读取课程信息
	courses, err := ReadCourse(cfg)
	if err == nil {
		return courses, nil
	}

	// 本地课程信息读取失败，从服务器获取课程信息
	log.Error("本地课程信息读取失败，正在从服务器获取课程信息...")
	log.Info("Notice！: 等待时间可能较长，请耐心等待...")
	// 在线获取课程信息
	XueNian := cfg.Time.XueNian
	intXueNian, err := strconv.Atoi(XueNian)
	if err != nil {
		return nil, errors.New("学年格式错误")
	}
	xnmc := fmt.Sprintf("%s-%s", XueNian, strconv.Itoa(intXueNian+1))
	xqmc := cfg.Time.XueQi
	courseResp, courseToExcelResp, err := GetCourseOnline(c, cfg, "")
	if err != nil {
		return nil, fmt.Errorf("在线获取课程信息失败: %w", err)
	}

	// 将课程信息保存为 Excel
	err = CourseRenameToExcel(courseToExcelResp, xnmc, xqmc)
	if err != nil {
		log.Error("保存课程信息到Excel失败: ", err)
	} else {
		log.Info("任务落实情况课程已导出到Excel文件中...")
	}

	// 保存课程信息到本地
	err = SaveCourse(courseResp)
	if err != nil {
		return nil, err
	}

	return courseResp, nil
}

// GetCourseOnline 在线获取课程
func GetCourseOnline(c *client.Client, cfg *config.Config, CourseName string) (*client.GetCourseResp, *client.GetCourseToExcelResp, error) {
	// 初始化请求
	XueNian := cfg.Time.XueNian
	intXueNian, err := strconv.Atoi(XueNian)
	if err != nil {
		return nil, nil, errors.New("学年格式错误")
	}
	xnmc := fmt.Sprintf("%s-%s", XueNian, strconv.Itoa(intXueNian+1))
	xqmc := cfg.Time.XueQi
	var xqm string
	if xqmc == "1" {
		xqm = "3"
	} else if xqmc == "2" {
		xqm = "12"
	} else {
		return nil, nil, errors.New("学期格式错误")
	}
	req := &client.GetCourseReq{
		// Cxfs:        "1",
		// Zymc:        "全部",
		Xnmc: xnmc,
		Xqmc: xqmc,
		// Kkxymc:      "全部",
		// Jgmc:        "全部",
		// Ywtk:        "0",
		// Skfs:        "0",
		Xnm:         XueNian,
		Xqm:         xqm,
		Search:      "false",
		Nd:          fmt.Sprintf("%d", time.Now().Unix()),
		ShowCount:   "9999",
		CurrentPage: "1",
		SortOrder:   "asc",
		Time:        "0",
		Jxbmc:       CourseName,
	}

	// 获取课程信息
	courseResp, courseToExcelResp, err := c.GetCourse(req)
	if err != nil {
		return nil, nil, err
	}
	return courseResp, courseToExcelResp, nil
}

// CourseRenameToExcel 将课程信息转换为Excel保存
func CourseRenameToExcel(course *client.GetCourseToExcelResp, xnmc string, xqmc string) error {
	// 创建 Excel 文件
	f := excelize.NewFile()

	// 创建工作表
	sheetName := "课程信息"
	index, _ := f.NewSheet(sheetName)

	// 设置表头
	headers := []string{
		"教学班名称", "课程号", "课程名称", "是否开课", "是否排课", "选课标记",
		"上课时间", "上课地点", "场地名称", "场地具体名称", "教职工信息",
		"教师出生日期", "教师性别", "开课部门", "学分", "授课学院",
		"教学班容量", "教学班人数", "选课人数", "面向对象", "授课班级",
		"授课详情", "开课类型", "课程归属", "讲课学时", "考核方式", "学科备注",
	}
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1) // 列索引从 1 开始
		cell := colName + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	for row, item := range course.Items {
		values := []string{
			item.Jxbmc, item.KchID, item.Kcmc, item.Kkztmc, item.Bpkbj, item.Xkbjmc,
			item.Sksj, item.Jxdd, item.Cdlbmc, item.Cdejlbmc, item.Jzgxx,
			item.Jscsrq, item.Jsxb, item.Kkbm, item.Xf, item.Zczymc,
			strconv.Itoa(item.Jxbrl), strconv.Itoa(item.Jxbrs), strconv.Itoa(item.Xkrs), item.Mxdx, item.Jxbzc,
			item.Skdxssxy, item.Kklxmc, item.Kcgsmc, item.Kczhxs, item.Khfsmc, item.Xkbz,
		}
		for col, value := range values {
			colName, _ := excelize.ColumnNumberToName(col + 1)
			cell := colName + strconv.Itoa(row+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置工作表为活动表
	f.SetActiveSheet(index)

	fileName := fmt.Sprintf("%s_%s_任务落实情况课程导出.xlsx", xnmc, xqmc)

	// 保存 Excel 文件
	if err := f.SaveAs(fileName); err != nil {
		return err
	}

	return nil
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
