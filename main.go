package main

import (
	"log"

	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	"go.etcd.io/etcd/server/v3/embed"
)

func main() {
	cfg := embed.NewConfig()
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
