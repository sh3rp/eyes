package probe

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func GetLatency(srcIP, dstIP string, dstPort uint16) (time.Duration, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var receiveTime time.Time
	var err error

	addrs, err := net.LookupHost(dstIP)
	if err != nil {
		log.Error().Msgf("Error resolving %s. %s\n", dstIP, err)
	}
	for _, addr := range addrs {
		if addr != "127.0.0.1" && !strings.Contains(addr, ":") {
			ipAddr := addr
			elements := strings.Split(ipAddr, "/")
			dstIP = elements[0]
			break
		}
	}

	go func() {
		receiveTime, err = WaitForResponse(srcIP, dstIP, dstPort)
		wg.Done()
	}()

	time.Sleep(1 * time.Millisecond)
	sendTime := SendPing(srcIP, dstIP, 0, dstPort)

	wg.Wait()
	return receiveTime.Sub(sendTime), nil
}

func SendPing(srcIP, dstIP string, srcPort, dstPort uint16) time.Time {
	if srcPort == 0 {
		srcPort = getNextLocalPort()
	}

	packet := TCPHeader{
		Src:        srcPort,
		Dst:        dstPort,
		Seq:        rand.Uint32(),
		Ack:        0,
		DataOffset: 5,
		Reserved:   0,
		ECN:        0,
		Ctrl:       2,
		Window:     0xaaaa,
		Checksum:   0,
		Urgent:     0,
		Options:    []TCPOption{},
	}

	data := packet.MarshalTCP()

	if !validIP(srcIP) {
		log.Info().Msgf("Invalid src IP: %v", srcIP)
		return time.Now()
	}

	if !validIP(dstIP) {
		log.Info().Msgf("Invalid dst IP: %v", srcIP)
		return time.Now()
	}

	packet.Checksum = Checksum(data, to4byte(srcIP), to4byte(dstIP))

	data = packet.MarshalTCP()

	conn, err := net.Dial("ip4:tcp", dstIP)
	if err != nil {
		log.Info().Msgf("Dial: %s\n", err)
		return time.Now()
	}

	sendTime := time.Now()

	numWrote, err := conn.Write(data)

	if err != nil {
		log.Error().Msgf("Write: %s\n", err)
		return time.Now()
	}

	if numWrote != len(data) {
		log.Error().Msgf("Error writing %d/%d bytes\n", numWrote, len(data))
		return time.Now()
	}

	conn.Close()

	return sendTime
}

func WaitForResponse(localAddress, remoteAddress string, port uint16) (time.Time, error) {
	netaddr, err := net.ResolveIPAddr("ip4", localAddress)
	if err != nil {
		log.Error().Msgf("ERROR: net.ResolveIPAddr: %s. %s\n", localAddress, netaddr)
		return time.Now(), err
	}

	conn, err := net.ListenIP("ip4:tcp", netaddr)
	if err != nil {
		log.Error().Msgf("ListenIP: %s\n", err)
		return time.Now(), err
	}

	conn.SetReadDeadline(time.Now().Add(time.Duration(3 * time.Second)))

	var receiveTime time.Time
	for {
		buf := make([]byte, 1024)
		numRead, raddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Error().Msgf("ReadFrom: %s\n", err)
			return time.Now(), err
		}
		if raddr.String() != remoteAddress {
			continue
		}
		receiveTime = time.Now()
		tcp := ParseTCP(buf[:numRead])
		if tcp.HasFlag(RST) || (tcp.HasFlag(SYN) && tcp.HasFlag(ACK)) {
			break
		}
	}
	return receiveTime, nil
}

// Grab first interface found and the first IP on it
func GetInterface() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error().Msgf("Error, no interfaces: %s", err)
		return ""
	}
	for _, iface := range interfaces {
		if strings.HasPrefix(iface.Name, "lo") {
			continue
		}
		addrs, err := iface.Addrs()

		if err != nil {
			log.Error().Msgf(" %s. %s", iface.Name, err)
			continue
		}
		var retAddr net.Addr
		for _, a := range addrs {
			if !strings.Contains(a.String(), ":") {
				retAddr = a
				break
			}
		}
		if retAddr != nil {
			return retAddr.String()[:strings.Index(retAddr.String(), "/")]
		}
	}

	return ""
}

func to4byte(addr string) [4]byte {
	parts := strings.Split(addr, ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Error().Msgf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])
	return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}

func getNextLocalPort() uint16 {
	return 0
}

func validIP(ip string) bool {
	if !strings.Contains(ip, ".") {
		return false
	}

	tokens := strings.Split(ip, ".")

	if len(tokens) != 4 {
		return false
	}

	for _, toke := range tokens {
		_, err := strconv.Atoi(toke)
		if err != nil {
			return false
		}
	}

	return true
}
