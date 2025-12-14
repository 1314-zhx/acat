import requests
import pytest #不能删，要用pytest命令来运行测试
import config

# 发送登录请求
"""
成功登录后不跳转，该测试脚本仅测试登录逻辑
"""
def login(phone, password):

    payload = {"phone": phone, "password": password}
    response = requests.post(
        f"{config.BASE_URL}{config.LOGIN_ENDPOINT}",
        json=payload,
        allow_redirects=False  # 不跟随重定向
    )
    return response


# ===== 测试用例 =====
# 正确手机号和密码：应返回 200，并设置 token cookie
def test_admin_login_success():
    resp = login("15229300775", "123456")

    assert resp.status_code == 200
    assert "token" in resp.cookies
    # 可选：检查响应体结构
    json_data = resp.json()
    assert "status" in json_data
    assert json_data["status"] == 200


# 登录验证错误后都不设置token
# 密码错误：应返回 400
def test_admin_login_wrong_password():

    resp = login("15229300775", "1")

    assert resp.status_code == 400
    assert "token" not in resp.cookies

# 手机号不存在：应返回 400
def test_admin_login_nonexistent_phone():

    resp = login("11111111111", "123456")

    assert resp.status_code == 400
    assert "token" not in resp.cookies

# 缺少 phone 字段：ShouldBindJSON 失败，返回 400
def test_admin_login_missing_phone():

    payload = {"password": "123456"}
    resp = requests.post(f"{config.BASE_URL}{config.LOGIN_ENDPOINT}", json=payload)

    assert resp.status_code == 400

# 缺少 password 字段：ShouldBindJSON 失败，返回 400
def test_admin_login_missing_password():
    payload = {"phone": "13800138000"}
    resp = requests.post(f"{config.BASE_URL}{config.LOGIN_ENDPOINT}", json=payload)

    assert resp.status_code == 400

# 空请求体：ShouldBindJSON 失败，返回 400
def test_admin_login_empty_body():
    resp = requests.post(f"{config.BASE_URL}{config.LOGIN_ENDPOINT}", json={})

    assert resp.status_code == 400

# 非 JSON 请求：应返回 400
def test_admin_login_invalid_json():

    resp = requests.post(
        f"{config.BASE_URL}{config.LOGIN_ENDPOINT}",
        data="not a json",
        headers={"Content-Type": "application/json"}
    )
    assert resp.status_code == 400
