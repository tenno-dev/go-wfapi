package helper

import (
	"github.com/bitti09/go-wfapi/datasources"
	"net/http"
	"io/ioutil"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)



// Getjson - Getjson
func Getjson(uri string) string {
	cachedir := "./_cache"
	cache := diskcache.New(cachedir)
	tp := httpcache.NewTransport(cache)
	client := &http.Client{Transport: tp}
	resp, err := client.Get(uri)
	if err != nil {
		panic(err)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(bytes)
}