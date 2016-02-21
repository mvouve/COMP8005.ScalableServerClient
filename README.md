# COMP8005.ScaleableServerClient
This project was written in Go and is intended to aid in comparing the preformance of epoll, select and a traditional multi-threaded network archetecture using go routines. This repository contains the client the servers can be found in
* [multithreaded](https://github.com/mvouve/COMP8005.ScalableServer)
* [select](https://github.com/mvouve/COMP8005.SelectScalableServer)
* [epoll](https://github.com/mvouve/COMP8005.EPollScalableServer)

This repositiory contains the documentation for all 3 servers as well as the client.

##Usage
This server can be envoked using the syntax of:
```bash
./COMP8005.ScalableServerClient -d [ammount of data] -r [connection per client] -i [itterations per connection] -c [number of clients] [host:port of server]
```

When terminated, the process will exit and generate an XLSX report listing useful information including ammount transfered to server per connection and round trip time per connection.

##Known Issues
For larger loads this program may crash when used with the EPoll server.
