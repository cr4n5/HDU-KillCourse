package client

import "net/url"

// LoginReq 登录请求
type LoginReq struct {
	Csrftoken string
	Username  string
	Password  string
}

// ToFormData 转换为表单数据
func (req *LoginReq) ToFormData() url.Values {
	return url.Values{
		"csrftoken": {req.Csrftoken},
		"yhm":       {req.Username},
		"mm":        {req.Password},
	}
}

// CasLoginReq cas登录请求
type CasLoginReq struct {
	Username    string
	Type        string
	EventID     string
	Geolocation string
	Execution   string
	CaptchaCode string
	Croypto     string
	Password    string
}

// ToFormData 转换为表单数据
func (req *CasLoginReq) ToFormData() url.Values {
	return url.Values{
		"username":     {req.Username},
		"type":         {req.Type},
		"_eventId":     {req.EventID},
		"geolocation":  {req.Geolocation},
		"execution":    {req.Execution},
		"captcha_code": {req.CaptchaCode},
		"croypto":      {req.Croypto},
		"password":     {req.Password},
	}
}

// GetCourseReq 获取课程请求
type GetCourseReq struct {
	Kkbm        string
	Kch         string
	Kcfzr       string
	Xsxy        string
	ZyhID       string
	BhID        string
	ZyfxID      string
	NjdmID      string
	Xsdm        string
	Jxdd        string
	Kklxdm      string
	XqhID       string
	Xkbj        string
	Kkzt        string
	Kclbdm      string
	Kcgsdm      string
	Kcxzdm      string
	Apksfsdm    string
	Ksfsdm      string
	Khfsdm      string
	Cxfs        string
	Jsssbm      string
	Zcm         string
	Xbdm        string
	CdlbID      string
	CdejlbID    string
	Jxbmc       string
	Sfzjxb      string
	Sfhbbj      string
	Zymc        string
	Xnmc        string
	Xqmc        string
	Kkxymc      string
	Jgmc        string
	Njmc        string
	Sfpk        string
	Sfwp        string
	Ywtk        string
	Skfs        string
	Dylx        string
	Jzglbm      string
	Jxms        string
	Skpt        string
	Sfhxkc      string
	Sfxwkc      string
	Sknr        string
	Bz          string
	Xkbz        string
	Sfzj        string
	Qsz         string
	ZykfkcbjCx  string
	SfgssxbkCx  string
	Zzz         string
	Xf          string
	JysID       string
	Xnm         string
	Xqm         string
	Js          string
	Kclxdm      string
	Search      string
	Nd          string
	ShowCount   string
	CurrentPage string
	QsortName   string
	SortOrder   string
	Time        string
}

// ToFormData 转换为表单数据
func (req *GetCourseReq) ToFormData() url.Values {
	return url.Values{
		"kkbm":                   {req.Kkbm},
		"kch":                    {req.Kch},
		"kcfzr":                  {req.Kcfzr},
		"xsxy":                   {req.Xsxy},
		"zyh_id":                 {req.ZyhID},
		"bh_id":                  {req.BhID},
		"zyfx_id":                {req.ZyfxID},
		"njdm_id":                {req.NjdmID},
		"xsdm":                   {req.Xsdm},
		"jxdd":                   {req.Jxdd},
		"kklxdm":                 {req.Kklxdm},
		"xqh_id":                 {req.XqhID},
		"xkbj":                   {req.Xkbj},
		"kkzt":                   {req.Kkzt},
		"kclbdm":                 {req.Kclbdm},
		"kcgsdm":                 {req.Kcgsdm},
		"kcxzdm":                 {req.Kcxzdm},
		"apksfsdm":               {req.Apksfsdm},
		"ksfsdm":                 {req.Ksfsdm},
		"khfsdm":                 {req.Khfsdm},
		"cxfs":                   {req.Cxfs},
		"jsssbm":                 {req.Jsssbm},
		"zcm":                    {req.Zcm},
		"xbdm":                   {req.Xbdm},
		"cdlb_id":                {req.CdlbID},
		"cdejlb_id":              {req.CdejlbID},
		"jxbmc":                  {req.Jxbmc},
		"sfzjxb":                 {req.Sfzjxb},
		"sfhbbj":                 {req.Sfhbbj},
		"zymc":                   {req.Zymc},
		"xnmc":                   {req.Xnmc},
		"xqmc":                   {req.Xqmc},
		"kkxymc":                 {req.Kkxymc},
		"jgmc":                   {req.Jgmc},
		"njmc":                   {req.Njmc},
		"sfpk":                   {req.Sfpk},
		"sfwp":                   {req.Sfwp},
		"ywtk":                   {req.Ywtk},
		"skfs":                   {req.Skfs},
		"dylx":                   {req.Dylx},
		"jzglbm":                 {req.Jzglbm},
		"jxms":                   {req.Jxms},
		"skpt":                   {req.Skpt},
		"sfhxkc":                 {req.Sfhxkc},
		"sfxwkc":                 {req.Sfxwkc},
		"sknr":                   {req.Sknr},
		"bz":                     {req.Bz},
		"xkbz":                   {req.Xkbz},
		"sfzj":                   {req.Sfzj},
		"qsz":                    {req.Qsz},
		"zykfkcbj_cx":            {req.ZykfkcbjCx},
		"sfgssxbk_cx":            {req.SfgssxbkCx},
		"zzz":                    {req.Zzz},
		"xf":                     {req.Xf},
		"jys_id":                 {req.JysID},
		"xnm":                    {req.Xnm},
		"xqm":                    {req.Xqm},
		"js":                     {req.Js},
		"kclxdm":                 {req.Kclxdm},
		"_search":                {req.Search},
		"nd":                     {req.Nd},
		"queryModel.showCount":   {req.ShowCount},
		"queryModel.currentPage": {req.CurrentPage},
		"queryModel.sortName":    {req.QsortName},
		"queryModel.sortOrder":   {req.SortOrder},
		"time":                   {req.Time},
	}
}

