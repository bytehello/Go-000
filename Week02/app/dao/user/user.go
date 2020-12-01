package user

import (
	goerrors "Go-000/Week02/app/dao/errors"
	"Go-000/Week02/app/model/user"
	databasesql "database/sql"
	"fmt"
	pkgerrors "github.com/pkg/errors"
)

type Dao struct {
	db int
}

func New() (d *Dao) {
	return &Dao{}
}

func (d *Dao) FindById(id int64) (rows *user.User, err error) {
	sql := fmt.Sprintf("SELECT * FROM user WHERE id = %d", id)
	rows, err = d.query(sql)
	// 没有出错，直接返回
	if err == nil {
		return rows, nil
	}

	if err == databasesql.ErrNoRows { // 没有查询到数据
		// wrap 包装
		return nil, pkgerrors.Wrap(&goerrors.QueryNoRowsError{ // 这里是自定义的错误类型
			Msg: "no result",
			Err: err,
		}, fmt.Sprintf("wrap message : failed sql is \"%s\" ", sql))
	} else {
		return nil, pkgerrors.Wrap(err, sql)
	}
}

func (d *Dao) query(sql string) (rows *user.User, err error) {
	return nil, databasesql.ErrNoRows
}
