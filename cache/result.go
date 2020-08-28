package cache

import (
	"fmt"
	"strings"
)

type IResult interface {
	Name() string
	Args() []interface{}
	stringArg(int) string
	Err() error
	SetVal(val interface{})
	SetError(err error)
}

type result struct {
	//val interface{}
	_args []interface{}
	err   error
}

func (r *result) Err() error {
	return r.err
}

func (r *result) Args() []interface{} {
	if len(r._args) > 1 {
		return r._args[1:]
	}
	return nil
}

func (r *result) stringArg(pos int) string {
	if pos < 0 || pos >= len(r._args) {
		return ""
	}
	s, _ := r._args[pos].(string)
	return s
}

func (r *result) Name() string {
	if len(r._args) > 0 {
		// Cmd name must be lower cased.
		s := strings.ToLower(r.stringArg(0))
		r._args[0] = s
		return s
	}
	return ""
}

func (r *result) SetError(err error) {
	r.err = err
}

//func (r *Result) Bytes() ([]byte, error){
//	return nil, nil
//}
//
//func (r *Result) String() (string, error) {
//
//	return "", nil
//}
//
//
//func (r *Result) Raw() (interface{}, error) {
//	return r.val, r.err
//}

type BoolResult struct {
	result
	val bool
}

func NewBoolResult(args ...interface{}) *BoolResult {
	return &BoolResult{
		result: result{_args: args},
	}
}

func (r *BoolResult) Result() (bool, error) {
	return r.val, r.err
}

func (r *BoolResult) SetVal(val interface{}) {
	boolVal, ok := val.(bool)
	if !ok {
		r.err = fmt.Errorf("%s need a %s type val", "BoolResult", "bool")
		return
	}
	r.val = boolVal
}

type BytesResult struct {
	result
	val []byte
}

func NewBytesResult(args ...interface{}) *BytesResult {
	return &BytesResult{
		result: result{_args: args},
	}
}

func (r *BytesResult) Result() ([]byte, error) {
	return r.val, r.err
}

func (r *BytesResult) SetVal(val interface{}) {
	boolVal, ok := val.([]byte)
	if !ok {
		r.err = fmt.Errorf("%s need a %s type val", "BytesResult", "[]byte")
		return
	}
	r.val = boolVal
}

type IntResult struct {
	result
	val int
}

func NewIntResult(args ...interface{}) *IntResult {
	return &IntResult{
		result: result{_args: args},
	}
}

func (r *IntResult) Result() (int, error) {
	return r.val, r.err
}

func (r *IntResult) SetVal(val interface{}) {
	intVal, ok := val.(int)
	if !ok {
		r.err = fmt.Errorf("%s need a %s type val", "BytesResult", "[]byte")
		return
	}
	r.val = intVal
}
