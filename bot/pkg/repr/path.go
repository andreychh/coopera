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

func (p path) Head() (Segment, error) {
	if p.raw == "" {
		return nil, fmt.Errorf("path is empty")
	}
	if p.raw[0] == '[' {
		end := strings.IndexByte(p.raw, ']')
		if end == -1 {
			return nil, fmt.Errorf("unclosed bracket in path")
		}
		indexStr := p.raw[1:end]
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid index %q: %w", indexStr, err)
		}
		return indexSegment{value: index}, nil
	}
	end := strings.IndexAny(p.raw, ".[")
	if end == -1 {
		return keySegment{value: p.raw}, nil
	}
	return keySegment{value: p.raw[:end]}, nil
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

func PathOf(s string) Path {
	return path{raw: s}
}

func EmptyPath() Path {
	return path{raw: ""}
}
