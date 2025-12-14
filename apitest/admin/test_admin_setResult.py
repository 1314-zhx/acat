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

@pytest.fixture
def auth_session():
    """为每个测试创建独立的已认证 session"""
    s = requests.Session()
    admin_login(s)
    return s


def set_result(session,slot_id,round):
    payload = {
        "slot_id":slot_id,
        "round":round
    }
    url = f"{config.BASE_URL}{config.SET_INTERVIEW_RESULT}"
    return session.post(url, json=payload, allow_redirects=False)

# 正向测试
def test_set_result_success(auth_session):
    resp = set_result(auth_session,4,2)
    assert resp.status_code==200
    json_data = resp.json()
    assert json_data.get("status")==200 # 业务码为200
    error = json_data.get("error")
    assert error  == ""


# 反向测试
def test_set_result_wrong_round(auth_session):
    resp = set_result(auth_session,4,3)
    assert resp.status_code==400
    json_data = resp.json()
    assert json_data.get("status")==400
    error = json_data.get("error")
    assert error != ""

def test_set_result_wrong_slotId(auth_session):
    resp = set_result(auth_session,0,2)
    assert resp.status_code==400
    json_data = resp.json()
    assert json_data.get("status")==400
    error = json_data.get("error")
    assert error != ""

