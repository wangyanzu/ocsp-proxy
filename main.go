package main

import (
	"context"
	"log"
	"net/http"

	"github.com/wangyanzu/ocsp-proxy/internal/cache"
	"github.com/wangyanzu/ocsp-proxy/internal/config"
	"github.com/wangyanzu/ocsp-proxy/internal/ocsp"
)

func main() {
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		cacheKey := request.RequestURI
		log.Println(request.RemoteAddr, "Requesting", cacheKey)
		var err error
		rc, ok, expired := cache.Get(cacheKey)
		// if no ocsp exist
		if ok {
			log.Println("hit cache: ", cacheKey)
		} else {
			rc, err = ocsp.RequestResponserAndCache(request)
			if err != nil {
				log.Println(err)
				responseWriter.WriteHeader(500)
				return
			}
			expired = false
		}
		for name, values := range rc.Header {
			responseWriter.Header()[name] = values
		}
		_, err = responseWriter.Write(rc.Body)
		if err != nil {
			log.Println(err)
			return
		}
		if expired {
			// set context nil
			log.Println("cache is expired", cacheKey)
			go ocsp.RequestResponserAndCache(request.WithContext(context.TODO()))
		}
	})
	http.ListenAndServe(config.Addr, nil)
}
