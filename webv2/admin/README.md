## ts-proto
```shell
windows:
protoc --plugin=protoc-gen-ts_proto=.\\node_modules\\.bin\\protoc-gen-ts_proto.cmd --ts_proto_opt=esModuleInterop=true --ts_proto_out=. ./xxx.proto 

mac: 
protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto --ts_proto_opt=esModuleInterop=true --ts_proto_out=. ./xxx.proto

```