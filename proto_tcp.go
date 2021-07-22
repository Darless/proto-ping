package protoping

import (
	"fmt"
	"net"
	"time"
)

func (entry *ConnEntry) TCP() error {
	full_host := fmt.Sprintf("%s:%d", entry.Host, entry.Port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", full_host, time.Duration(time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to %s", full_host)
	}
	duration := time.Since(start)
	entry.Stats.Rtt.Last = duration
	conn.Close()
	return nil
}
