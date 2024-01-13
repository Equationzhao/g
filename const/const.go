package constval

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Hashable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string | uintptr | ~unsafe.Pointer
}
