/*
因为管理员不能在公开场合注册，应该手动在数据库中插入
该文件用于，生成管理员的加密后的密码，将其手动插入进数据库中
*/
package main

import (
	"acat/util/encryption"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run gen_password.go <你的密码>")
		os.Exit(1)
	}
	password := os.Args[1]
	hash, err := encryption.HashPassword(password)
	if err != nil {
		panic(err)
	}
	fmt.Printf("生成的加密后的密码:\n%s\n", hash)
}
