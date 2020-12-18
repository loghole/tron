package http

import (
	"net/textproto"
)

func headerMatcher(key string) (string, bool) {
	return textproto.CanonicalMIMEHeaderKey(key), true
}
