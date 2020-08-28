package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test1", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetNil(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Get("test1").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Del("test1").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HGet("test1", "111").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HDel("test1").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HDel("test1", "111").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGet(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test2", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := cache.Get("test2").Result()
	if string(res) != "1" || err != nil {
		t.Fatal("get error")
	}
}

func TestDelString(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test3", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := cache.Get("test3").Result()
	if string(res) != "1" || err != nil {
		t.Fatal("get error")
	}

	intVal, err := cache.Del("test3").Result()
	if intVal != 1 || err != nil {
		t.Fatal("del string result error")
	}

	intVal, err = cache.Del("test3").Result()
	if intVal != 0 || err != nil {
		t.Fatal("del empty result error")
	}

	res, err = cache.Get("test3").Result()
	if string(res) != "" || err != nil {
		t.Fatal("get error")
	}
}

func TestExpire(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test2", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := cache.Get("test2").Result()
	if string(res) != "1" || err != nil {
		t.Fatal("get error")
	}
	resExpire, err := cache.Expire("test2", 1).Result()
	if resExpire != 1 && err != nil {
		t.Fatal("expire result error")
	}
	resExpire, err = cache.Expire("testNotExist", 1).Result()
	if resExpire != 0 && err != nil {
		t.Fatal("expire result error")
	}
	time.Sleep(time.Duration(2) * time.Second)
	resA, err := cache.Get("test2").Result()
	if string(resA) != "" || err != nil {
		t.Fatal("expire func error")
	}
}

func TestHSet(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := cache.HSet("key1", "field1", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if res != 1 {
		t.Fatal("hset1 result error")
	}
	res2, err := cache.HSet("key1", "field1", []byte("2")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if res2 != 0 {
		t.Fatal("hset2 result error")
	}
}

func TestHGet(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HSet("key1", "field1", []byte("100")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res2, err := cache.HGet("key1", "field1").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(res2) != "100" {
		t.Fatal("hget result error")
	}
}

func TestHDel(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HSet("key1", "field1", []byte("100")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res2, err := cache.HDel("key1", "field1").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if res2 != 1 {
		t.Fatal("Hdel result error")
	}
}

func TestHDelMulti(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HSet("key1", "field1", []byte("100")).Result()
	_, err = cache.HSet("key1", "field2", []byte("200")).Result()
	_, err = cache.HSet("key1", "field3", []byte("300")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res2, err := cache.HDel("key1", "field1", "field2", "field3", "notexist").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if res2 != 3 {
		t.Fatal("Del Multi result error")
	}
}

func TestDelMulti(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test1", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.HSet("key1", "field1", []byte("100")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	res2, err := cache.Del("key1", "test1", "notexist").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if res2 != 2 {
		t.Fatal("del multi result error")
	}

	getRes, err := cache.Get("test2").Result()
	if string(getRes) != "" || err != nil {
		t.Fatal("get error")
	}
	hgetRes, err := cache.HGet("key1", "field1").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(hgetRes) != "" {
		t.Fatal("hget result error")
	}
}

func TestLimit(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 1})
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test1", []byte("1")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = cache.Set("test1", []byte("2")).Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	// should have error
	_, err = cache.Set("test2", []byte("out of limit")).Result()
	if err == nil {
		t.Fatal("should have error,  but no error")
	}
}

func TestSAdd(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	res1, err := cache.SAdd("key1", "111", "222").Result()
	if err != nil || res1 != 2 {
		t.Fatal("res1 error")
	}
	res2, err := cache.SAdd("key1", "111", "222", "333").Result()
	if err != nil || res2 != 1 {
		t.Fatal("res2 error")
	}
	t.Log("Sadd: res1=", res1, ", res2=", res2)
}

func TestSIsMember(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 10})
	if err != nil {
		t.Fatal(err.Error())
	}
	res0, err := cache.SIsMember("key1", "111").Result()
	if err != nil || res0 != 0 {
		t.Fatal("res0 error")
	}
	res1, err := cache.SAdd("key1", "111", "222").Result()
	if err != nil || res1 != 2 {
		t.Fatal("res1 error")
	}
	res2, err := cache.SIsMember("key1", "111").Result()
	if err != nil || res2 != 1 {
		t.Fatal("res2 error")
	}
	res3, err := cache.SIsMember("key1", "333").Result()
	if err != nil || res3 != 0 {
		t.Fatal("res3 error")
	}
	cache.Del("key1")
	res4, err := cache.SIsMember("key1", "111").Result()
	if err != nil || res4 != 0 {
		t.Fatal("res4 error")
	}

	t.Log("Sadd: res0=", res0, ", res1=", res1, ", res2=", res2, ", res3=", res3, ", res4=", res4)
}

func TestGetBench(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{MaxSize: 175000})
	if err != nil {
		t.Fatal(err.Error())
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 175000; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			_, err = cache.Set(strconv.Itoa(count), []byte("1")).Result()
			if err != nil {
				t.Fatal(err.Error())
			}
			res, err := cache.Get(strconv.Itoa(count)).Result()
			if string(res) != "1" || err != nil {
				t.Fatal("res=", res)
			}
		}(i)
	}
	wg.Wait()
}

func TestExpireBench(t *testing.T) {
	cache, err := NewMemCache(&CacheConf{
		MaxSize:             1000,
		TtlPeriodMillSecond: 500,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			_, err = cache.Set(strconv.Itoa(count), []byte("1")).Result()
			if err != nil {
				t.Fatal(err.Error())
			}
			cache.Expire(strconv.Itoa(count), 1)
			time.Sleep(time.Second * 2)
			res, err := cache.Get(strconv.Itoa(count)).Result()
			if string(res) != "" || err != nil {
				t.Fatal("res=", res)
			}
		}(i)
	}
	wg.Wait()
}
