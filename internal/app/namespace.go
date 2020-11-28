package app

import (
	"path"
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

// ValuesPath returns config values path for current namespace.
func (n Namespace) ValuesPath() string {
	var name string

	switch n {
	case NamespaceDev:
		name = ValuesDevName
	case NamespaceStage:
		name = ValuesStgName
	case NamespaceProd:
		name = ValuesProdName
	case NamespaceLocal:
		name = ValuesLocalName
	default:
		panic("unknown namespace")
	}

	return path.Join(DeploymentsDir, ValuesDir, strings.Join([]string{name, ValuesExt}, "."))
}
