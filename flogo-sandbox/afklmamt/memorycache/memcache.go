package memorychache

import (
	"sync"
	"time"
)

type CachedObject struct {
	created time.Time
	expires time.Time
	value   interface{}
}

type CachedHTTPResponse struct {
	code    int
	headers map[string]string
	body    interface{}
}

var singletonCache map[string]CachedObject
var mu sync.Mutex

/**
 * Performs the lazy initialization of caching of the map.
 */
func locateCache() map[string]CachedObject {
	if singletonCache == nil {
		mu.Lock()
		defer mu.Unlock()

		if singletonCache == nil {
			singletonCache = make(map[string]CachedObject)
		}
	}

	return singletonCache
}

func Get(key string) interface{} {
	cacheRec, ok := locateCache()[key]
	if ok {
		return cacheRec.value
	} else {
		return nil
	}
}

func Set(key string, value interface{}) {
	locateCache()[key] = CachedObject{time.Now(), nil, value}
}
