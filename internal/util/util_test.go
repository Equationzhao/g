package util

import "testing"

func TestSplitNumberAndUnit(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 string
	}{
		{
			name: "123bit",
			args: args{
				s: "123bit",
			},
			want:  123,
			want1: "bit",
		},
		{
			name: "123",
			args: args{
				s: "123",
			},
			want:  123,
			want1: "",
		},
		{
			name: "1,234.321bit",
			args: args{
				s: "1,234.321bit",
			},
			want:  1234.321,
			want1: "bit",
		},
		{
			name: "-1,234.321bit",
			args: args{
				s: "-1,234.321bit",
			},
			want:  -1234.321,
			want1: "bit",
		},
		{
			name: "bit",
			args: args{
				s: "bit",
			},
			want:  0,
			want1: "bit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitNumberAndUnit(tt.args.s)
			if got != tt.want {
				t.Errorf("SplitNumberAndUnit() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitNumberAndUnit() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
