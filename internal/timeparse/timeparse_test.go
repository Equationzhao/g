package timeparse

import "testing"

func TestTransform(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name         string
		args         args
		wantGoFormat string
	}{
		{
			name: "test",
			args: args{
				format: "%F %T",
			},
			wantGoFormat: "2006-01-02 15:04:05",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotGoFormat := Transform(tt.args.format); gotGoFormat != tt.wantGoFormat {
				t.Errorf("Transform() = %v, want %v", gotGoFormat, tt.wantGoFormat)
			}
		})
	}
}
