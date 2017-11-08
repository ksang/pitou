package main

import "flag"

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "d", "pitou.conf", "configuration file location.")
}

func main() {
	flag.Parse()
}
