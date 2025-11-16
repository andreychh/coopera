package repr

type keySegment struct {
	value string
}

func (k keySegment) Index() (int, bool) {
	return 0, false
}

func (k keySegment) Key() (string, bool) {
	return k.value, true
}
