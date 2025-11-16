package repr

import (
	"fmt"
	"strconv"
	"strings"
)

type path struct {
	raw string
}

func (p path) Empty() bool {
	return p.raw == ""
}

func (p path) Index() (int, error) {
	if p.raw == "" {
		return 0, fmt.Errorf("path is empty")
	}
	if p.raw[0] == '[' {
		end := strings.IndexByte(p.raw, ']')
		if end == -1 {
			return 0, fmt.Errorf("unclosed bracket in path")
		}
		indexStr := p.raw[1:end]
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return 0, fmt.Errorf("invalid index: %s", indexStr)
		}
		return index, nil
	}
	return 0, fmt.Errorf("expected index, got key")
}

func (p path) Key() (string, error) {
	if p.raw == "" {
		return "", fmt.Errorf("path is empty")
	}
	if p.raw[0] == '[' {
		return "", fmt.Errorf("expected key, got index")
	}
	end := strings.IndexAny(p.raw, ".[")
	if end == -1 {
		return p.raw, nil
	}
	return p.raw[:end], nil
}

func (p path) Tail() Path {
	if p.raw == "" {
		return path{raw: ""}
	}
	if p.raw[0] == '[' {
		end := strings.IndexByte(p.raw, ']')
		if end == -1 {
			return path{raw: ""}
		}
		rest := p.raw[end+1:]
		if len(rest) > 0 && rest[0] == '.' {
			rest = rest[1:]
		}
		return path{raw: rest}
	}
	end := strings.IndexAny(p.raw, ".[")
	if end == -1 {
		return path{raw: ""}
	}
	rest := p.raw[end:]
	if len(rest) > 0 && rest[0] == '.' {
		rest = rest[1:]
	}
	return path{raw: rest}
}

func NewPath(s string) Path {
	return path{raw: s}
}
