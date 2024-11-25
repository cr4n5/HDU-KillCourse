# HDU-KillCourse
杭电 抢课×选课√  脚本

## 简介

支持主修，选修，体育课程，特殊课程

## 环境

python3

## 使用

1. 安装依赖

```shell
git clone --recursive https://github.com/cr4n5/HDU-KillCourse.git
cd HDU-KillCourse
pip install -r requirements.txt
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
    }
    //课程按顺序执行
}
```

3. 获取课程信息

需在任务落实查询开放后，并在选课之前（省去在选课时查询课程请求，不对土豆服务器造成过多压力）获取课程信息

```shell
python get_course.py
```

4. 选课

选课之前，可先去<a href='https://github.com/cr4n5/HDU-course_list'>杭电课程导出</a>，排好课表，获取课程教学班名称

```shell
python kill_course.py
```
