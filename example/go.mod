module example

go 1.16

require (
	github.com/gadavy/ozzo-validation/v4 v4.3.2
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0
	github.com/lissteron/simplerr v0.8.0
	github.com/loghole/tron v0.15.2
	github.com/spf13/viper v1.7.1
	google.golang.org/genproto v0.0.0-20210122163508-8081c04a3579
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/loghole/tron => ../