// GetDoJxbIdReq 获取do_jxb_id请求
type GetDoJxbIdReq struct {
	BklxID   string
	NjdmID   string
	Xkxnm    string
	Xkxqm    string
	Kklxdm   string
	KchID    string
	XkkzID   string
	Xsbj     string
	Ccdm     string
	Xz       string
	Mzm      string
	Xslbdm   string
	Xbm      string
	BhID     string
	ZyfxID   string
	JgID     string
	XqhID    string
	NjdmIDXs string
	ZyhIDXs  string
}

// ToFormData 转换为表单数据
func (req *GetDoJxbIdReq) ToFormData() url.Values {
	return url.Values{
		"bklx_id":    {req.BklxID},
		"njdm_id":    {req.NjdmID},
		"xkxnm":      {req.Xkxnm},
		"xkxqm":      {req.Xkxqm},
		"kklxdm":     {req.Kklxdm},
		"kch_id":     {req.KchID},
		"xkkz_id":    {req.XkkzID},
		"xsbj":       {req.Xsbj},
		"ccdm":       {req.Ccdm},
		"xz":         {req.Xz},
		"mzm":        {req.Mzm},
		"xslbdm":     {req.Xslbdm},
		"xbm":        {req.Xbm},
		"bh_id":      {req.BhID},
		"zyfx_id":    {req.ZyfxID},
		"jg_id":      {req.JgID},
		"xqh_id":     {req.XqhID},
		"njdm_id_xs": {req.NjdmIDXs},
		"zyh_id_xs":  {req.ZyhIDXs},
	}
}

// SelectCourseReq 选课请求
type SelectCourseReq struct {
	JxbIDs   string
	KchID    string
	Qz       string
	NjdmID   string
	ZyhID    string
	NjdmIDXs string
	ZyhIDXs  string
	XkkzID   string
}

// ToFormData 转换为表单数据
func (req *SelectCourseReq) ToFormData() url.Values {
	return url.Values{
		"jxb_ids":    {req.JxbIDs},
		"kch_id":     {req.KchID},
		"qz":         {req.Qz},
		"njdm_id":    {req.NjdmID},
		"zyh_id":     {req.ZyhID},
		"njdm_id_xs": {req.NjdmIDXs},
		"zyh_id_xs":  {req.ZyhIDXs},
		"xkkz_id":    {req.XkkzID},
	}
}

// CancelCourseReq 退课请求
type CancelCourseReq struct {
	JxbIDs string
	KchID  string
	Xkxnm  string
	Xkxqm  string
}

// ToFormData 转换为表单数据
func (req *CancelCourseReq) ToFormData() url.Values {
	return url.Values{
		"jxb_ids": {req.JxbIDs},
		"kch_id":  {req.KchID},
		"xkxnm":   {req.Xkxnm},
		"xkxqm":   {req.Xkxqm},
	}
}

// SearchCourseReq 搜索课程请求
type SearchCourseReq struct {
	Xkxnm      string
	Xkxqm      string
	Kklxdm     string
	Jspage     string
	Kspage     string
	Yllist     string // 是否有余量
	Filterlist string // 搜索内容
	NjdmIDXs   string
	ZyhIDXs    string
}

// ToFormData 转换为表单数据
func (req *SearchCourseReq) ToFormData() url.Values {
	return url.Values{
		"xkxnm":          {req.Xkxnm},
		"xkxqm":          {req.Xkxqm},
		"kklxdm":         {req.Kklxdm},
		"jspage":         {req.Jspage},
		"kspage":         {req.Kspage},
		"yl_list[0]":     {req.Yllist},
		"filter_list[0]": {req.Filterlist},
		"njdm_id_xs":     {req.NjdmIDXs},
		"zyh_id_xs":      {req.ZyhIDXs},
	}
}
