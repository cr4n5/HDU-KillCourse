# HDU-KillCourse

>本项目仅供学习和研究使用请于24小时内删除。使用本项目所产生的任何后果由使用者自行承担。在使用本项目之前，请确保您已充分了解相关法律法规，并确保您的行为符合所在国家或地区的法律要求。未经授权的情况下，请勿将本项目用于商业用途或其他非法用途。转载使用请标明出处。

- <img src="./Doc/img/香草蛋糕.jpg" width="100" height="100" alt="huohuo">

`杭电 抢课×选课√`

## 简介

支持主修，选修，体育课程，特殊课程

> [!TIP]
>
> If you are good at using it, you'll discover some pleasant surprises.

## 环境

Go 1.23

## 使用

1. 下载编译文件

- 在 [Releases](https://github.com/cr4n5/HDU-KillCourse/releases)中，下载对应系统的可执行文件。

- Or

```shell
go build
```

2. 修改配置

- 下载 [config.example.json](./config.example.json) 文件。
- 进入 [config.example.json](./config.example.json) 文件，修改对应内容。
- 配置名更改为 config.json。

```
{
    "cas_login": {
        "username": "2201xxxx",//杭电统一身份认证账号密码
        "password": "xxxxxxxx",
        "level:" : "0" //优先级
    },
    "newjw_login": {
        "username": "2201xxxx",//正方教务系统账号密码
        "password": "xxxxxxxx",
        "level:" : "1" //优先级
    }, // 0<1 所以优先使用cas登录 所以0比1大 数学天才
    "cookies": {
        "JSESSIONID": "",// 每次登录cookie参数都会自动更新
        "route": "",
        "enabled": "1"//如若登录过期，将enabled改为0，将不会使用cookies登录
    },
    "time": {
        "XueNian": "2024",//所选课程所在的学年学期，如2024-2025-1
        "XueQi": "1"
    },
    //课程教学班名称，如(2024-2025-1)-C2092011-01
    "course" : {
        "(2024-2025-1)-C2092011-01" : "1",//1为选课，0为退课
        "(2024-2025-1)-T1300019-04" : "1",
        "(2024-2025-1)-T1300019-05" : "1",
        "(2024-2025-1)-B2700380-02" : "0",
        "(2024-2025-1)-C2892008-02" : "1",
        "(2024-2025-1)-W0001321-06" : "0"
    },
    //课程按顺序执行
    "start_time": "2024-07-25 12:00:00",//程序开始时间
    "ClientBodyConfigEnabled": "0",//默认
    "DontTouchForDebug": "0"//默认
}
```

3. 选课

- 选课之前，可先去<a href='https://github.com/cr4n5/HDU-course_list'>杭电课程导出</a>，排好课表，获取课程教学班名称

> [!NOTE]
>
> 需在任务落实查询开放后，并在选课之前（省去在选课时查询课程请求）执行一次可执行文件获取课程信息

- 保证可执行文件和config.json在同一级目录下，然后在开始前几分钟执行可执行文件即可

## 协议

[Apache License 2.0](./LICENSE)
