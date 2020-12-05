module example

go 1.15

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/loghole/tron v0.0.0-00010101000000-000000000000
	github.com/pelletier/go-toml v1.5.0 // indirect
	github.com/spf13/viper v1.7.1
	google.golang.org/genproto v0.0.0-20201204160425-06b3db808446
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/loghole/tron => ../
