package store

import (
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/embed"
)

const (
	EtcdDataFolder = "etcd_data_"
)

func (s *Server) genConfig(c *embed.Config) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Failed to get hostname:", err)
	}
	c.Name = hostname
	if len(s.PeerURLs) > 0 {
		c.LPUrls = s.PeerURLs
		c.APUrls = s.PeerURLs
	}
	if len(s.ClientURLs) > 0 {
		c.LCUrls = s.ClientURLs
		c.ACUrls = s.ClientURLs
	}
	if len(s.Name) > 0 {
		c.Name = s.Name
	}
	if len(s.Cluster) > 0 {
		c.InitialCluster = s.Cluster
	} else {
		c.InitialCluster = c.InitialClusterFromName(c.Name)
	}
	c.Dir = EtcdDataFolder + c.Name
	c.LogPkgLevels = "all=FATAL"
}

func (s *Server) Start() error {
	var err error
	if s.etcd == nil {
		inCfg := embed.NewConfig()
		s.genConfig(inCfg)
		if err := inCfg.Validate(); err != nil {
			return err
		}
		s.etcd, err = embed.StartEtcd(inCfg)
	} else {
		s.etcd.Server.Start()
	}
	return err
}

func (s *Server) StartAndServe() error {
	if err := s.Start(); err != nil {
		return err
	}
	defer s.etcd.Close()
	select {
	case <-s.etcd.Server.ReadyNotify():
		log.Println("etcd server is ready!")
	case <-time.After(60 * time.Second):
		s.etcd.Server.Stop() // trigger a shutdown
		log.Println("etcd server took too long to start, stopping")
	}
	return <-s.etcd.Err()
}

func (s *Server) Stop() {
	if s.etcd != nil {
		s.etcd.Close()
	}
}
