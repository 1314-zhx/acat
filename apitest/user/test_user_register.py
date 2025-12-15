import pytest
import requests
import config

def user_register(name,phone,password,re_password,email,stu_id,gender,direction):
    payload = {
        "name":name,
        "phone":phone,
        "password":password,
        "re_password":re_password,
        "email":email,
        "stu_id":stu_id,
        "gender":gender,
        "direction":direction
    }
    url = f"{config.BASE_URL}{config.REGISTER_ROUTE}"
    return requests.post(url,json=payload,allow_redirects=False)

def user_register_missing_email(name,phone,password,re_password,stu_id,gender,direction):
    payload = {
        "name":name,
        "phone":phone,
        "password":password,
        "re_password":re_password,
        "stu_id":stu_id,
        "gender":gender,
        "direction":direction
    }
    url = f"{config.BASE_URL}{config.REGISTER_ROUTE}"
    return requests.post(url,json=payload,allow_redirects=False)

# 正向测试
def test_user_register_success():
    resp = user_register("zhx","15339300775","123456","123456","2998759818@ww.com","2400413083",1,1)
    assert resp.status_code == 200
    json_data = resp.json()
    assert json_data.get("status") == 200
    assert json_data.get("error") == ""

# 反向测试
def test_user_register_missing_email():
    resp = user_register_missing_email("zhx","15339300775","123456","123456","2400413083",1,1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") ==  "缺少必填字段"

def tset_user_register_wrong_param():
    resp = user_register("zhx","15339300775","123456","654321","2998759818@ww.com","2400413083",1,1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "参数错误"

def test_user_register_invalid_param():
    resp = user_register("zhx","1","123456","123456","2998759818@ww.com","2400413083",1,1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") =="格式不正确"

def test_user_register_exists():
    resp = user_register("zhx","15229300775","123456","123456","2998759818@ww.com","2400413083",1,1)
    assert resp.status_code == 400
    json_data = resp.json()
    assert json_data.get("status") == 400
    assert json_data.get("error") == "该手机号已被注册"