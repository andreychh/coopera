package hsm

type Children interface {
	Initial() Spec
	All() []Spec
}

type childGroup struct {
	initial Spec
	others  []Spec
}

func (c childGroup) Initial() Spec {
	return c.initial
}

func (c childGroup) All() []Spec {
	res := make([]Spec, 0, 1+len(c.others))
	res = append(res, c.initial)
	res = append(res, c.others...)
	return res
}

func Group(initial Spec, others ...Spec) Children {
	return childGroup{
		initial: initial,
		others:  others,
	}
}
