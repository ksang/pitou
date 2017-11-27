package puppet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/ksang/pitou/store"
	"github.com/ksang/pitou/util"
)

// SwitchREST is switch rest type of puppet includes internal used members
type SwitchREST struct {
	// Node is puppet configuration
	Node      *Puppet
	errCh     chan error
	quitCh    chan struct{}
	portCount int
	instance  string
	store     *store.Client
}

// NewSwitchREST returns a SwitchREST collector of puppet configuration provided
func NewSwitchREST(p Puppet) Collector {
	return &SwitchREST{
		Node:     &p,
		errCh:    make(chan error, 1),
		quitCh:   make(chan struct{}, 1),
		instance: util.RemoveScheme(p.Address),
		store:    p.Store,
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
					log.Println("collect error:", err)
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
	if err := s.collectPorts(s.portCount); err != nil {
		return err
	}
	return nil
}

func (s *SwitchREST) collectSystemResp() ([]byte, error) {
	if s.Node == nil {
		return []byte{}, nil
	}
	resp, err := http.Get(s.Node.Address + "/system/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *SwitchREST) updateSystemResp(sysResp []byte) error {
	sr := &SwitchRESTSystemResp{
		Data: SwitchRESTSystemRespData{},
	}
	if err := json.Unmarshal(sysResp, &sr); err != nil {
		return err
	}
	s.portCount = sr.Data.PortCount
	if s.store != nil {
		s.store.Put("/nodes/"+s.instance+"/port_count", strconv.FormatInt(int64(s.portCount), 10))
	}
	return nil
}

func (s *SwitchREST) collectPorts(portNum int) error {
	for i := 0; i < portNum; i++ {
		pdr, err := s.collectPortDetailResp(i)
		if err != nil {
			return err
		}
		if err := s.updatePortDetail(pdr); err != nil {
			return err
		}
	}
	return nil
}

func (s *SwitchREST) collectPortDetailResp(portId int) ([]byte, error) {
	if s.Node == nil {
		return []byte{}, nil
	}
	addr := fmt.Sprintf("%s/ports/%d/", s.Node.Address, portId)
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *SwitchREST) updatePortDetail(pdResp []byte) error {
	sr := &SwitchRESTPortDetailResp{}
	if err := json.Unmarshal([]byte(pdResp), &sr); err != nil {
		return err
	}
	if len(sr.Data.Error) > 0 {
		return errors.New(sr.Data.Error)
	}
	prefix := fmt.Sprintf("/nodes/%s/ports/%d/", s.instance, sr.Data.PortID)
	if s.store != nil {
		// port detail info
		s.store.Put(prefix+"admin_speed", sr.Data.AdminSpeed)
		s.store.Put(prefix+"admin_state", sr.Data.AdminState)
		s.store.Put(prefix+"oper_speed", strconv.FormatInt(int64(sr.Data.OperationSpeed), 10))
		s.store.Put(prefix+"oper_state", sr.Data.OperationState)
		s.store.Put(prefix+"mtu", strconv.FormatInt(int64(sr.Data.MTU), 10))
		s.store.Put(prefix+"label_port", strconv.FormatInt(int64(sr.Data.PortLable), 10))
		// port counters
		counters := reflect.Indirect(reflect.ValueOf(sr.Data.Counters))
		for i := 0; i < counters.NumField(); i++ {
			if counters.Field(i).Kind() == reflect.Int {
				value := counters.Field(i).Interface().(int)
				fmt.Printf("Counters field: %s, value: %v\n", counters.Type().Field(i).Name, value)
				s.store.Put(prefix+"counters/"+counters.Type().Field(i).Name, strconv.FormatInt(int64(value), 10))
			}
		}
	}
	return nil
}
