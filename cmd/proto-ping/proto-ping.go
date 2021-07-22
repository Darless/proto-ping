package main

import (
	"flag"
	"time"

	"github.com/darless/protoping"
)

func main() {
	entries := []protoping.ConnEntry{}
	var loop = flag.Bool("l", false, "Loop - Send pings forever")
	var interval = flag.Int("i", 1000, "How often to ping in ms")
	var count = flag.Int("c", 0, "Stop after sending this amount of pings")
	flag.Parse()
	for _, arg := range flag.Args() {
		entry := protoping.ParseConnection(arg)
		entries = append(entries, entry)
	}
	index := 0
	for {
		for ix, _ := range entries {
			entry := &entries[ix]
			rtt, err := entry.Ping()
			entry.Stats.Error = err
			if err == nil {
				entry.Stats.Rtt.Last = rtt
				entry.Stats.Rtt.Total += rtt
				entry.Stats.Received += 1
			} else {
				entry.Stats.Rtt.Last = -1
			}
			entry.PrintResult()
		}
		index += 1
		if *count > 0 {
			if index >= *count {
				break
			}
		} else if !*loop {
			break
		}
		time.Sleep(time.Duration(*interval * int(time.Millisecond)))
	}
}
