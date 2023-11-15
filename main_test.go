package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"

	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

type EmbednsTestSuite struct {
	suite.Suite
}

func (suite *EmbednsTestSuite) SetupTest() {
	go func() {
		cfg := embed.NewConfig()
		cfg.LogLevel = "debug"
		cfg.LogOutputs = []string{"./etcd.log"}
		cfg.Dir = "default.etcd"
		_, err := embed.StartEtcd(cfg)

		suite.Nil(err)
		coremain.Run()
	}()
}

func (suite *EmbednsTestSuite) TestARecord() {
	ctx := context.Background()
	etcdAddresses := "localhost:2379"

	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   strings.Split(etcdAddresses, ","),
			DialTimeout: 3 * time.Second,
			Context:     ctx,
		},
	)
	suite.Nil(err)

	defer cli.Close()

	hostname := namesgenerator.GetRandomName(12)
	suite.T().Logf("hostname is: %s", hostname)
	key := fmt.Sprintf("/coredns/%s/x1", hostname)
	hostip := fmt.Sprintf("10.233.%d.%d", rand.Intn(128), rand.Intn(100))

	val := fmt.Sprintf(`{"host": "%s","ttl":60}`, hostip)
	ctx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer func() {
		cancelFunc()
	}()

	_, err = cli.Put(ctx, key, val)
	suite.Nil(err)

	// nolint
	defer cli.Delete(ctx, key)

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "localhost:53")
		},
	}

	ip, err := r.LookupHost(context.Background(), hostname)
	suite.Nil(err)
	suite.Equal(ip[0], hostip)
}

func (suite *EmbednsTestSuite) TearDownSuite() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(EmbednsTestSuite))
}
