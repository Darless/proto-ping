package protoping

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ValidProtocol(protocol string) bool {
	switch protocol {
	case "icmp":
		return true
	case "tcp":
		return true
	default:
		return false
	}
}

/**
 * Parse a argument connection
 * Format: key=val[,key=val,...]
 * proto=Protocol, defaults to ICMP
 * host=Host to ping (IPv4, IPv6, FQDN)
 * port=Destination port or applicable
 *
 * This supports additional keyword arguments that other
 * protocols may use, they will be accessible via the attributes
 * variable.
 */
func ParseConnection(content string) ConnEntry {
	entry := ConnEntry{
		Content:  content,
		Protocol: "icmp",
		Host:     "",
		Port:     53,
	}
	arr := strings.Split(content, ",")
	for _, kv_data := range arr {
		kv_arr := strings.SplitN(kv_data, "=", 2)
		if len(kv_arr) != 2 {
			fmt.Printf("Not key=val attribute: %s", kv_data)
			panic("key/val parsing")
		}
		key := kv_arr[0]
		val := kv_arr[1]
		attrib := ConnAttrib{
			Key: key,
			Val: val,
		}
		entry.Attributes = append(entry.Attributes, attrib)

		if key == "proto" {
			entry.Protocol = val
		} else if key == "host" {
			entry.Host = val
		} else if key == "port" {
			entry.Port, _ = strconv.Atoi(val)
		}
	}
	if !ValidProtocol(entry.Protocol) {
		panic(fmt.Sprintf("Not a valid protocol: %s", entry.Protocol))
	}
	return entry
}

type ConnAttrib struct {
	Key string
	Val string
}

type RTT struct {
	Last  time.Duration
	Min   time.Duration
	Max   time.Duration
	Total time.Duration
}

type Statistics struct {
	Sent     int
	Received int
	Rtt      RTT
	Error    error
}
type ConnEntry struct {
	Content    string
	Protocol   string
	Host       string
	Port       int
	Attributes []ConnAttrib
	Stats      Statistics
}

func (entry *ConnEntry) Ping() (time.Duration, error) {
	entry.Stats.Sent++
	var rtt time.Duration
	var err error

	if entry.Protocol == "icmp" {
		err = entry.ICMP()
	} else if entry.Protocol == "tcp" {
		err = entry.TCP()
	} else {
		err = fmt.Errorf("failed")
	}
	if err == nil {
		rtt = entry.Stats.Rtt.Last
	}
	return rtt, err
}

func (entry *ConnEntry) PrintResult() {
	loss := (entry.Stats.Sent - entry.Stats.Received) / entry.Stats.Sent * 100
	if entry.Stats.Error != nil {
		fmt.Printf("%s : [%d], Error %s (%d%% loss)\n",
			entry.Content,
			entry.Stats.Sent,
			entry.Stats.Error,
			loss)
	} else {
		fmt.Printf("%s : [%d], %d ms (%d%% loss)\n",
			entry.Content,
			entry.Stats.Sent,
			entry.Stats.Rtt.Last.Milliseconds(),
			loss)
	}
}
