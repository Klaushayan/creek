package brook

import (
	"encoding/json"
	"net/http"
)

func (f *Firewall) ServeConnectedIPs(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(f.connectedIPs)
	if err != nil {
		f.log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(j)
}

func (f *Firewall) ServeBlockedIPs(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(f.blockedIPs)
	if err != nil {
		f.log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(j)
}

func (f *Firewall) ServeConfig(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(
		struct {
			MaxConnections int `json:"max_connections"`
			BlockPeriod    int `json:"block_period"`
			ConnectionCool int `json:"connection_cool"`
		}{
			MaxConnections: f.MaxConnections,
			BlockPeriod:    int(f.BlockPeriod / 60),
			ConnectionCool: int(f.ConnectionCooldown / 60),
		},
	)
	if err != nil {
		f.log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(j)
}

func (f *Firewall) ChangeConfig(w http.ResponseWriter, r *http.Request) {
	var config struct {
		MaxConnections int `json:"max_connections"`
		BlockPeriod    int `json:"block_period"`
		ConnectionCool int `json:"connection_cool"`
	}
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		f.log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f.MaxConnections = config.MaxConnections
	f.BlockPeriod = int64(config.BlockPeriod * 60)
	f.ConnectionCooldown = int64(config.ConnectionCool * 60)
}

func (f *Firewall) BlockIPList(w http.ResponseWriter, r *http.Request) {
	var ips []string
	err := json.NewDecoder(r.Body).Decode(&ips)
	if err != nil {
		f.log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var bl []IP
	for _, ip := range ips {
		bl = append(bl, NewIP(ip))
	}
	f.log("Blocking IPs by request: ", IPList(bl).String())
	f.Block(bl...)
}
