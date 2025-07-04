//go:build theme

package theme

import (
	"testing"

	"github.com/Equationzhao/g/internal/global"
)

func TestStyle_ToReadable(t *testing.T) {
	tests := []struct {
		name   string
		before Style
		want   Style
	}{
		{
			name: "TestStyle_ToReadable",
			before: Style{
				Color: global.BrightBlue,
			},
			want: Style{
				Color: "bright-blue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.ToReadable()
			if got.Color != tt.want.Color {
				t.Errorf("Style.ToReadable() = %v, want %v", got.Color, tt.want.Color)
			}
		})
	}
}
