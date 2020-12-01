package user

import (
	"Go-000/Week02/app/model/user"
	"database/sql"
)

type Dao struct {
	db int
}

func New() (d *Dao) {
	return &Dao{}
}

func (d *Dao) FindById(id int64) (rows *user.User, err error) {
	if err == sql.ErrNoRows {

	}
	return nil, nil
}
