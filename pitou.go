package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/ksang/pitou/puppet"
	"github.com/ksang/pitou/store"
	"github.com/ksang/pitou/util"
	"github.com/olebedev/config"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "d", "pitou.conf", "configuration file location")
}

func StartStore(cfg *config.Config) *store.Client {
	storeCfg, err := cfg.Map("store")
	if err != nil {
		log.Fatal("Failed to get store config", err)
	}
	var (
		pua      []url.URL
		cua      []url.URL
		namea    string
		clustera string
		timeouta = time.Second
	)
	pu, ok := storeCfg["PeerURL"].(string)
	if ok {
		pua, err = util.StringToUrls(pu)
		if err != nil {
			log.Fatal("failed to parse peer url:", err)
		}
	}
	cu, ok := storeCfg["ClientURL"].(string)
	if ok {
		cua, err = util.StringToUrls(cu)
		if err != nil {
			log.Fatal("failed to parse peer url:", err)
		}
	}
	name, ok := storeCfg["Name"].(string)
	if ok {
		namea = name
	}
	cluster, ok := storeCfg["Cluster"].(string)
	if ok {
		clustera = cluster
	}
	s := &store.Server{
		Name:       namea,
		PeerURLs:   pua,
		ClientURLs: cua,
		Cluster:    clustera,
	}
	if err := s.Start(); err != nil {
		log.Fatal("Failed to start store:", err)
	}
	to, ok := storeCfg["Timeout"].(string)
	if ok {
		timeouta, err = time.ParseDuration(to)
		if err != nil {
			log.Fatal("Failed to parse store timeout:", err)
		}
	}
	c := &store.Client{
		Server:  s,
		Timeout: timeouta,
	}
	if err := c.Init(); err != nil {
		log.Fatal("Failed to init store client:", err)
	}
	return c
}

func StartPuppetMgr(cfg *config.Config, cli *store.Client) {
	mgr := puppet.NewManager()
	puppets, err := cfg.List("puppets")
	if err != nil {
		log.Println("failed to get puppets from config")
		return
	}
	for i, p := range puppets {
		interval := puppet.DefaultInterval
		addr, ok := p.(map[string]interface{})["Address"]
		if !ok {
			log.Println("no address configured for puppet #%d", i)
			continue
		}
		inter, ok := p.(map[string]interface{})["Interval"]
		if ok {
			interval, err = time.ParseDuration(inter.(string))
			if err != nil {
				interval = puppet.DefaultInterval
			}
		}
		c := puppet.NewSwitchREST(puppet.Puppet{
			Address:  addr.(string),
			Interval: interval,
			Store:    cli,
		})
		mgr.Add(c)
	}
}
func main() {
	flag.Parse()
	cfgs, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.ParseYaml(string(cfgs))
	if err != nil {
		log.Fatal("Failed to parse yaml config, err:", err)
	}
	cli := StartStore(cfg)
	StartPuppetMgr(cfg, cli)
	time.Sleep(999 * time.Second)
}
