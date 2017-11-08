/*
Package puppet provides functionalities interacting with SDN nodes
*/
package puppet

import "time"

/*
Collector is the abstraction of collecting metrics and telemetry from SDN nodes.
*/
type Collector interface {
	Start() chan error
	Stop() error
}

// NodeType is the enum type of puppet node
type NodeType uint8

const (
	// NodeSwitchREST REST service running on Mellanox SwitchX
	NodeSwitchREST NodeType = iota
	// NodeUNKOWN type
	NodeUNKOWN
)

// Puppet describes properties of a puppet node
type Puppet struct {
	Type     NodeType
	Address  string
	Interval time.Duration
}
