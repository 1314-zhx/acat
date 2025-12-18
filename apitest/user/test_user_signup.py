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

def signup_interview(session,name,direction,slot_id):
    payload = {
        "name":name,
        "direction":direction,
        "slot_id":slot_id
    }
    url = f"{config.BASE_URL}{config.SIGNUP_INTERVIEW}"
    return session.post(url, json=payload, allow_redirects=False)

@pytest.fixture
def auth_session():
    s = requests.Session()
    user_login(s)
    yield s
    s.close()

# 正向测试
def test_signup_interview_success(auth_session):
    resp = signup_interview(auth_session,"张皓翔",1,10)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200

# 反向测试
def test_signup_interview_wrong_round(auth_session):
    resp = signup_interview(auth_session,"张皓翔",1,0)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "查询面试表失败"