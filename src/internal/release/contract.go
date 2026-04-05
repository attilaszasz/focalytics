package release

import (
	"fmt"
	"strings"
)

const (
	BinaryName = "focalytics"
)

type Target struct {
	GOOS               string
	GOARCH             string
	ArchiveExtension   string
	WingetArchitecture string
}

func DefaultTargets() []Target {
	return []Target{
		{GOOS: "darwin", GOARCH: "amd64", ArchiveExtension: ".tar.gz", WingetArchitecture: ""},
		{GOOS: "darwin", GOARCH: "arm64", ArchiveExtension: ".tar.gz", WingetArchitecture: ""},
		{GOOS: "linux", GOARCH: "amd64", ArchiveExtension: ".tar.gz", WingetArchitecture: ""},
		{GOOS: "linux", GOARCH: "arm64", ArchiveExtension: ".tar.gz", WingetArchitecture: ""},
		{GOOS: "windows", GOARCH: "amd64", ArchiveExtension: ".zip", WingetArchitecture: "x64"},
		{GOOS: "windows", GOARCH: "arm64", ArchiveExtension: ".zip", WingetArchitecture: "arm64"},
	}
}

func NormalizeVersion(version string) (string, error) {
	trimmed := strings.TrimSpace(version)
	trimmed = strings.TrimPrefix(trimmed, "refs/tags/")
	trimmed = strings.TrimPrefix(trimmed, "v")
	if trimmed == "" {
		return "", fmt.Errorf("version is required")
	}

	return trimmed, nil
}

func ReleaseTag(version string) (string, error) {
	normalized, err := NormalizeVersion(version)
	if err != nil {
		return "", err
	}

	return "v" + normalized, nil
}

func TargetByPlatform(goos, goarch string) (Target, error) {
	for _, target := range DefaultTargets() {
		if target.GOOS == goos && target.GOARCH == goarch {
			return target, nil
		}
	}

	return Target{}, fmt.Errorf("unsupported target %s/%s", goos, goarch)
}

func ReleaseAssetBase(version string, target Target) (string, error) {
	normalized, err := NormalizeVersion(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_v%s_%s_%s", BinaryName, normalized, target.GOOS, target.GOARCH), nil
}

func ReleaseAssetName(version string, target Target) (string, error) {
	base, err := ReleaseAssetBase(version, target)
	if err != nil {
		return "", err
	}

	return base + target.ArchiveExtension, nil
}

func BinaryFileName(target Target) string {
	if target.GOOS == "windows" {
		return BinaryName + ".exe"
	}

	return BinaryName
}

func ChecksumManifestName(version string) (string, error) {
	normalized, err := NormalizeVersion(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_v%s_checksums.txt", BinaryName, normalized), nil
}

func ReleaseAssetURL(repository, version, assetName string) (string, error) {
	repo := strings.TrimSpace(repository)
	if repo == "" {
		return "", fmt.Errorf("repository is required")
	}

	tag, err := ReleaseTag(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, tag, assetName), nil
}
