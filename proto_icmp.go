package protoping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func (entry *ConnEntry) ICMP() error {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("could not start listening")
	}

	// Make a new ICMP message
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  entry.Stats.Sent,
			Data: []byte(""),
		},
	}

	// Get binary encoding of the ICMP message
	b, err := msg.Marshal(nil)
	if err != nil {
		conn.Close()
		return fmt.Errorf("could not marshall")
	}

	// Resolve any DNS (if used) and get the UDP Address of the target
	full := fmt.Sprintf("%s:%d", entry.Host, entry.Port)
	dst, err := net.ResolveUDPAddr("udp4", full)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to resolve DNS for %s", full)
	}

	// Send it
	start := time.Now()
	n, err := conn.WriteTo(b, dst)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to write to %s: %s", dst, err)
	} else if n != len(b) {
		conn.Close()
		return fmt.Errorf("got %v; want %v", n, len(b))
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		conn.Close()
		return fmt.Errorf("error in setting read deadline")
	}

	n, peer, err := conn.ReadFrom(reply)
	conn.Close()
	duration := time.Since(start)
	entry.Stats.Rtt.Last = duration
	if err != nil {
		return fmt.Errorf("Timeout")
	}

	icmp_proto_int := 1
	reply_msg, err := icmp.ParseMessage(icmp_proto_int, reply[:n])
	if err != nil {
		return fmt.Errorf("failed to parse reply message")
	}
	switch reply_msg.Type {
	case ipv4.ICMPTypeEchoReply:
		break
	default:
		return fmt.Errorf("got %+v from %v; want echo reply", reply_msg, peer)
	}
	return nil
}
