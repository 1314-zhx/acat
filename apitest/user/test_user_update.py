import pytest
import requests
import config

def user_login(session: requests.Session):
    payload = {
        "phone": "15229300775",
        "password": "123456"
    }
    url = f"{config.BASE_URL}{config.LOGIN_ROUTE}"
    resp = session.post(url, json=payload)
    assert resp.status_code == 200, f"登录失败: {resp.text}"
    json_data = resp.json()
    assert json_data.get("status") == 200, f"业务错误: {json_data.get('error', '未知')}"
    return resp

def update(session,name,direction,slot_id,is_delete):
    payload = {
        "name":name,
        "direction":direction,
        "slot_id":slot_id,
        "is_delete":is_delete
    }
    url = f"{config.BASE_URL}{config.SIGNUP_INTERVIEW}"
    return session.post(url, json=payload, allow_redirects=False)

@pytest.fixture
def auth_session():
    s = requests.Session()
    user_login(s)
    return s

# 正向测试更新
def test_user_update_success(auth_session):
    resp = update(auth_session,"张皓翔","2",10,0)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200

def test_user_delete_success(auth_session):
    resp = update(auth_session,"张皓翔","2",10,1)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200

# 反向测试
def test_user_update_no_slot(auth_session):
    resp = update(auth_session,"张皓翔","2",0,1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "没有该面试时段"