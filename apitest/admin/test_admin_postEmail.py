import requests
import pytest
import config

def admin_login(session):
    url = f"{config.BASE_URL}{config.LOGIN_ENDPOINT}"
    payload = {
        "phone": "15229300775",
        "password": "123456"
    }
    resp = session.post(url, json=payload)
    assert resp.status_code == 200, f"Login failed: {resp.text}"
    return resp

# 用来进行测试前的用户身份登录，因为该模块需要用户身份验证
# 并由pytest将其结果注入到与函数同名的参数中
@pytest.fixture
def auth_session():
    """为每个测试创建独立的已认证 session"""
    s = requests.Session() # 为可变变量，可用在其它函数中被改变并影响此处
    admin_login(s)
    return s
def post_email(session,user_id,name,round,email,customize,content):
    payload = {
        "user_id":user_id,
        "name":name,
        "round":round,
        "email":email,
        "customize":customize,
        "content":content
    }
    url = f"{config.BASE_URL}{config.POST_EMAIL}"
    response = session.post(url,json=payload,allow_redirects=False)
