package version

import (
	"fmt"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/cr4n5/HDU-KillCourse/vars"
)

func VersionUpdate() {
	c := client.NewClient(nil)
	// 获取最新版本信息
	releaseResp, err := c.GetReleases()
	if err != nil {
		log.Error("获取最新版本信息失败: ", err)
		return
	}

	// 检查当前版本是否需要更新
	if releaseResp.TagName == vars.Version {
		return
	}

	log.Info(log.ErrorColor(fmt.Sprintf("Notice！: 有新版本 %s 可用，当前版本 %s", releaseResp.TagName, vars.Version)))
	log.Info(log.ErrorColor("\nNotice！: 更新内容: " + releaseResp.Body))
}
