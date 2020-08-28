package cache

import (
	"fmt"
	"time"
)

type str map[string][]byte

func initStr() str {
	return make(str)
}

// register cmd when add a operate
func commandString(db *MemCacheDB) {
	db.register("set", db.set)
	db.register("get", db.get)
	db.register("del", db.del)
	db.register("expire", db.expire)
}

// return a string
func (db *MemCacheDB) get(result IResult) {
	if len(result.Args()) != 1 {
		result.SetError(fmt.Errorf("get need 1 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("get argument 1 shuold be string"))
		return
	}
	err := db.doBeforeProcess(arg0, STRING)
	if err != nil {
		result.SetError(err)
		return
	}
	val := db.s[arg0]
	result.SetVal(val)
}

// return true
func (db *MemCacheDB) set(result IResult) {
	if len(result.Args()) != 2 {
		result.SetError(fmt.Errorf("set need 2 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("set argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].([]byte)
	if !ok {
		result.SetError(fmt.Errorf("set argument 2 shuold be []byte"))
		return
	}
	err := db.doBeforeProcess(arg0, STRING)
	if err != nil {
		result.SetError(err)
		return
	}
	_, err = db.addKey(arg0, STRING)
	if err != nil {
		result.SetError(err)
		return
	}
	db.s[arg0] = arg1
	result.SetVal(true)
}

// keys can be multi, return count that keys be deleted
func (db *MemCacheDB) del(result IResult) {
	if len(result.Args()) < 1 {
		result.SetError(fmt.Errorf("del need at least 1 argument"))
		return
	}
	keys, ok := result.Args()[0].([]string)
	if !ok {
		result.SetError(fmt.Errorf("del keys must be []string"))
		return
	}
	for _, key := range keys {
		err := db.doBeforeProcess(key, DEFAULT)
		if err != nil {
			result.SetError(err)
			return
		}
	}

	res := 0
	for _, elem := range keys {
		d, _ := db.delKey(elem, true)
		if d {
			res++
		}
	}
	result.SetVal(res)
}

func (db *MemCacheDB) expire(result IResult) {
	if len(result.Args()) != 2 {
		result.SetError(fmt.Errorf("expire need 2 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("expire argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].(int)
	if !ok {
		result.SetError(fmt.Errorf("expire argument 2 shuold be integer in [1, 2147483647]"))
		return
	}
	if arg1 <= 0 {
		result.SetError(fmt.Errorf("expire seconds can't <= 0, should be integer in [1, 2147483647]"))
	}
	if db.keys[arg0] == DEFAULT { //key don't exist
		result.SetVal(0)
		return
	}
	db.ttl[arg0] = time.Now().Add(time.Duration(arg1) * time.Second)
	result.SetVal(1)
}
