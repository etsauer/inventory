package database

import (
	"context"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const timeout = 5

var cli clientv3.Client
var kvc clientv3.KV

func init() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		println("Failed to create etcd client" + err.Error())
	}
	defer cli.Close()

	kvc = clientv3.NewKV(cli)
}

func get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	gresp, err := kvc.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	return gresp.Kvs[0].Value, nil
}

func put(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_, err := cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		// handle error!
		return err
	}
	return nil
}
