package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ReleaseInfo stores information about a release
// nolint:tagliatelle
type ReleaseInfo struct {
	Version     string    `json:"tag_name"`
	URL         string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
}

// CheckForUpdate checks whether this software has had a newer release on GitHub
func CheckForUpdate(repo, currentVersion string) (*ReleaseInfo, error) {
	releaseInfo, err := getLatestReleaseInfo(repo)
	if err != nil {
		return nil, err
	}

	if releaseInfo.Version != currentVersion {
		return releaseInfo, nil
	}
	return nil, nil
}

func getLatestReleaseInfo(repo string) (*ReleaseInfo, error) {
	var latestRelease ReleaseInfo
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response status code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &latestRelease)
	if err != nil {
		return nil, err
	}
	return &latestRelease, nil
}
