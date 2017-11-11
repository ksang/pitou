package store

import (
	"context"
	"errors"

	"github.com/coreos/etcd/clientv3"
	"github.com/ksang/pitou/util"
)

var (
	ErrNotInit = errors.New("client no initialized")
)

func (c *Client) Init() error {
	if c.Server == nil {
		return errors.New("server is empty")
	}
	endpoints := util.UrlsToStrings(c.Server.ClientURLs)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: c.Timeout,
	})
	if err != nil {
		return err
	}
	c.cli = cli
	c.kv = clientv3.NewKV(cli)
	return nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) Put(key string, value string) error {
	if c.kv == nil {
		return ErrNotInit
	}
	op := clientv3.OpPut(key, value)
	_, err := c.kv.Do(context.TODO(), op)
	return err
}

func (c *Client) Get(key string) (map[string]string, error) {
	if c.kv == nil {
		return nil, ErrNotInit
	}
	resp, err := c.kv.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for _, ev := range resp.Kvs {
		ret[string(ev.Key)] = string(ev.Value)
	}
	return ret, nil
}

func (c *Client) GetSortedPrefix(key string) (map[string]string, error) {
	if c.kv == nil {
		return nil, ErrNotInit
	}
	resp, err := c.kv.Get(context.TODO(), key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for _, ev := range resp.Kvs {
		ret[string(ev.Key)] = string(ev.Value)
	}
	return ret, nil
}

func (c *Client) Del(key string) error {
	if c.kv == nil {
		return ErrNotInit
	}
	_, err := c.kv.Delete(context.TODO(), key)
	return err
}
