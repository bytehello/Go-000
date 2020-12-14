## 答案

```go
dao: 

 return errors.Wrapf(code.NotFound, fmt.Sprintf("sql: %s error: %v", sql, err))


biz:

if errors.Is(err, code.NotFound} {

}
```

## 课后作业

> 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
  
Q: 
首先 考虑一种常用分层架构 api - service - dao

1. dao 层作为调用其他三方库 database/sql 的最后一层，应该需要通过wrap保存根因（课上老师明确讲了哦）
2. 因为 **HANDLE ONCE** 原则，所以 dao 包了一层传递给 service 层时，service 层也需要原封不动的继续网上抛，交给 api 层处理。当然这里也可以 service 降级处理,我个人认为还是抛错误到上层比较好
3. 考虑底层数据库可能更换(mongodb/mysql...), 不能直接返回 sql.ErrNoRows，这里我在 dao 层利用**不透明错误**的玩法，自己定义了错误类型，并向外提供 IsQueryNoRowsError 方法,最外层通过 pkg/errors 的Cause获取根因，然后再调用IsQueryNoRowsError就可以判断是否是查询没结果

dao 层代码 Week02/app/dao/user/user.go 具体如下:
```go
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

```
dao 层 error types 代码 Week02/app/dao/errors/errors.go
这里实践了 不透明错误
```go
package errors

import (
	"fmt"
)

// NoRows interface
type NoRows interface {
	IsNoRowsError() bool
}

// QueryNoRowsError
type QueryNoRowsError struct {
	Msg string
	Err error
}

func (q *QueryNoRowsError) Error() string {
	return fmt.Sprintf("msg is %v, err is \"%v\"", q.Msg, q.Err)
}

func (q *QueryNoRowsError) IsNoRowsError() bool {
	return true
}

// IsQueryNoRowsError determines if err is due to no rows in result set
func IsQueryNoRowsError(err error) bool {
	te, ok := err.(NoRows)
	return ok && te.IsNoRowsError()
}


```

api 层代码 Week02/main.go
 ```go
package main

import (
	goerrors "Go-000/Week02/app/dao/errors"
	"Go-000/Week02/app/service/user"
	"fmt"
)

func main() {
	s := user.Init()
	if _, err := s.FindById(0); err != nil {
		if goerrors.IsQueryNoRowsError(err) { // 遇到没有查询到数据
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


``` 

output is 
```
HTTP/1.1 404 
stack strace : 
msg is no result, err is "sql: no rows in result set"
wrap message : failed sql is "SELECT * FROM user WHERE id = 0" 
Go-000/Week02/app/dao/user.(*Dao).FindById
        /Users/Gechanghang/Code/Go-000/Week02/app/dao/user/user.go:29
Go-000/Week02/app/service/user.(*Service).FindById
        /Users/Gechanghang/Code/Go-000/Week02/app/service/user/user.go:19
main.main
        /Users/Gechanghang/Code/Go-000/Week02/main.go:12
runtime.main
        /usr/local/go/src/runtime/proc.go:200
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1337

```


## go error 学习笔记

#### 定义错误的各种方法

Opaque errors > error types > sentinel errors

不透明错误的处理除了老师课上讲得通过自定义error，也可以参考这里 https://github.com/kubernetes/apimachinery/blob/v0.19.4/pkg/api/errors/errors.go#L622

#### 原则：

错误只处理一次(handle once)

#### 用法：
1. 应用代码中，用pkg/error errors.New errors.Errorf 保存堆栈数据
2. 调用包内其他函数，直接简单返回，避免因为wrap 多倍的堆栈信息
3. 调用其他第三方库，调用基础库标准，自己的库，应该是wrap保存根因
4. 直接返回操作，不用导出打日志
5. 工作的goroutine顶部使用 %+v打印

#### 总结：

1. 作为kit/标准库（被很多人使用）、可重用性高的包只能返回根错误，eg sql/database就是如此，肯定不会wrap给你的，
2. 不打算处理的包/不降级，wrap起来，往上抛 ，举例：调用dao层报错，dao层作为与db库交互的最后一层， 会 wrap 错误上抛

#### 其他
底层如 db层出现一个 错误该如何处理, 是直接返回error record not found 还是包装一层在抛给上层

0. 考虑分层、考虑你的使用者
1. 应该考虑换数据库的情况，mongo、mysql的情况下
2. 返回空数据，二义性，
