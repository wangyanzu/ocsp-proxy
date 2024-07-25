package config

import (
	"flag"
	"log"
	"os"
)

var Addr string
var OcspHost string
var Interval uint

func init() {
	flag.StringVar(&OcspHost, "ocsphost", "", "OCSP server to proxy requests to")
	flag.StringVar(&Addr, "http", ":8080", "HTTP host:port to listen to")
	flag.UintVar(&Interval, "interval", 1800, "cache refresh time in second")
	flag.Parse()
	if OcspHost == "" {
		OcspHost = os.Getenv("OCSP_HOST")
	}
	if OcspHost == "" {
		log.Fatal("need ocsphost parameter")
	}
}
