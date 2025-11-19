package callbacks

import (
	"fmt"
	"strings"
)

type CallbackData interface {
	Encode() string
	With(key, value string) CallbackData
}

type PrefixedCD interface {
	Prefix() (string, error)
	Value(key string) (string, error)
	ExistsValue(key string) bool
}

// TODO refactor prefixedCD
type prefixedCD struct {
	content  string
	prefix   string
	params   map[string]string
	parsed   bool
	parseErr error
}

func NewPrefixedCD(data string) prefixedCD {
	cd := prefixedCD{content: data}
	cd.parse()
	return cd
}

func (s *prefixedCD) parse() {
	if s.parsed {
		return
	}

	parts := strings.SplitN(s.content, ":", 2)
	if len(parts) == 0 {
		s.parseErr = fmt.Errorf("callback data is empty")
		s.parsed = true
		return
	}

	s.prefix = parts[0]
	s.params = make(map[string]string)

	if len(parts) == 2 && parts[1] != "" {
		paramStrings := strings.Split(parts[1], ":")
		for _, param := range paramStrings {
			kv := strings.SplitN(param, "=", 2)
			if len(kv) == 2 {
				s.params[kv[0]] = kv[1]
			} else if len(kv) == 1 && kv[0] != "" {
				s.params[kv[0]] = ""
			}
		}
	}

	s.parsed = true
}

func (s *prefixedCD) Prefix() (string, error) {
	if !s.parsed {
		s.parse()
	}
	return s.prefix, s.parseErr
}

func (s *prefixedCD) Value(key string) (string, error) {
	if !s.parsed {
		s.parse()
	}
	if s.parseErr != nil {
		return "", s.parseErr
	}

	if val, ok := s.params[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key %q not found in callback parameters", key)
}

func (s *prefixedCD) ExistsValue(key string) bool {
	if !s.parsed {
		s.parse()
	}
	if s.parseErr != nil {
		return false
	}

	_, ok := s.params[key]
	return ok
}

func PrefixedData(content string) PrefixedCD {
	return &prefixedCD{content: content}
}
