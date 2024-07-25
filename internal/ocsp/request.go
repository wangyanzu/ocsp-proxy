package ocsp

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wangyanzu/ocsp-proxy/internal/cache"
	"github.com/wangyanzu/ocsp-proxy/internal/config"
)

func RequestResponserAndCache(request *http.Request) (rc *cache.RespCache, err error) {
	cacheKey := request.RequestURI
	request.URL.Scheme = "http"
	request.URL.Host = config.OcspHost
	request.Host = request.URL.Host

	request.RequestURI = ""
	c := &http.Client{Timeout: time.Second * 30}
	resp, err := c.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	// cache ocsp response
	var body bytes.Buffer
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Println(err)
	}
	rc = &cache.RespCache{
		Header: resp.Header,
		Body:   body.Bytes(),
	}
	// log.Println("set cache: ", cacheKey)
	log.Printf("[%d] set cache %s \nBody: %s \n", resp.StatusCode, cacheKey, body.Bytes())
	cache.Set(cacheKey, rc)

	resp.Body = io.NopCloser(&body)
	return
}
