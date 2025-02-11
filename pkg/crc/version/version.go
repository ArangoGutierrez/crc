/*
Copyright (C) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"io/ioutil"
	"net/http"
	"time"
)

// The following variables are private fields and should be set when compiling with ldflags, for example --ldflags="-X github.com/code-ready/crc/pkg/version.crcVersion=vX.Y.Z
var (
	// The current version of minishift
	crcVersion = "0.0.0-unset"

	// The SHA-1 of the commit this binary is build off
	commitSha = "sha-unset"

	// Bundle version which used for the release.
	bundleVersion = "0.0.0-unset"
)

const (
	releaseInfoLink = "https://mirror.openshift.com/pub/openshift-v4/clients/crc/latest/release-info.json"
)

type CrcReleaseInfo struct {
	Version struct {
		LatestVersion string `json:"crcVersion"`
	}
}

func GetCRCVersion() string {
	return crcVersion
}

func GetCommitSha() string {
	return commitSha
}

func GetBundleVersion() string {
	return bundleVersion
}

func getCRCLatestVersionFromMirror() (*semver.Version, error) {
	var releaseInfo CrcReleaseInfo
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(releaseInfoLink)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	releaseMetaData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(releaseMetaData, &releaseInfo)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON metadata: %v", err)
	}
	version, err := semver.NewVersion(releaseInfo.Version.LatestVersion)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func NewVersionAvailable() (bool, string, error) {
	latestVersion, err := getCRCLatestVersionFromMirror()
	if err != nil {
		return false, "", err
	}
	currentVersion, err := semver.NewVersion(GetCRCVersion())
	if err != nil {
		return false, "", err
	}
	return latestVersion.GreaterThan(currentVersion), latestVersion.String(), nil
}
