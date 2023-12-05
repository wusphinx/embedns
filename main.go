package main

import (
	"log"
	"net/url"

	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	"go.etcd.io/etcd/server/v3/embed"
)

func main() {
	cfg := embed.NewConfig()
	lp := "http://0.0.0.0:2379"
	lcurl, _ := url.Parse(lp)
	cfg.ListenClientUrls = []url.URL{*lcurl}
	cfg.LogLevel = "debug"
	cfg.LogOutputs = []string{"./etcd.log"}
	cfg.Dir = "default.etcd"
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer e.Close()
	coremain.Run()
}
