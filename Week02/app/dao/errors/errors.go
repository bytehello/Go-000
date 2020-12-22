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
