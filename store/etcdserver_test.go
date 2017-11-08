package store

import (
	"testing"
	"time"

	"github.com/coreos/etcd/embed"
	"github.com/ksang/pitou/util"
)

func TestNewConfig(t *testing.T) {
	cfg := embed.NewConfig()
	t.Logf("%#v", *cfg)
}

func TestStartServer(t *testing.T) {
	pu, _ := util.StringToUrls("http://127.0.0.1:1234")
	cu, _ := util.StringToUrls("http://127.0.0.1:1235")
	ss := []Server{
		Server{},
		Server{
			Name:       "test",
			PeerURLs:   pu,
			ClientURLs: cu,
		},
	}
	for _, s := range ss {
		if err := s.Start(); err != nil {
			t.Errorf("%s", err)
			return
		}
		defer s.Stop()
		select {
		case <-s.etcd.Server.ReadyNotify():
			t.Logf("etcd Server is ready!")
		case <-time.After(60 * time.Second):
			s.etcd.Server.Stop() // trigger a shutdown
			t.Errorf("Server took too long to start!")
		}

		select {
		case err := <-s.etcd.Err():
			t.Errorf("etcd server error: %s", err)
		case <-time.After(10 * time.Second):
			t.Logf("etcd start server ok")
		}
		s.etcd.Close()
	}
}
