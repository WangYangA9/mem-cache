package cache

import "fmt"

type hmap map[string]map[string][]byte

func initHmap() hmap {
	return make(hmap)
}

// register cmd when add a operate
func commandHashMap(db *MemCacheDB) {
	db.register("hset", db.hset)
	db.register("hget", db.hget)
	db.register("hdel", db.hdel)
}

// field exist return 0， new field return 1
func (db *MemCacheDB) hset(result IResult) {
	if len(result.Args()) != 3 {
		result.SetError(fmt.Errorf("hset need 3 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("hset argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].(string)
	if !ok {
		result.SetError(fmt.Errorf("hset argument 2 shuold be string"))
		return
	}
	arg2, ok := result.Args()[2].([]byte)
	if !ok {
		result.SetError(fmt.Errorf("hset argument 3 shuold be []byte"))
		return
	}
	err := db.doBeforeProcess(arg0, HASH)
	if err != nil {
		result.SetError(err)
		return
	}
	res := 0
	if db.hm[arg0] == nil {
		db.addKey(arg0, HASH)
		db.hm[arg0] = make(map[string][]byte)
	}
	if db.hm[arg0][arg1] == nil {
		res = 1
	}
	db.hm[arg0][arg1] = arg2

	result.SetVal(res)
}

// return string value
func (db *MemCacheDB) hget(result IResult) {
	if len(result.Args()) != 2 {
		result.SetError(fmt.Errorf("hget need 2 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("hgel argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].(string)
	if !ok {
		result.SetError(fmt.Errorf("hget argument 2 shuold be string"))
		return
	}
	err := db.doBeforeProcess(arg0, HASH)
	if err != nil {
		result.SetError(err)
		return
	}
	if db.hm[arg0] == nil {
		result.SetVal([]byte(""))
		return
	}
	result.SetVal(db.hm[arg0][arg1])
}

// field can be multi, return field count that del successful
func (db *MemCacheDB) hdel(result IResult) {
	if len(result.Args()) < 2 {
		result.SetError(fmt.Errorf("hdel need at least 2 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("hdel argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].([]string)
	if !ok {
		result.SetError(fmt.Errorf("hdel argument 2 shuold be []string"))
		return
	}
	err := db.doBeforeProcess(arg0, HASH)
	if err != nil {
		result.SetError(err)
		return
	}
	res := 0
	key := arg0
	if db.hm[key] == nil {
		result.SetVal(res)
		return
	}
	keys := arg1
	for _, fieldTemp := range keys {
		if db.hm[key][fieldTemp] != nil {
			res++
			delete(db.hm[key], fieldTemp)
		}
	}
	result.SetVal(res)
}
