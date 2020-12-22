package main

import (
	goerrors "Go-000/Week02/app/dao/errors"
	"Go-000/Week02/app/service/user"
	"fmt"
	pkgerrors "github.com/pkg/errors"
)

func main() {
	s := user.Init()
	if _, err := s.FindById(0); err != nil {
		if goerrors.IsQueryNoRowsError(pkgerrors.Cause(err)) { // 遇到没有查询到数据
			fmt.Println("HTTP/1.1 404 ")
		} else { // 其他错误处理
			fmt.Println("HTTP/1.1 500 ")
		}
		fmt.Printf("stack strace : \n%+v\n", err)
		return
	}
	// 正常处理
	fmt.Println("HTTP/1.1 200 OK")
}
