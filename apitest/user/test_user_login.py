import pytest
import requests
import config
# 工具函数
def user_login(phone,password):
    payload = {
        "phone":phone,
        "password":password
    }
    url = f"{config.BASE_URL}{config.LOGIN_ROUTE}"
    return requests.post(url,json=payload,allow_redirects=False)

def user_login_missing_phone(password):
    payload = {
        "password":password
    }
    url = f"{config.BASE_URL}{config.LOGIN_ROUTE}"
    return requests.post(url,json=payload,allow_redirects=False)
# 正向
def test_user_login_success():
    resp = user_login("15229300775","123456")
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200
    assert json_data.get("error") == ""

# 反向
def test_user_login_wrong_phone():
    resp = user_login("15229377777","123456")
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "手机号或密码错误"

def test_user_login_invalid_param():
    resp = user_login("1","123456")
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "参数格式错误"

def test_user_login_missing_param():
    resp = user_login_missing_phone("123456")
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "缺少必填字段"