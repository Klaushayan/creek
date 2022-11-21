package brook

import "time"

type Firewall struct {
	MaxConnections     int
	BlockPeriod        int
	ConnectionCooldown int

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
		BlockPeriod:        blockPeriod,
		ConnectionCooldown: connectionCooldown,
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

func (f *Firewall) Block(ip string) {
	f.blockedIPs = append(f.blockedIPs, IP{Address: ip})
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