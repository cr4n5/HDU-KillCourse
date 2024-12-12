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
	} `json:"items"`
}

type GetDoJxbIdResp struct {
	JxbID   string `json:"jxb_id"`
	DoJxbID string `json:"do_jxb_id"`
}
