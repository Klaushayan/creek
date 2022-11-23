package brook

import (
	"log"
	"os"
	"strings"
	"time"
)

type Firewall struct {
	MaxConnections     int
	BlockPeriod        int64
	ConnectionCooldown int64

	connectedIPs []IP
	blockedIPs   []IP
	interval	 *Job
	logFile   	 *Job
}

type IP struct {
	Address         string
	New             bool
	FirstConnection int64
	LastConnection  int64
}

type IPList []IP

func (ip IP) String() string {
	return ip.Address
}

func (ips IPList) String() string {
	var sb strings.Builder
	for _, ip := range ips {
		sb.WriteString(ip.Address)
		sb.WriteString(", ")
	}
	return sb.String()
}

func NewFirewall(maxConnections, blockPeriod, connectionCooldown int) *Firewall {

	o, err := os.OpenFile("creek.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		o = os.Stdout
	}

	defer o.Close()

	f := &Firewall{
		MaxConnections:     maxConnections,
		BlockPeriod:        int64(blockPeriod * 60),
		ConnectionCooldown: int64(connectionCooldown * 60),
	}

	lf := log.New(o, "", log.Ldate|log.Ltime|log.Lshortfile)
	f.logFile = NewJobWithArgument(lf.Println)
	f.logFile.PingChan = make(chan string, 100)

	j := func() {
		f.coolDownCheck()
		f.blockCheck()
		f.unblockCheck()
	}

	f.interval = NewJob(j)

	err = f.interval.StartWithTicker(1 * time.Minute)
	if err != nil {
		log.Println(err)
	}

	err = f.logFile.StartWithArgument()
	if err != nil {
		log.Println(err)
	}

	f.log("Firewall started")

	return f
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
	f.log("Blocking IPs:", IPList(ip).String())
	f.blockedIPs = append(f.blockedIPs, ip...)
}

func (f *Firewall) Unblock(ip string) {
	f.log("Unblocking IP:", ip)
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
			f.log("Removing IP from connected IPs:", c.Address)
			f.connectedIPs = append(f.connectedIPs[:i], f.connectedIPs[i+1:]...)
		} else if c.New && now-c.FirstConnection >= f.ConnectionCooldown {
			f.log("IP is no longer new:", c.Address)
			f.connectedIPs[i].New = false
		}
	}
}

func (f *Firewall) blockCheck() {
	if len(f.connectedIPs) > f.MaxConnections {
		tb := f.connectedIPs[f.MaxConnections:]
		f.Block(tb...)
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
	if f.IsBlocked(ip) {
		f.log("IP is blocked:", ip)
		return false
	}
	if f.IsConnected(ip) {
		f.log("IP is connected:", ip)
		f.updateLastConnection(ip)
		return true
	}
	f.log("IP is new:", ip)
	f.connectedIPs = append(f.connectedIPs, IP{Address: ip, FirstConnection: time.Now().Unix(), New: true})
	return true
}

func (f *Firewall) log(s ...string) {
	var sb strings.Builder
	for _, str := range s {
		sb.WriteString(str)
	}
	f.logFile.PingChan <- sb.String()
}