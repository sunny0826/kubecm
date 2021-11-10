package cmd

import (
	"testing"
)

func Test_checkFlags(t *testing.T) {
	type args struct {
		provider string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "exist",
			args: args{
				provider: "alibabacloud",
			},
			want: 0,
		},
		{
			name: "notExist",
			args: args{
				provider: "alibabaclou",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFlags(tt.args.provider); got != tt.want {
				t.Errorf("checkFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}
