package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type ObjectContent interface {
	AsObject() repr.Object
}

type ArrayContent interface {
	AsArray() repr.Array
}
