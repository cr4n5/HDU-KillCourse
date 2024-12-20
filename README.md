# HDU-KillCourse
杭电 抢课×选课√  脚本

## 简介

支持主修，选修，体育课程，特殊课程

> [!NOTE]
>
> 本项目经过Go重构，暂未测试。Python版本见[HDU-KillCourse_Python](./HDU-KillCourse_Python/README.md)

## 环境

Go 1.23

Python 3版本见[HDU-KillCourse_Python](./HDU-KillCourse_Python/README.md)

## 使用

1. 下载编译文件

- 在 [Releases](https://github.com/cr4n5/HDU-KillCourse/releases)中，下载对应系统的可执行文件。

- Or

```shell
go build
```

2. 修改配置

- 进入 [config.example.json](./config.example.json) 文件，修改对应内容。
- 配置名更改为 config.json。

```
{
    "login": {
        "username": "2201xxxx",//教务系统账号密码（非数字杭电统一认证账号密码！！）
        "password": "xxxxxxxx"
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
    "start_time": "2024-07-25 12:00:00"//程序开始时间
}
```

3. 选课

- 选课之前，可先去<a href='https://github.com/cr4n5/HDU-course_list'>杭电课程导出</a>，排好课表，获取课程教学班名称

> [!NOTE]
>
> 需在任务落实查询开放后，并在选课之前（省去在选课时查询课程请求）执行一次可执行文件获取课程信息

- 保证可执行文件和config.json在同一级目录下，然后执行可执行文件即可
