import requests
import config
import pytest

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
# 工具函数
def set_pass(session, user_id, slot_id, round, is_pass):
    payload = {
        "user_id": user_id,
        "slot_id": slot_id,
        "round": round,
        "is_pass": is_pass
    }
    url = f"{config.BASE_URL}{config.SET_INTERVIEW_PASS}"
    return session.post(url, json=payload, allow_redirects=False)

# 正向测试
def test_set_pass_success(auth_session):
    resp = set_pass(auth_session,10,6,1,1)
    assert resp.status_code==200 # gin 框架的http 码是200
    json_data = resp.json()
    assert json_data.get("status")==200 # 业务码是200
    assert json_data.get("error") == "" # 没有错误

# 反向测试
def test_set_pass_wrong_user(auth_session):
    resp = set_pass(auth_session,-1,6,1,1)
    assert resp.status_code==400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") != ""

def test_set_pass_wrong_sotId(auth_session):
    resp = set_pass(auth_session,10,-1,1,1)
    assert resp.status_code==400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") != ""

def test_set_pass_wrong_round(auth_session):
    resp = set_pass(auth_session,10,6,3,1)
    assert resp.status_code==400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") != ""

def test_set_pass_wrong_pass(auth_session):
    resp = set_pass(auth_session,10,6,3,3)
    assert resp.status_code==400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") != ""