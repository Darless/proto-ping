# protoping

NOTE: This is a brand new project and will be fluid for commands and arguments.

The purpose of this repository is to have a Golang tool that can ping using any supported protocol.

Main featurues:
- Ability to ping multiple hosts
- Ability to define different setting per host

Usage:

Perform an ICMP ping to host 4.2.2.3:
```
proto-ping.go host=4.2.2.3
host=4.2.2.3 : [1], 21 ms (0% loss)
```

Perform an ICMP ping to multiple hosts:
```
proto-ping.go host=1.2.3.4 host=4.2.2.3
host=1.2.3.4 : [1], Error Timeout (100% loss)
host=4.2.2.3 : [1], 21 ms (0% loss)
```

Perform a TCP ping and ICMP ping in 1 call
```
proto-ping.go host=4.2.2.2 proto=tcp,host=4.2.2.2,port=80 proto=tcp,host=google.com,port=80

host=4.2.2.2 : [1], 22 ms (0% loss)
proto=tcp,host=4.2.2.2,port=80 : [1], Error failed to connect to 4.2.2.2:80 (100% loss)
proto=tcp,host=google.com,port=80 : [1], 26 ms (0% loss)
```

