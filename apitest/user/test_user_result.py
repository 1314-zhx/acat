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

def get_interview_result(session: requests.Session, round_num: int):
    payload = {
        "round": round_num
    }
    url = f"{config.BASE_URL}{config.RESULT_INTERVIEW}"
    return session.post(url, json=payload, allow_redirects=False)

@pytest.fixture
def auth_session():
    s = requests.Session()
    user_login(s)
    return s

# 正向测试
def test_user_result_success(auth_session):
    resp = get_interview_result(auth_session, 1)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200

# 反向测试
def test_user_result_wrong_round(auth_session):
    resp = get_interview_result(auth_session, 3)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "无效参数"
