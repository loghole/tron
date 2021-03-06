package app

import (
	"strings"
)

// Namespace type.
type Namespace string

// ParseNamespace return app Namespace.
func ParseNamespace(s string) Namespace {
	if localConfigEnabled {
		return NamespaceLocal
	}

	switch strings.ToLower(s) {
	case "d", "dev", "develop":
		return NamespaceDev
	case "s", "stg", "stage", "demo":
		return NamespaceStage
	case "p", "prod", "production":
		return NamespaceProd
	default:
		panic("unknown namespace")
	}
}

func (n Namespace) String() string {
	return string(n)
}

// ValuesName returns config values name for current namespace.
func (n Namespace) ValuesName() string {
	switch n {
	case NamespaceDev:
		return ValuesDevName
	case NamespaceStage:
		return ValuesStgName
	case NamespaceProd:
		return ValuesProdName
	case NamespaceLocal:
		return ValuesLocalName
	default:
		panic("unknown namespace")
	}
}
