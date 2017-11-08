package puppet

import (
	"encoding/json"
	"errors"
	"time"
)

// SwitchREST is switch rest type of puppet includes internal used members
type SwitchREST struct {
	// Node is puppet configuration
	Node      Puppet
	errCh     chan error
	quitCh    chan struct{}
	portCount int
}

// NewSwitchREST returns a SwitchREST collector of puppet configuration provided
func NewSwitchREST(p Puppet) Collector {
	return &SwitchREST{
		Node:   p,
		errCh:  make(chan error, 1),
		quitCh: make(chan struct{}, 1),
	}
}

// Start collector
func (s *SwitchREST) Start() chan error {
	go func() {
		defer func() {
			close(s.errCh)
		}()
		for {
			select {
			case <-time.After(s.Node.Interval):
				err := s.collect()
				if err != nil {
					s.errCh <- err
				}
			case <-s.quitCh:
				s.errCh <- errors.New("stopped")
				return
			}
		}
	}()
	return s.errCh
}

// Stop collector
func (s *SwitchREST) Stop() error {
	s.quitCh <- struct{}{}
	return nil
}

// SwitchRESTSystemRespData defines SwitchREST system reponse data field
type SwitchRESTSystemRespData struct {
	Architecture string `json:"architecture"`
	Error        string `json:"error"`
	PortCount    int    `json:"port_count"`
}

// SwitchRESTSystemResp defines SwitchREST system reponse
type SwitchRESTSystemResp struct {
	Data      SwitchRESTSystemRespData `json:"data"`
	Directory string                   `json:"directory"`
	Timestamp float32                  `json:"timestamp"`
}

// SwitchRESTPortDetailRespCounters defines SwitchREST port detail response counters field
type SwitchRESTPortDetailRespCounters struct {
	RxBroadcast   int `json:"port_rx_broadcast"`
	RxFcsErrors   int `json:"port_rx_fcs_errors"`
	RxFrames      int `json:"port_rx_frames"`
	RxJumbo       int `json:"port_rx_jumbo"`
	RxMulticast   int `json:"port_rx_multicast"`
	RxNoBuffer    int `json:"port_rx_no_buffer"`
	RxOctets      int `json:"port_rx_octets"`
	RxOtherErrors int `json:"port_rx_other_errors"`
	RxRunt        int `json:"port_rx_runt"`
	RxUnicast     int `json:"port_rx_unicast"`
	TxBroadcast   int `json:"port_tx_broadcast"`
	TxFrames      int `json:"port_tx_frames"`
	TxJumbo       int `json:"port_tx_jumbo"`
	TxMulticast   int `json:"port_tx_multicast"`
	TxOctets      int `json:"port_tx_octets"`
	TxUnicast     int `json:"port_tx_unicast"`
	TxErrors      int `json:"port_tx_errors"`
}

// SwitchRESTPortDetailRespData defines SwitchREST port detail response data field
type SwitchRESTPortDetailRespData struct {
	AdminSpeed     string                           `json:"admin_speed"`
	AdminState     string                           `json:"admin_state"`
	ModuleState    string                           `json:"module_state"`
	MTU            int                              `json:"mtu"`
	OperationSpeed int                              `json:"oper_speed"`
	OperationState string                           `json:"oper_state"`
	PortID         int                              `json:"port_id"`
	PVID           int                              `json:"pvid"`
	Error          string                           `json:"error"`
	PortLable      int                              `json:"label_port"`
	PortLocal      int                              `json:"local_port"`
	PortLog        string                           `json:"log_port"`
	Counters       SwitchRESTPortDetailRespCounters `json:"counters"`
}

// SwitchRESTPortDetailResp defines SwitchREST port detail message
type SwitchRESTPortDetailResp struct {
	Data      SwitchRESTPortDetailRespData `json:"data"`
	Directory string                       `json:"directory"`
	Timestamp float32                      `json:"timestamp"`
}

func (s *SwitchREST) collect() error {
	sysResp, err := s.collectSystemResp()
	if err != nil {
		return err
	}
	if err := s.updateSystemResp(sysResp); err != nil {
		return err
	}
	return nil
}

func (s *SwitchREST) collectSystemResp() ([]byte, error) {
	return []byte{}, nil
}

func (s *SwitchREST) updateSystemResp(sysResp []byte) error {
	sr := &SwitchRESTSystemResp{
		Data: SwitchRESTSystemRespData{},
	}
	if err := json.Unmarshal(sysResp, &sr); err != nil {
		return err
	}
	s.portCount = sr.Data.PortCount
	return nil
}
