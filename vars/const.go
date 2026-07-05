package vars

var (
	// 不需要debug的url
	NoDebugUrl = map[string]bool{
		"https://sso.hdu.edu.cn/login":                                                                   true,
		"https://newjw.hdu.edu.cn/sso/driot4login":                                                       true,
		"https://newjw.hdu.edu.cn/jwglxt/xtgl/login_slogin.html":                                         true,
		"https://newjw.hdu.edu.cn/jwglxt/rwlscx/rwlscx_cxRwlsIndex.html?doType=query&gnmkdm=N1548":       true, // 获取课程
		"https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzb_cxZzxkYzbIndex.html?gnmkdm=N253512&layout=default": true, // 选课配置
		"https://api.github.com/repos/cr4n5/HDU-KillCourse/releases/latest":                              true, // 仓库
	}
)

// Version 当前版本
const Version = "v1.4.9"

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 Edg/149.0.0.0"
