package repr

type Encodable interface {
	Encode() ([]byte, error)
}

type Array interface {
	Encodable
	WithElement(element Encodable) Array
	Extend(array Array) Array
	AsSlice() []Encodable
}

type Object interface {
	Encodable
	WithField(key string, value Encodable) Object
}
