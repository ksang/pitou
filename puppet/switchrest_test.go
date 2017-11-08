package puppet

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	srSystemResp = `{
    "data": {
        "architecture": "switchx",
        "error": "",
        "port_count": 35
    },
    "directory": "system",
    "timestamp": 1504174356.99182
}`
	srPortsResp = `{
    "data": {
        "error": "",
        "port_count": 35
    },
    "directory": "ports",
    "timestamp": 1504174518.370088
}`
	srPortsDetailResp = `{
    "data": {
        "admin_speed": "mode_100GB_CR4|mode_100GB_SR4|mode_100GB_LR4_ER4",
        "admin_state": "UP",
        "counters": {
            "port_rx_broadcast": 29,
            "port_rx_fcs_errors": 0,
            "port_rx_frames": 638,
            "port_rx_jumbo": 0,
            "port_rx_multicast": 499,
            "port_rx_no_buffer": 0,
            "port_rx_octets": 78107,
            "port_rx_other_errors": 0,
            "port_rx_runt": 0,
            "port_rx_unicast": 0,
            "port_tx_broadcast": 402200,
            "port_tx_errors": 0,
            "port_tx_frames": 3024557,
            "port_tx_jumbo": 0,
            "port_tx_multicast": 2622012,
            "port_tx_octets": 292960408,
            "port_tx_unicast": 345
        },
        "error": "",
        "label_port": 1,
        "local_port": 61,
        "log_port": "0x13d00",
        "module_state": "PLUGGED",
        "mtu": 9238,
        "oper_speed": 16,
        "oper_state": "UP",
        "port_id": 32,
        "pvid": 1
    },
    "directory": "ports",
    "timestamp": 1504174572.226944
}`
	srTeleSamplingResp = `{
    "data": {
        "error": "",
        "histogram": {
            "0": {
                "legend": "<2976",
                "value": 1080681
            },
            "1": {
                "legend": "27552",
                "value": 0
            },
            "2": {
                "legend": "52128",
                "value": 0
            },
            "3": {
                "legend": "76704",
                "value": 0
            },
            "4": {
                "legend": "101280",
                "value": 0
            },
            "5": {
                "legend": "125856",
                "value": 0
            },
            "6": {
                "legend": "150432",
                "value": 0
            },
            "7": {
                "legend": "175008",
                "value": 0
            },
            "8": {
                "legend": "199584",
                "value": 0
            },
            "9": {
                "legend": "199584<",
                "value": 0
            }
        },
        "label_port": 2,
        "log_port": "0x13f00",
        "port_id": 33,
        "traffic_class": 0
    },
    "directory": "telemetry",
    "timestamp": 1504174629.824418
}`
)

func TestUpdateSystemResp(t *testing.T) {
	s := SwitchREST{}
	err := s.updateSystemResp([]byte(srSystemResp))
	if err != nil {
		t.Errorf("%v", err)
	}
	if s.portCount != 35 {
		t.Errorf("port count incorrect, expected: %d, actual: %d", 35, s.portCount)
	}
}

func TestUpdatePortDetailResp(t *testing.T) {
	sr := &SwitchRESTPortDetailResp{}
	if err := json.Unmarshal([]byte(srPortsDetailResp), &sr); err != nil {
		t.Error(err)
	}
	fmt.Printf("SR: %#v", sr)
}
