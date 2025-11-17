package repr

type indexSegment struct {
	value int
}

func (i indexSegment) Index() (int, bool) {
	return i.value, true
}

func (i indexSegment) Key() (string, bool) {
	return "", false
}
