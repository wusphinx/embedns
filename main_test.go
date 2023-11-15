package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestARecord(t *testing.T) {
	ctx := context.Background()
	etcdAddresses := "localhost:2379"

	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   strings.Split(etcdAddresses, ","),
			DialTimeout: 3 * time.Second,
			Context:     ctx,
		},
	)
	assert.Nil(t, err)

	defer cli.Close()

	hostname := namesgenerator.GetRandomName(12)
	t.Logf("hostname is: %s", hostname)
	key := fmt.Sprintf("/coredns/%s/x1", hostname)
	hostip := fmt.Sprintf("10.233.%d.%d", rand.Intn(128), rand.Intn(100))

	val := fmt.Sprintf(`{"host": "%s","ttl":60}`, hostip)
	ctx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer func() {
		cancelFunc()
	}()

	_, err = cli.Put(ctx, key, val)
	assert.Nil(t, err)

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
	assert.Nil(t, err)
	assert.Equal(t, ip[0], hostip)
}
