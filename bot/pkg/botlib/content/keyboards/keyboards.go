package keyboards

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type Keyboard interface {
	AsObject() repr.Object
}
