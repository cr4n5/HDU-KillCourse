import base64   
import binascii
import rsa
import lxml.etree
import requests
import time
import json

def encrypt(plain,n):
    try:
        plain=plain.encode("utf-8")
        n=base64.b64decode(n)
        n=binascii.hexlify(n)
        pubkey=rsa.PublicKey(int(n,16),65537)
        cipher=rsa.encrypt(plain,pubkey)
        output=''.join([("%x" % x).zfill(2) for x in cipher])
        output=binascii.unhexlify(output)
        output=base64.b64encode(output)
        return output.decode()
    except Exception as e:
        print(f"Encryption error: {e}")
        return None

def login(session : requests.Session) -> requests.Session:

    print("开始登录...")
    # 读取配置文件
    try:
        with open("config.json", "r") as f:
            config = json.load(f)
    except FileNotFoundError:
        print("Error: config.json file not found.")
        return None
    except json.JSONDecodeError:
        print("Error: Failed to decode JSON from config.json.")
        return None
    
    # 登录session配置
    url="https://newjw.hdu.edu.cn/jwglxt/xtgl/login_slogin.html"
    session.headers.update({
        "Accept": "text/html, application/xhtml+xml, application/xml; q=0.9, */*; q=0.8",
        "Accept-Language": "zh_CN",
        "Connection": "keep-alive",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363",
    })

    print("正在获取csrftoken...")
    # 获取csrftoken
    try:
        response = session.get(url)
        response.raise_for_status()
    except requests.RequestException as e:
        print(f"HTTP request error: {e}")
        return None
    
    # 解析csrftoken
    try:
        tree = lxml.etree.HTML(response.text)
        csrftoken = tree.xpath('//input[@name="csrftoken"]/@value')[0]
    except IndexError:
        print("Error: CSRF token not found.")
        return None
    
    print("正在获取公钥...")
    # 获取公钥
    pubkey_url = "https://newjw.hdu.edu.cn/jwglxt/xtgl/login_getPublicKey.html?time={}".format(int(time.time()))
    try:
        response = session.get(pubkey_url)
        response.raise_for_status()
        pubkey = response.json()
        n = pubkey["modulus"]
    except requests.RequestException as e:
        print(f"HTTP request error: {e}")
        return None
    except KeyError:
        print("Error: Public key not found in response.")
        return None
    
    # 加密密码
    password=config["login"]["password"]
    mm=encrypt(password,n)
    if mm is None:
        print("Error: Password encryption failed.")
        return None
    yhm=config["login"]["username"]
    data={
        "csrftoken":csrftoken,
        "yhm":yhm,
        "mm":mm,
    }

    print("正在登录...")
    # 登录
    try:
        response = session.post(url, data=data)
        response.raise_for_status()
    except requests.RequestException as e:
        print(f"HTTP request error: {e}")
        return None
    
    if "用户登录" in response.text:
        print("登录失败,请检查用户名和密码")
        return None
    else:
        print("登陆成功")
        return session

if __name__ == "__main__":
    kc_url="https://newjw.hdu.edu.cn/jwglxt/rwlscx/rwlscx_cxRwlsIndex.html?doType=query&gnmkdm=N1548"
    session=requests.Session()
    session=login(session)
   