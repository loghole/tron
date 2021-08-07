module github.com/loghole/tron/example

go 1.16

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/loghole/tron v0.0.0
	github.com/spf13/viper v1.8.1
	google.golang.org/genproto v0.0.0-20210805201207-89edb61ffb67
	google.golang.org/grpc v1.39.1
	google.golang.org/protobuf v1.27.1
)

replace github.com/loghole/tron => ../
