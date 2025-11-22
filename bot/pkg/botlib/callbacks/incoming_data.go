package callbacks

import "strings"

type incomingData struct {
	data string
}

func (i incomingData) Prefix() string {
	index := strings.Index(i.data, ":")
	if index == -1 {
		return i.data
	}
	return i.data[:index]
}

func (i incomingData) Value(key string) (string, bool) {
	parts := strings.Split(i.data, ":")
	for i := 1; i < len(parts); i++ {
		segment := strings.SplitN(parts[i], "=", 2)
		if len(segment) != 2 {
			continue
		}
		if segment[0] == key {
			return segment[1], true
		}
	}
	return "", false
}

func IncomingData(data string) Incoming {
	return incomingData{data: data}
}
