package update

import (
	"reflect"
	"testing"
	"time"
)

func Test_getLatestReleaseInfo(t *testing.T) {
	type args struct {
		repo string
	}
	tests := []struct {
		name    string
		args    args
		want    *ReleaseInfo
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test lastest release info",
			args: args{
				repo: "sunny0826/kubectl-pod-lens",
			},
			want: &ReleaseInfo{
				Version:     "v0.2.2",
				URL:         "https://github.com/sunny0826/kubectl-pod-lens/releases/tag/v0.2.2",
				PublishedAt: time.Date(2021, 8, 29, 10, 19, 42, 0, time.UTC),
			},
		},
		{
			name: "repo not found",
			args: args{
				repo: "bar/foo",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLatestReleaseInfo(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestReleaseInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLatestReleaseInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckForUpdate(t *testing.T) {
	type args struct {
		repo           string
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    *ReleaseInfo
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "need update",
			args: args{
				repo:           "sunny0826/kubectl-pod-lens",
				currentVersion: "v0.2.1",
			},
			want: &ReleaseInfo{
				Version:     "v0.2.2",
				URL:         "https://github.com/sunny0826/kubectl-pod-lens/releases/tag/v0.2.2",
				PublishedAt: time.Date(2021, 8, 29, 10, 19, 42, 0, time.UTC),
			},
		},
		{
			name: "do not need update",
			args: args{
				repo:           "sunny0826/kubectl-pod-lens",
				currentVersion: "v0.2.2",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckForUpdate(tt.args.repo, tt.args.currentVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckForUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckForUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
