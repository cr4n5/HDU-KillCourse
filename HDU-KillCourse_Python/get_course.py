import requests
import time
import json
from HDU_Login import newjw_login

with open('config.json', 'r') as f:
    config = json.load(f)

url="https://newjw.hdu.edu.cn/jwglxt/rwlscx/rwlscx_cxRwlsIndex.html?doType=query&gnmkdm=N1548"

nd=int(time.time())
XueNian=config["time"]["XueNian"]
xnmc="{}-{}".format(XueNian,str(int(XueNian)+1))
xqmc=config["time"]["XueQi"]
if xqmc=='1':
    xqm="3"
else:
    xqm="12"

data={
    "kkbm": "",
    "kch": "",
    "kcfzr": "",
    "xsxy": "",
    "zyh_id": "",
    "bh_id": "",
    "zyfx_id": "",
    "njdm_id": "",
    "xsdm": "",
    "jxdd": "",
    "kklxdm": "",
    "xqh_id": "",
    "xkbj": "",
    "kkzt": "",
    "kclbdm": "",
    "kcgsdm": "",
    "kcxzdm": "",
    "apksfsdm": "",
    "ksfsdm": "",
    "khfsdm": "",
    "cxfs": "1",
    "jsssbm": "",
    "zcm": "",
    "xbdm": "",
    "cdlb_id": "",
    "cdejlb_id": "",
    "jxbmc": "",
    "sfzjxb": "",
    "sfhbbj": "",
    "zymc": "全部",
    "xnmc": xnmc,
    "xqmc": xqmc,
    "kkxymc": "全部",
    "jgmc": "全部",
    "njmc": "",
    "sfpk": "",
    "sfwp": "",
    "ywtk": "0",
    "skfs": "0",
    "dylx": "",
    "jzglbm": "",
    "jxms": "",
    "skpt": "",
    "sfhxkc": "",
    "sfxwkc": "",
    "sknr": "",
    "bz": "",
    "xkbz": "",
    "sfzj": "",
    "qsz": "",
    "zykfkcbj_cx": "",
    "sfgssxbk_cx": "",
    "zzz": "",
    "xf": "",
    "jys_id": "",
    "xnm": XueNian,
    "xqm": xqm,
    "js": "",
    "kclxdm": "",
    "_search": "false",
    "nd": nd,
    "queryModel.showCount": "9999",
    "queryModel.currentPage": "1",
    "queryModel.sortName": "",
    "queryModel.sortOrder": "asc",
    "time": "0"
}

session=requests.Session()
session=newjw_login.login(session)
if session is None:
    print("登录失败！")
    exit(1)

print("正在获取课程信息...")
response=session.post(url,data=data)
print("获取成功！")
response=response.json()
items=response["items"]
filtered_items = []
for item in items:
    filtered_item = {
        "jxbmc": item["jxbmc"],
        "kch_id": item["kch_id"],
        "jxb_id": item["jxb_id"],
        "jxbzc": item["jxbzc"],
        "kklxmc": item["kklxmc"],
    }
    filtered_items.append(filtered_item)
with open("hdu_kc.json","w") as f:
    json.dump(filtered_items,f,ensure_ascii=False,indent=2)
print("数据已保存到 hdu_kc.json")