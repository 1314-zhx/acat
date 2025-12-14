# test_admin_setSchedule.py

import requests
import pytest
import config

# 统一使用无秒、无时区的 datetime-local 格式
VALID_START = "2025-12-14T17:45"
VALID_END   = "2025-12-14T18:45"


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


def set_interview(session, start_time, end_time, max_num, interview_round):
    payload = {
        "start_time": start_time,
        "end_time": end_time,
        "max_num": max_num,
        "round": interview_round
    }
    url = f"{config.BASE_URL}{config.SET_INTERVIEW_ENDPOINT}"
    return session.post(url, json=payload, allow_redirects=False)



# 正向测试：所有参数合法

def test_admin_set_success(auth_session):
    resp = set_interview(auth_session, VALID_START, VALID_END, 50, 1)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200
    assert "msg" in json_data

# 反向测试
def test_admin_set_round_less_than_1(auth_session):
    resp = set_interview(auth_session, VALID_START, VALID_END, 50, 0)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data["status"] == 400
    assert json_data.get("error") not in (None, "")



def test_admin_set_round_greater_than_2(auth_session):
    resp = set_interview(auth_session, VALID_START, VALID_END, 50, 3)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data["status"] == 400
    assert json_data.get("error") not in (None, "")


def test_admin_set_max_num_less_than_1(auth_session):
    resp = set_interview(auth_session, VALID_START, VALID_END, 0, 1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data["status"] == 400
    assert json_data.get("error") not in (None, "")



def test_admin_set_max_num_greater_than_100(auth_session):
    resp = set_interview(auth_session, VALID_START, VALID_END, 51, 1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data["status"] == 400
    assert json_data.get("error") not in (None, "")


def test_admin_set_start_after_end(auth_session):
    resp = set_interview(auth_session, "2025-12-14T19:00", "2025-12-14T18:00", 50, 1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data["status"] == 400
    assert json_data.get("error") not in (None, "")