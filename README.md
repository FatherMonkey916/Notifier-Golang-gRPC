# Notifier - A server to streaming realtime data with Golang and Grpc

### Why?
- Want to build a server to stream realtime data via Websocket with Golang!

### What?
- Realtime twitter data by topic via websocket
- Technologies: 
   + Languages: Server: Golang
   + Proto: Grpc, websocket

### How?
- Backend side by Java will use Kafka to connect to Twitter API to fetch realtime data by reigistering topics
- Java side will send data via GRPC via stub to GRPC server(Golang)
- Server side(Golang) will distribute the twitter data to client via websocket proto

1. Add git sub-module
```bash
$ git submodule add https://github.com/chariot9/proto-notifier.git proto
```

1. Compile proto files
```bash
$ protoc -I proto/ proto/notifier/src/twitter/*.proto --go_out=plugins=grpc:grpc
```

```bash
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:../ \
  notifier/src/twitter/*.proto
```
