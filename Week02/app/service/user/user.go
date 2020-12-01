package user

import (
	userdao "Go-000/Week02/app/dao/user"
	"Go-000/Week02/app/model/user"
)

type Service struct {
	dao *userdao.Dao
}

func Init() (s *Service) {
	return &Service{
		dao: userdao.New(),
	}
}

func (s *Service) FindById(id int64) (*user.User, error) {
	return s.dao.FindById(id)
}
