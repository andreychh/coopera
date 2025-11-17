package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type Content interface {
	Structure() repr.Structure
	Method() string
}
