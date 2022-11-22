package brook

import (
	"log"
	"time"
)

type Firewall struct {
	MaxConnections     int
	BlockPeriod        int64
	ConnectionCooldown int64

	connectedIPs []IP
	blockedIPs   []IP
}

type IP struct {
	Address         string
	New             bool
	FirstConnection int64
	LastConnection  int64
}

func NewFirewall(maxConnections, blockPeriod, connectionCooldown int) *Firewall {
	return &Firewall{
		MaxConnections:     maxConnections,
		BlockPeriod:        int64(blockPeriod * 60),
		ConnectionCooldown: int64(connectionCooldown * 60),
	}
}

func (f *Firewall) IsBlocked(ip string) bool {
	for _, b := range f.blockedIPs {
		if b.Address == ip {
			return true
		}
	}
	return false
}

func (f *Firewall) IsConnected(ip string) bool {
	for _, c := range f.connectedIPs {
		if c.Address == ip {
			return true
		}
	}
	return false
}

func (f *Firewall) Block(ip ...IP) {
	f.blockedIPs = append(f.blockedIPs, ip...)
}

func (f *Firewall) Unblock(ip string) {
	for i, b := range f.blockedIPs {
		if b.Address == ip {
			f.blockedIPs = append(f.blockedIPs[:i], f.blockedIPs[i+1:]...)
			return
		}
	}
}

func (f *Firewall) updateLastConnection(ip string) {
	for i, c := range f.connectedIPs {
		if c.Address == ip {
			f.connectedIPs[i].LastConnection = time.Now().Unix()
			return
		}
	}
}

func (f *Firewall) coolDownCheck() {
	if f.MaxConnections == 0 {
		return
	}
	var now = time.Now().Unix()
	cs := f.connectedIPs
	for i, c := range cs {
		if now-c.LastConnection >= f.ConnectionCooldown {
			f.connectedIPs = append(f.connectedIPs[:i], f.connectedIPs[i+1:]...)
		}
		if c.New && now-c.FirstConnection >= f.ConnectionCooldown {
			f.connectedIPs[i].New = false
		}
	}
}

func (f *Firewall) blockCheck() {
	if len(f.connectedIPs) > f.MaxConnections {
		tb := f.connectedIPs[f.MaxConnections:]
		f.Block(tb...)
		log.Println("Blocked IPs:", tb)
		f.connectedIPs = f.connectedIPs[:f.MaxConnections]
	}
}

func (f *Firewall) unblockCheck() {
	var now = time.Now().Unix()
	bs := f.blockedIPs
	for i, b := range bs {
		if now-b.LastConnection >= f.BlockPeriod {
			f.blockedIPs = append(f.blockedIPs[:i], f.blockedIPs[i+1:]...)
		}
	}
}

func (f *Firewall) Verify(ip string) bool {
	if f.MaxConnections == 0 {
		return true
	}
	if len(f.connectedIPs) >= f.MaxConnections {
		return false
	}
	if f.IsBlocked(ip) {
		return false
	}
	if f.IsConnected(ip) {
		f.updateLastConnection(ip)
		return true
	}
	f.connectedIPs = append(f.connectedIPs, IP{Address: ip, FirstConnection: time.Now().Unix(), New: true})
	return true
}
