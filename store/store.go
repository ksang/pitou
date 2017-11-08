/*
Package store provides distributed KV store facility
*/
package store

import (
	"net/url"

	"github.com/coreos/etcd/embed"
)

type Server struct {
	// URL for peer listen and advertise, by default would be http://localhost:2380
	PeerURLs []url.URL
	// URL for client listen and advertise , by default would be http://localhost:2379
	ClientURLs []url.URL
	// Name of local node, will use hostname if not provided
	Name string
	// Initial Cluster string, by default it would be local node only
	// example: etcd0=http://1.1.1.1:2380,etcd1=http://2.2.2.2:2380
	Cluster string
	etcd    *embed.Etcd
}
