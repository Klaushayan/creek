package brook

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
)

type CreekManager struct {
	WSServers []WSServerJSON
	Relays    []RelayJSON
	Port	  int
}

type WSServerJSON struct {
	Addr               string   `json:"addr"`
	Password           string   `json:"password"`
	Domain             string   `json:"domain"`
	Path               string   `json:"path"`
	TCPTimeout         int      `json:"tcp_timeout"`
	UDPTimeout         int      `json:"udp_timeout"`
	BlockDomainList    string   `json:"block_domain_list"`
	BlockCIDR4List     string   `json:"block_cidr4_list"`
	BlockCIDR6List     string   `json:"block_cidr6_list"`
	UpdateListInterval int64    `json:"update_list_interval"`
	BlockGeoIP         []string `json:"block_geoip"`
	MaxClients         int      `json:"max_clients"`
	ConnectionCooldown int      `json:"connection_cooldown"`
	BlockTimeout       int      `json:"block_timeout"`
	Object			   *WSServer
}


type RelayJSON struct {
	Server     string `json:"server"`
	From       string `json:"from"`
	To         string `json:"to"`
	Remote     string `json:"remote"`
	Password   string `json:"password"`
	TCPTimeout int    `json:"tcp_timeout"`
	UDPTimeout int    `json:"udp_timeout"`
	Object	   *Map
}

func (cm *CreekManager) NewCreekManager(port int) (*CreekManager, error) {
	cm.Port = port
	return cm, nil
}

func (cm *CreekManager) AddWSServer(w http.ResponseWriter, r *http.Request) {
	var wsj WSServerJSON
	err := json.NewDecoder(r.Body).Decode(&wsj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = cm.NewWSServer(&wsj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (cm *CreekManager) AddRelay(w http.ResponseWriter, r *http.Request) {
	var rj RelayJSON
	err := json.NewDecoder(r.Body).Decode(&rj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = cm.NewMap(&rj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (cm *CreekManager) NewWSServer(wsj *WSServerJSON) (*WSServer, error) {

	if wsj.Addr == "" || wsj.Password == "" {
		return nil, errors.New("addr or password is empty")
	}

	if wsj.BlockDomainList != "" && !strings.HasPrefix(wsj.BlockDomainList, "http://") && !strings.HasPrefix(wsj.BlockDomainList, "https://") && !filepath.IsAbs(wsj.BlockDomainList) {
		return nil, errors.New("blockDomainList must be with absolute path")
	}
	if wsj.BlockCIDR4List != "" && !strings.HasPrefix(wsj.BlockCIDR4List, "http://") && !strings.HasPrefix(wsj.BlockCIDR4List, "https://") && !filepath.IsAbs(wsj.BlockCIDR4List) {
		return nil, errors.New("blockCIDR4List must be with absolute path")
	}
	if wsj.BlockCIDR6List != "" && !strings.HasPrefix(wsj.BlockCIDR6List, "http://") && !strings.HasPrefix(wsj.BlockCIDR6List, "https://") && !filepath.IsAbs(wsj.BlockCIDR6List) {
		return nil, errors.New("blockCIDR6List must be with absolute path")
	}

	ws, err := NewWSServer(wsj.Addr, wsj.Password, wsj.Domain, wsj.Path, wsj.TCPTimeout, wsj.UDPTimeout, wsj.BlockDomainList, wsj.BlockCIDR4List, wsj.BlockCIDR6List, wsj.UpdateListInterval, wsj.BlockGeoIP, wsj.MaxClients, wsj.ConnectionCooldown, wsj.BlockTimeout)
	if err != nil {
		return nil, err
	}

	wsj.Object = ws

	cm.WSServers = append(cm.WSServers, *wsj)
	return ws, nil
}

func (cm *CreekManager) NewMap(rj *RelayJSON) (*Map, error) {
	if rj.From == "" || rj.To == "" || rj.Server == "" || rj.Password == "" {
		return nil, errors.New("from or to or server or password is empty")
	}

	m, err := NewMap(rj.From, rj.To, rj.Remote, rj.Password, rj.TCPTimeout, rj.UDPTimeout)
	if err != nil {
		return nil, err
	}

	rj.Object = m

	cm.Relays = append(cm.Relays, *rj)
	return m, nil
}