module example

go 1.15

require (
	github.com/gadavy/ozzo-validation/v4 v4.3.2
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/lissteron/simplerr v0.7.0
	github.com/loghole/tron v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.7.1
	google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/loghole/tron => ../
