package language

import (
	"net/http"
	"strings"
)

const (
	English = "en"
	Spanish = "es"
)

func ResolveHeader(value string) string {
	for _, candidate := range strings.Split(value, ",") {
		current := strings.TrimSpace(candidate)
		if current == "" {
			continue
		}
		if index := strings.IndexByte(current, ';'); index >= 0 {
			current = current[:index]
		}
		current = strings.TrimSpace(strings.ToLower(current))
		if current == "" {
			continue
		}
		if index := strings.IndexAny(current, "-_"); index >= 0 {
			current = current[:index]
		}
		switch current {
		case Spanish:
			return Spanish
		case English:
			return English
		}
	}
	return English
}

func ResolveRequest(r *http.Request) string {
	if r == nil {
		return English
	}
	return ResolveHeader(r.Header.Get("Accept-Language"))
}
