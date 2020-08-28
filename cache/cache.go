package cache

import (
	"fmt"
	"sync"
	"time"
)

type ValueType int32

const (
	DEFAULT ValueType = iota
	STRING
	HASH
	Set
)

type MemCacheDB struct {
	//all keys
	keys map[string]ValueType
	ttl  map[string]time.Time
	//Data Structure
	s  str
	hm hmap
	hs hset
	// internal function
	name2func map[string]Cmd
	// storage limit
	count int
	msize int
}

type MemCache struct {
	l  sync.Mutex
	db *MemCacheDB
}

type Cmd func(result IResult)

func (db *MemCacheDB) doBeforeProcess(key string, cmdType ValueType) error {
	if db.count >= db.msize && db.keys[key] == DEFAULT {
		return fmt.Errorf("keys count limit: %d", db.msize)
	}
	valueType := db.keys[key]
	if cmdType != DEFAULT && valueType != DEFAULT && valueType != cmdType {
		return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	//if ttl exist, and NOW > ttl, lazy del key
	expireTime := db.ttl[key]
	if !expireTime.IsZero() && time.Now().After(expireTime) {
		_, err := db.delKey(key, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// valueType param can't be DEFAULT
func (db *MemCacheDB) addKey(key string, valueType ValueType) (bool, error) {
	// key don't exist before addKey
	if db.keys[key] == DEFAULT {
		db.count++
	}
	db.keys[key] = valueType
	return true, nil
}

func (db *MemCacheDB) delKey(key string, ttl bool) (bool, error) {
	valueType := db.keys[key]
	if db.keys[key] != DEFAULT {
		db.count--
		delete(db.keys, key)
	}
	if ttl {
		delete(db.ttl, key)
	}
	if valueType == STRING && db.s[key] != nil {
		delete(db.s, key)
		return true, nil
	} else if valueType == HASH && db.hm[key] != nil {
		delete(db.hm, key)
		return true, nil
	} else if valueType == Set && db.hs[key] != nil {
		delete(db.hs, key)
		return true, nil
	}
	return false, nil
}

func (db *MemCacheDB) register(cmd string, f Cmd) error {
	db.name2func[cmd] = f
	return nil
}

func NewMemCache(conf *CacheConf) (*MemCache, error) {
	s := &MemCache{
		l: sync.Mutex{},
		db: &MemCacheDB{
			keys:      make(map[string]ValueType),
			ttl:       make(map[string]time.Time),
			s:         initStr(),
			hm:        initHmap(),
			hs:        initHset(),
			name2func: map[string]Cmd{},
			msize:     conf.MaxSize,
			count:     0,
		},
	}
	// add a command init function when add a new data structure
	commandString(s.db)
	commandHashMap(s.db)
	commandHashSet(s.db)
	// ttl policy
	go func() {
		ttlPeriodMillSecond := conf.TtlPeriodMillSecond
		if ttlPeriodMillSecond <= 0 {
			ttlPeriodMillSecond = 100 //default 100ms
		}
		for {
			s.l.Lock()
			volatileRange(s.db)
			s.l.Unlock()
			time.Sleep(time.Duration(ttlPeriodMillSecond) * time.Millisecond)
		}
	}()

	return s, nil
}

func (s *MemCache) doWithTransaction(r IResult) {
	s.l.Lock()
	defer s.l.Unlock()
	cmdName := r.Name()
	s.db.name2func[cmdName](r)
}

//string api
//********************************************************************
func (s *MemCache) Set(key string, value []byte) *BoolResult {
	//Todo: check param.
	cmd := NewBoolResult("set", key, value)
	s.doWithTransaction(cmd)
	return cmd
}

func (s *MemCache) Get(key string) *BytesResult {
	cmd := NewBytesResult("get", key)
	s.doWithTransaction(cmd)
	return cmd
}

func (s *MemCache) Del(keys ...string) *IntResult {
	cmd := NewIntResult("del", keys)
	s.doWithTransaction(cmd)
	return cmd
}

func (s *MemCache) Expire(key string, seconds int) *IntResult {
	cmd := NewIntResult("expire", key, seconds)
	s.doWithTransaction(cmd)
	return cmd
}

//hashmap api
//********************************************************************

func (s *MemCache) HSet(key, field string, value []byte) *IntResult {
	cmd := NewIntResult("hset", key, field, value)
	s.doWithTransaction(cmd)
	return cmd
}

func (s *MemCache) HGet(key, field string) *BytesResult {
	cmd := NewBytesResult("hget", key, field)
	s.doWithTransaction(cmd)
	return cmd
}

func (s *MemCache) HDel(key string, field ...string) *IntResult {
	cmd := NewIntResult("hdel", key, field)
	s.doWithTransaction(cmd)
	return cmd
}

//hashset api
//********************************************************************

func (s *MemCache) SAdd(key string, members ...string) *IntResult {
	cmd := NewIntResult("sadd", key, members)
	s.doWithTransaction(cmd)
	return cmd
}

func (s *MemCache) SIsMember(key, member string) *IntResult {
	cmd := NewIntResult("sismember", key, member)
	s.doWithTransaction(cmd)
	return cmd
}
