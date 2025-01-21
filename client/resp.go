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
	} `json:"tmpList"`
}
