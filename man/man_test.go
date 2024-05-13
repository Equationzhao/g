package man

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGenMDAndMan(t *testing.T) {
	GenMDAndMan(afero.NewMemMapFs())
}
