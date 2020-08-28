package cache

import "fmt"

type hset map[string]map[string]float64

func initHset() hset {
	return make(hset)
}

// register cmd when add a operate
func commandHashSet(db *MemCacheDB) {
	db.register("sadd", db.sAdd)
	db.register("sismember", db.sIsMember)
}

// member exist return 0， new member return new member count
func (db *MemCacheDB) sAdd(result IResult) {
	if len(result.Args()) < 2 {
		result.SetError(fmt.Errorf("sadd need at least 2 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("sadd argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].([]string)
	if !ok {
		result.SetError(fmt.Errorf("sadd argument 2 shuold be []string"))
		return
	}
	err := db.doBeforeProcess(arg0, Set)
	if err != nil {
		result.SetError(err)
		return
	}
	res := 0
	if db.hs[arg0] == nil {
		db.addKey(arg0, Set)
		db.hs[arg0] = make(map[string]float64)
	}
	keys := arg1
	for _, member := range keys {
		// if member not exist, save & res++
		if db.hs[arg0][member] == 0 {
			res++
			db.hs[arg0][member] = 1
		}
	}
	result.SetVal(res)
}

// member exist return 0， new member return new member count
func (db *MemCacheDB) sIsMember(result IResult) {
	if len(result.Args()) != 2 {
		result.SetError(fmt.Errorf("sismember need 2 argument"))
		return
	}
	arg0, ok := result.Args()[0].(string)
	if !ok {
		result.SetError(fmt.Errorf("sismember argument 1 shuold be string"))
		return
	}
	arg1, ok := result.Args()[1].(string)
	if !ok {
		result.SetError(fmt.Errorf("sismember argument 2 shuold be string"))
		return
	}
	err := db.doBeforeProcess(arg0, Set)
	if err != nil {
		result.SetError(err)
		return
	}
	//key not exist or member not exist
	if db.hs[arg0] == nil || db.hs[arg0][arg1] == 0 {
		result.SetVal(0)
		return
	}
	result.SetVal(1)
}
