package update

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type releaseTransport func(*http.Request) (*http.Response, error)

func (f releaseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func stubReleaseAPI(t *testing.T) {
	t.Helper()
	originalClient := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = originalClient })
	http.DefaultClient = &http.Client{Transport: releaseTransport(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet || req.URL.Scheme != "https" || req.URL.Host != "api.github.com" {
			return nil, fmt.Errorf("unexpected release request: %s %s", req.Method, req.URL)
		}
		status := http.StatusOK
		body := `{"tag_name":"v0.2.2","html_url":"https://github.com/sunny0826/kubectl-pod-lens/releases/tag/v0.2.2","published_at":"2021-08-29T10:19:42Z"}`
		switch req.URL.Path {
		case "/repos/sunny0826/kubectl-pod-lens/releases/latest":
		case "/repos/bar/foo/releases/latest":
			status = http.StatusNotFound
			body = `{"message":"Not Found"}`
		default:
			return nil, fmt.Errorf("unexpected release path: %s", req.URL.Path)
		}
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
}

func Test_getLatestReleaseInfo(t *testing.T) {
	stubReleaseAPI(t)
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
	stubReleaseAPI(t)
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
				currentVersion: "0.2.1",
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
				currentVersion: "0.2.2",
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
