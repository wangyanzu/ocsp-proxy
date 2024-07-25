package cache

import (
	"net/http"
	"sync"
	"time"

	"github.com/wangyanzu/ocsp-proxy/internal/config"
)

var lock sync.RWMutex
var cache map[string]*ocspCache

type RespCache struct {
	Header http.Header
	Body   []byte
}

type ocspCache struct {
	Expireat time.Time
	Cache    *RespCache
}

func init() {
	cache = make(map[string]*ocspCache, 1)
}

func Get(key string) (resp *RespCache, ok bool, expired bool) {
	lock.RLock()
	defer lock.RUnlock()
	oc, ok := cache[key]
	if ok {
		expired = time.Now().After(oc.Expireat)
	} else {
		expired = true
		return nil, ok, expired
	}
	resp = oc.Cache
	return
}

func Set(key string, rc *RespCache) {
	lock.Lock()
	defer lock.Unlock()
	cache[key] = &ocspCache{
		Expireat: time.Now().Add(time.Duration(config.Interval) * time.Second),
		Cache:    rc,
	}
}
