package app

// Application internal config keys.
const (
	NamespaceEnv           = "NAMESPACE"
	LoggerLevelEnv         = "LOGGER_LEVEL"
	LoggerCollectorAddrEnv = "LOGGER_COLLECTOR_ADDR"
	LoggerDisableStdoutEnv = "LOGGER_DISABLE_STDOUT"
	JaegerAddrEnv          = "JAEGER_ADDR"
	JaegerSamplerType      = "JAEGER_SAMPLER_TYPE"
	JaegerSamplerParam     = "JAEGER_SAMPLER_PARAM"
	HTTPPortEnv            = "SERVICE_PORT_HTTP"
	GRPCPortEnv            = "SERVICE_PORT_GRPC"
	AdminPortEnv           = "SERVICE_PORT_ADMIN"
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
