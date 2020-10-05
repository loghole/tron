package app

const (
	NamespaceEnv = "namespace"

	LoggerLevelEnv         = "logger_level"
	LoggerCollectorAddrEnv = "logger_collector_addr"
	LoggerDisableStdoutEnv = "logger_disable_stdout"

	JaegerAddrEnv = "jaeger_addr"

	HTTPPortEnv  = "service_port_http"
	GRPCPortEnv  = "service_port_grpc"
	AdminPortEnv = "service_port_admin"
)

const (
	ValuesPath     = ".deploy/config"
	ValuesExt      = "yaml"
	ValuesBaseName = "values"
	ValuesDevName  = "values_dev"
	ValuesStgName  = "values_stg"
	ValuesProdName = "values_prod"
)

const (
	NamespaceDev   Namespace = "dev"
	NamespaceLocal Namespace = "local"
	NamespaceStage Namespace = "stage"
	NamespaceProd  Namespace = "prod"
)
