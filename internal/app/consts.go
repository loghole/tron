package app

// Application internal config keys.
const (
	NamespaceEnv           = "namespace"
	LoggerLevelEnv         = "logger_level"
	LoggerCollectorAddrEnv = "logger_collector_addr"
	LoggerDisableStdoutEnv = "logger_disable_stdout"
	JaegerAddrEnv          = "jaeger_addr"
	HTTPPortEnv            = "service_port_http"
	GRPCPortEnv            = "service_port_grpc"
	AdminPortEnv           = "service_port_admin"
)

// Default names for application files and directories.
const (
	DeploymentsDir  = ".deploy"
	ValuesDir       = "config"
	ValuesExt       = "yaml"
	ValuesBaseName  = "values"
	ValuesDevName   = "values_dev"
	ValuesLocalName = "values_local"
	ValuesStgName   = "values_stage"
	ValuesProdName  = "values_prod"
)

// Namespace values.
const (
	NamespaceDev   Namespace = "dev"
	NamespaceLocal Namespace = "local"
	NamespaceStage Namespace = "stage"
	NamespaceProd  Namespace = "prod"
)
