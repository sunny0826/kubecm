package cloud

import (
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRegionID(t *testing.T) {
	regions, err := GetRegionID()
	require.NoError(t, err)
	assert.NotEmpty(t, regions)
	assert.Contains(t, regions, "us-east-1")
	assert.Contains(t, regions, "eu-west-1")
}

func TestGetRegionID_IsSorted(t *testing.T) {
	regions, err := GetRegionID()
	require.NoError(t, err)
	assert.True(t, sort.StringsAreSorted(regions), "regions should be sorted alphabetically")
}

func TestAWS_GetAWSConfig_InvalidAuthMode(t *testing.T) {
	a := &AWS{
		AuthMode: AWSAuth(99),
		RegionID: "us-east-1",
	}
	_, err := a.getAWSConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid AWS auth mode")
}

func TestAWS_GetAWSConfig_StaticCredentials(t *testing.T) {
	a := &AWS{
		AuthMode:        AWSAuthStaticCredentials,
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		AccessKeySecret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		RegionID:        "us-east-1",
	}
	cfg, err := a.getAWSConfig()
	require.NoError(t, err)
	assert.Equal(t, "us-east-1", cfg.Region)
}

func TestAWS_GetAWSConfig_Caching(t *testing.T) {
	a := &AWS{
		AuthMode:        AWSAuthStaticCredentials,
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		AccessKeySecret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		RegionID:        "eu-west-1",
	}
	cfg1, err := a.getAWSConfig()
	require.NoError(t, err)

	cfg2, err := a.getAWSConfig()
	require.NoError(t, err)

	// Both calls should return the same cached region
	assert.Equal(t, cfg1.Region, cfg2.Region)
	assert.NotNil(t, a.cfg, "config should be cached")
}

func TestAWS_GetAWSConfig_DefaultWithProfile(t *testing.T) {
	// Create a temporary AWS config file with a test profile
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config"
	err := os.WriteFile(configFile, []byte("[profile my-test-profile]\nregion = ap-northeast-1\n"), 0600)
	require.NoError(t, err)

	// Point the SDK at our temp config
	t.Setenv("AWS_CONFIG_FILE", configFile)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", tmpDir+"/credentials")

	a := &AWS{
		AuthMode: AWSAuthDefault,
		Profile:  "my-test-profile",
		RegionID: "ap-northeast-1",
	}
	cfg, err := a.getAWSConfig()
	require.NoError(t, err)
	assert.Equal(t, "ap-northeast-1", cfg.Region)
}
