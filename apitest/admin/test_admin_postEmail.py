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
def post_email(session,user_id,name,round,email,customize,content,test_mode):
    payload = {
        "user_id":user_id,
        "name":name,
        "round":round,
        "email":email,
        "customize":customize,
        "content":content,
        "test_mode": test_mode
    }
    url = f"{config.BASE_URL}{config.POST_EMAIL}"
    return session.post(url,json=payload,allow_redirects=False)
# 正向
def test_post_email_success(auth_session):
    resp = post_email(auth_session,10,"张皓翔",1,"2998759818@qq.com",True,"测试",True)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200
    assert json_data.get("error") == ""
# 反向
def test_post_email_wrong_user(auth_session):
    resp = post_email(auth_session,100,"张皓翔",1,"2998759818@qq.com",True,"测试",True)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "user not found"

def test_post_email_wrong_email(auth_session):
    resp = post_email(auth_session,10,"张皓翔",1,"2998759818qq.com",True,"测试",True)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "邮件格式不对"

def test_post_email_exceed_content(auth_session):
    resp = post_email(auth_session,10,"张皓翔",1,"2998759818@qq.com",True,
                      "1111111111111111111111111111111111111111111111111111111111111111",True)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "正文超过50字"
