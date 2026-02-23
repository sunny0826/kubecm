package cmd

import (
	"os"
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

func Test_checkEnvForSecret(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "err_input",
			args: args{
				num: -1,
			},
			want:  "",
			want1: "",
		},
		{
			name: "ali_env",
			args: args{
				num: 0,
			},
			want:  "aliyun_env_id",
			want1: "aliyun_env_sec",
		},
		{
			name: "ten_env",
			args: args{
				num: 1,
			},
			want:  "ten_env_id",
			want1: "ten_env_sec",
		},
		{
			name: "aws_env",
			args: args{
				num: 3,
			},
			want:  "aws_env_id",
			want1: "aws_env_sec",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "ali_env":
				os.Setenv("ACCESS_KEY_ID", "aliyun_env_id")
				os.Setenv("ACCESS_KEY_SECRET", "aliyun_env_sec")
			case "ten_env":
				os.Setenv("TENCENTCLOUD_SECRET_ID", "ten_env_id")
				os.Setenv("TENCENTCLOUD_SECRET_KEY", "ten_env_sec")
			case "aws_env":
				os.Setenv("AWS_ACCESS_KEY_ID", "aws_env_id")
				os.Setenv("AWS_SECRET_ACCESS_KEY", "aws_env_sec")
			}
			got, got1 := checkEnvForSecret(tt.args.num)
			if got != tt.want {
				t.Errorf("checkEnvForSecret() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkEnvForSecret() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
