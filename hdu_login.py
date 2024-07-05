import base64   
import binascii
import lxml.html
import rsa
import lxml.etree
import requests
import time

def encrypt(plain,n):
    plain=plain.encode("utf-8")
    n=base64.b64decode(n)
    n=binascii.hexlify(n)
    pubkey=rsa.PublicKey(int(n,16),65537)
    cipher=rsa.encrypt(plain,pubkey)
    output=''.join([("%x" % x).zfill(2) for x in cipher])
    output=binascii.unhexlify(output)
    output=base64.b64encode(output)
    return output.decode()

def login(session,username,password):
    print("开始登陆")
    url="https://newjw.hdu.edu.cn/jwglxt/xtgl/login_slogin.html"
    session.headers.update({
        "Accept": "text/html, application/xhtml+xml, application/xml; q=0.9, */*; q=0.8",
        "Accept-Language": "zh_CN",
        "Connection": "keep-alive",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363",
    })
    response=session.get(url)
    tree=lxml.etree.HTML(response.text)
    csrftoken=tree.xpath('//input[@name="csrftoken"]/@value')[0]
    print("csrftoken获取")
    time.sleep(0.5)
    pubkey_url="https://newjw.hdu.edu.cn/jwglxt/xtgl/login_getPublicKey.html?time={}".format(int(time.time()))
    response=session.get(pubkey_url)
    pubkey=response.json()
    print("pubkey获取")
    time.sleep(0.5)
    n=pubkey["modulus"]
    mm=encrypt(password,n)
    yhm=username
    data={
        "csrftoken":csrftoken,
        "yhm":yhm,
        "mm":mm,
    }
    response=session.post(url,data=data)
    print("登陆成功")
    return session

if __name__ == "__main__":
    session=requests.Session()
    session=login(session)
   