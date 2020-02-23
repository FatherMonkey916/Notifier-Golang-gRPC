# proto-notifier

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