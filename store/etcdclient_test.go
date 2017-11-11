package store

import (
	"reflect"
	"testing"
	"time"

	"github.com/coreos/pkg/capnslog"
	"github.com/ksang/pitou/util"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
}

func TestInitEtcdClient(t *testing.T) {
	cu, _ := util.StringToUrls("http://127.0.0.1:1235")
	cs := []Client{
		Client{},
		Client{
			Server: &Server{},
		},
		Client{
			Server: &Server{
				ClientURLs: cu,
			},
		},
	}
	errs := []bool{true, true, false}
	for i, c := range cs {
		if err := c.Init(); err != nil {
			if errs[i] == false {
				t.Errorf("Client Init test case #%d failed, %#v, error %v",
					i, c, err)
			}
		} else {
			if errs[i] == true {
				t.Errorf("Client Init test case #%d failed, %#v, expecting error",
					i, c)
			}
		}
	}
}

func TestClientOps(t *testing.T) {
	tests := []struct {
		op       string
		key      string
		value    string
		expected map[string]string
	}{
		{
			"put",
			"/test/1",
			"value1",
			nil,
		},
		{
			"put",
			"/test/1/2",
			"value2",
			nil,
		},
		{
			"get",
			"/test/1",
			"",
			map[string]string{"/test/1": "value1"},
		},
		{
			"del",
			"/test/1",
			"",
			nil,
		},
		{
			"get",
			"/test/1",
			"",
			map[string]string{},
		},
	}

	pu, _ := util.StringToUrls("http://127.0.0.1:2234")
	cu, _ := util.StringToUrls("http://127.0.0.1:2235")
	s := Server{
		Name:       "test",
		PeerURLs:   pu,
		ClientURLs: cu,
	}
	defer s.Stop()
	if err := s.Start(); err != nil {
		t.Errorf("%s", err)
		return
	}
	select {
	case <-s.etcd.Server.ReadyNotify():
		t.Logf("etcd Server is ready!")
	case <-time.After(10 * time.Second):
		s.Stop()
		t.Errorf("Server took too long to start!")
	}
	c := Client{
		Server:  &s,
		Timeout: time.Second,
	}
	c.Init()
	for i, te := range tests {
		switch te.op {
		case "put":
			if err := c.Put(te.key, te.value); err != nil {
				t.Errorf("test case #%d failed, error: %v", i, err)
			}
		case "get":
			r, err := c.Get(te.key)
			if err != nil {
				t.Errorf("test case #%d failed, error: %v", i, err)
			}
			if !reflect.DeepEqual(r, te.expected) {
				t.Errorf("test case #%d check result error, key %s, expected %s, actual %s",
					i, te.key, te.expected, r)
			}
		case "del":
			if err := c.Del(te.key); err != nil {
				t.Errorf("test case #%d failed, error: %v", i, err)
			}
		}

		root, err := c.GetSortedPrefix("/")
		if err != nil {
			t.Errorf("test case #%d get failed, error: %v", i, err)
			continue
		}
		t.Logf("test case #%d, root: %v", i, root)
	}
}
