package cmd

import (
	"testing"

	"github.com/sunny0826/kubecm/pkg/cloud"
)

func Test_printListTable(t *testing.T) {
	type args struct {
		clusters []cloud.ClusterInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "exist", args: args{
			clusters: []cloud.ClusterInfo{
				{
					ID:         "id",
					Name:       "test-cluster",
					RegionID:   "cn-shanghai",
					K8sVersion: "v0.20.0",
					ConsoleURL: "https://www.test-cluster.com",
				},
			}}},
		{name: "not exist", args: args{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := printListTable(tt.args.clusters); (err != nil) != tt.wantErr {
				t.Errorf("printListTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
