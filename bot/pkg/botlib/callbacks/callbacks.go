package callbacks

type Incoming interface {
	Prefix() string
	Value(key string) (string, bool)
}

type Outgoing interface {
	With(key, value string) Outgoing
	String() string
}
