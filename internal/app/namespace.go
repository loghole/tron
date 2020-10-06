package app

import (
	"path"
	"strings"
)

type Namespace string

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
		return NamespaceDev
	}
}

func (n Namespace) String() string {
	return string(n)
}

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
