import requests
import lxml.etree
import time
import json
import re
import hdu_login

def get_course_sort_id_true(list):
    pattern = r"queryCourse\(this,'(\d+)'"
    pattern1 = r"queryCourse\(this,'(?:[^']*)','(\w+)'"
    extracted_texts = {}
    for text in list:
        match = re.search(pattern, text)
        match1 = re.search(pattern1, text)
        if match:
            extracted_text = match1.group(1)
            id = match.group(1)
            extracted_texts[id] = extracted_text
    return extracted_texts

def get_xkkz_id(html):
    global ccdm,bh_id,jg_id,xsbj,xz,mzm,xslbdm,xbm,zyfx_id,xqh_id
    tree = lxml.etree.HTML(html)
    course_sort_id = tree.xpath('//a[@role="tab"]/@onclick')
    ccdm=tree.xpath('//input[@name="ccdm"]/@value')[0]
    bh_id=tree.xpath('//input[@name="bh_id"]/@value')[0]
    jg_id=tree.xpath('//input[@name="jg_id_1"]/@value')[0]
    xsbj=tree.xpath('//input[@name="xsbj"]/@value')[0]
    xz=tree.xpath('//input[@name="xz"]/@value')[0]
    mzm=tree.xpath('//input[@name="mzm"]/@value')[0]
    xslbdm=tree.xpath('//input[@name="xslbdm"]/@value')[0]
    xbm=tree.xpath('//input[@name="xbm"]/@value')[0]
    zyfx_id=tree.xpath('//input[@name="zyfx_id"]/@value')[0]
    xqh_id=tree.xpath('//input[@name="xqh_id"]/@value')[0]
    return get_course_sort_id_true(course_sort_id)


def get_do_jxb_id(kch_id,jxb_id,kklxdm):
    url="https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzbjk_cxJxbWithKchZzxkYzb.html"
    data={
        "bklx_id" : "0",
        "njdm_id" : njdm_id,
        "xkxnm" : XueNian,
        "xkxqm" : xqm,
        "kklxdm" : kklxdm,
        "kch_id" : kch_id,
        "xkkz_id" : xkkz_id[kklxdm],
        'xsbj': xsbj, 
        'ccdm': ccdm,
        'xz': xz, 
        'mzm': mzm, 
        'xslbdm': xslbdm, 
        'xbm': xbm, 
        'bh_id': bh_id,
        'zyfx_id': zyfx_id,
        'jg_id': jg_id,
        'xqh_id': xqh_id, 
    }
    response=session.post(url,data=data)
    response=response.json()
    for i in response:
        if i["jxb_id"]==jxb_id:
            return i["do_jxb_id"]

def choose(jxb_ids,kch_id,kklxdm,jxbzc):
    url="https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzbjk_xkBcZyZzxkYzb.html"
    data={
        "jxb_ids" : jxb_ids,
        "kch_id" : kch_id,
        "qz" : "0",
    }
    if kklxdm=="01":
        data["njdm_id"]="20"+jxbzc[0:2]
        data["zyh_id"]=jxbzc[2:6]
    response=session.post(url,data=data)
    print(response.text)

def not_choose(jxb_ids,kch_id):
    url="https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzb_tuikBcZzxkYzb.html"
    data={
        "jxb_ids" : jxb_ids,
        "kch_id" : kch_id,
    }
    response=session.post(url,data=data)
    print(response.text)

if __name__ == "__main__":
    with open('config.json', 'r') as f:
        config = json.load(f)

    XueNian=config["time"]["XueNian"]
    xnmc="{}-{}".format(XueNian,str(int(XueNian)+1))
    xqmc=config["time"]["XueQi"]
    if xqmc=='1':
        xqm="3"
    else:
        xqm="12"

    njdm_id="20"+config["login"]["username"][0:2]
    course=config["course"]

    session=requests.Session()
    session=hdu_login.login(session,config["login"]["username"],config["login"]["password"])

    xk_url="https://newjw.hdu.edu.cn/jwglxt/xsxk/zzxkyzb_cxZzxkYzbIndex.html?gnmkdm=N253512&layout=default"
    response=session.get(xk_url)
    print("xkkz_id获取")
    time.sleep(0.5)
    xkkz_id=get_xkkz_id(response.text)
    with open("hdu_xkkz_id.json","w") as f:
        json.dump(xkkz_id,f,ensure_ascii=False,indent=2)

    with open('hdu_kc.json', 'r') as f:
        courses = json.load(f)

    for i in course:
        for j in range(len(courses)):
            if courses[j]["jxbmc"]==i:
                kch_id=courses[j]["kch_id"]
                jxb_id=courses[j]["jxb_id"]
                jxbzc=courses[j]["jxbzc"]
                kklxmc=courses[j]["kklxmc"]
                break

        if kklxmc=="主修课程":
            kklxdm="01"
        elif kklxmc=="通识选修课":
            kklxdm="10"
        elif kklxmc=="体育分项":
            kklxdm="05"
        elif kklxmc=="特殊课程":
            kklxdm="09"
        
        do_jxb_id=get_do_jxb_id(kch_id,jxb_id,kklxdm)
        print("do_jxb_id获取")
        time.sleep(0.5)
        if course[i]=="1":
            choose(do_jxb_id,kch_id,kklxdm,jxbzc)
        else:
            not_choose(do_jxb_id,kch_id)
        time.sleep(3)