package cache

import "time"

func volatileRange(db *MemCacheDB) {
	delRate := 1.0
	currentTime := time.Now()

	for delRate > 0.25 {
		delKey := make([]string, 0, 100)
		allCount := 1
		for key, value := range db.ttl {
			if allCount >= 100 {
				break
			}
			allCount++
			if currentTime.After(value) {
				delKey = append(delKey, key)
			}
		}
		delRate = float64(len(delKey)) / float64(allCount)
		for _, key := range delKey {
			db.delKey(key, true)
		}
	}
}
