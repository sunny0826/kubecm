package cloud

import (
	"testing"
)

func TestGetRegionID(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRegionID()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GetRegionID() is empty.")
			}
		})
	}
}
