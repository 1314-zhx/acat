import pytest
import requests
import config

# 工具函数，模拟请求
def reset_password(account,password,code):
    payload = {
        "account":account,
        "new_password":password,
        "code":code,
    }
    url = f"{config.BASE_URL}{config.RESET_PASSWORD_ROUTE}"
    return requests.post(url,json=payload,allow_redirects=False)


# 工具函数，获得验证码，随便测试 forget 接口
@pytest.fixture
def get_code():
    s = requests.Session()
    payload = {
        "param": "2998759818@qq.com",
        "test_mode": True
    }
    url = f"{config.BASE_URL}{config.FORGET_ROUTE}"
    resp = s.post(url, json=payload)
    assert resp.status_code == 200, f"请求失败: {resp.text}"
    json_data = resp.json()
    assert json_data.get("status") == 200, f"业务错误: {json_data.get('error', '未知错误')}"
    code = json_data.get("data", {})
    assert code is not None, "响应中未包含验证码（请确认后端在 test_mode 下返回 code）"
    print(code)
    s.close()
    return  code


# 正向测试
def test_reset_password_success(get_code):
    resp = reset_password("2998759818@qq.com","123456",get_code)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200
    assert json_data.get("error") == ""

# 反向测试
def test_reset_password_wrong_email(get_code):
    resp = reset_password("2998759818qq.com","123456",get_code)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "无效账号格式"

def test_reset_password_missing_param(get_code):
    resp = reset_password("","123456",get_code)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "缺少必要参数"

def test_reset_password_wrong_code(get_code):
    resp = reset_password("2998759818@qq.com","123456","1")
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "验证码不匹配"
