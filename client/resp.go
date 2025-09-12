package client

type GetPublicKeyResp struct {
	Modules  string `json:"modulus"`
	Exponent string `json:"exponent"`
}

type GetCourseResp struct {
	Items []struct {
		Jxbmc  string `json:"jxbmc"`
		KchID  string `json:"kch_id"`
		JxbID  string `json:"jxb_id"`
		Jxbzc  string `json:"jxbzc"`
		Kklxmc string `json:"kklxmc"`
		Kcmc   string `json:"kcmc"` // 课程名称
		Sksj   string `json:"sksj"` // 上课时间
	} `json:"items"`
}

type GetCourseToExcelResp struct {
	Items []struct {
		Jxbmc    string `json:"jxbmc"`    // 教学班名称
		KchID    string `json:"kch_id"`   // 课程号
		Kcmc     string `json:"kcmc"`     // 课程名称
		Kkztmc   string `json:"kkztmc"`   // 是否开课
		Bpkbj    string `json:"bpkbj"`    // 是否排课
		Xkbjmc   string `json:"xkbjmc"`   // 选课标记
		Sksj     string `json:"sksj"`     // 上课时间
		Jxdd     string `json:"jxdd"`     // 上课地点
		Cdlbmc   string `json:"cdlbmc"`   // 场地名称
		Cdejlbmc string `json:"cdejlbmc"` // 场地具体名称
		Jzgxx    string `json:"jzgxx"`    // 教职工信息
		Jscsrq   string `json:"jscsrq"`   // 教师出生日期
		Jsxb     string `json:"jsxb"`     // 教师性别
		Kkbm     string `json:"kkbm"`     // 开课部门
		Xf       string `json:"xf"`       // 学分
		Zczymc   string `json:"zczymc"`   // 授课学院
		Jxbrl    int    `json:"jxbrl"`    // 教学班容量
		Jxbrs    int    `json:"jxbrs"`    // 教学班人数
		Xkrs     int    `json:"xkrs"`     // 选课人数
		Mxdx     string `json:"mxdx"`     // 面向对象
		Jxbzc    string `json:"jxbzc"`    // 授课班级
		Skdxssxy string `json:"skdxssxy"` // 授课详情
		Kklxmc   string `json:"kklxmc"`   // 开课类型
		Kcgsmc   string `json:"kcgsmc"`   // 课程归属
		Kczhxs   string `json:"kczhxs"`   // 讲课学时
		Khfsmc   string `json:"khfsmc"`   // 考核方式
		Xkbz     string `json:"xkbz"`     // 学科备注
	} `json:"items"`
}

type GetDoJxbIdResp struct {
	JxbID   string `json:"jxb_id"`
	DoJxbID string `json:"do_jxb_id"`
}

type SelectCourseResq struct {
	Flag string `json:"flag"`
	Msg  string `json:"msg"`
}

type SearchCourseResp struct {
	TmpList []struct {
		Jxbmc string `json:"jxbmc"`
		// Kklxmc string `json:"kklxmc"`
		// KchID  string `json:"kch_id"`
		// JxbID  string `json:"jxb_id"`
		// Jxbzc  string `json:"jxbzc"`
		// Kcmc   string `json:"kcmc"` // 课程名称
		// Sksj   string `json:"sksj"` // 上课时间
	} `json:"tmpList"`
}

type QrLoginIdResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type QrLoginStatusResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// GetReleasesResp 获取最新版本信息的响应结构体
type GetReleasesResp struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

type GetStuInfoResp struct {
	Xsxx struct {
		NJDMID string `json:"NJDM_ID"`
		ZYHID  string `json:"ZYH_ID"`
	} `json:"xsxx"`
}

type GetZyhIdByBhResp []struct {
	ZyhID string `json:"zyh_id"`
}
