package main

import (
	"Go-000/Week02/app/service/user"
	"fmt"
)

func main() {
	s := user.Init()
	if _, err := s.FindById(0); err != nil {
		fmt.Println("%+v", err)
	}
	fmt.Println("SUCCESS")
}
