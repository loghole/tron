package models

import (
	"strings"
)

type Options map[string]struct{}

func ParseTronOptions(s string) (Options, bool) {
	m := TronOptions.FindStringSubmatch(s)
	if len(m) == 0 {
		return nil, false
	}

	opts := make(Options)

	for _, val := range strings.Split(m[1], TronOptionsSep) {
		opts[strings.ToLower(val)] = struct{}{}
	}

	return opts, len(opts) > 0
}

func (opts Options) Apply(data string) string {
	if _, ok := opts[TronOptionJSON]; ok {
		data = strings.Replace(data, "[]byte", "json.RawMessage", 1)
	}

	return data
}
