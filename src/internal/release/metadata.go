package release

import "fmt"

type AssetReference struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	ArchiveName  string `json:"archive_name"`
	BinaryName   string `json:"binary_name"`
	URL          string `json:"url"`
	SHA256       string `json:"sha256"`
}

type HomebrewInput struct {
	FormulaName         string           `json:"formula_name"`
	Description         string           `json:"description"`
	Homepage            string           `json:"homepage"`
	License             string           `json:"license"`
	Version             string           `json:"version"`
	ReleaseTag          string           `json:"release_tag"`
	ChecksumManifestURL string           `json:"checksum_manifest_url"`
	Artifacts           []AssetReference `json:"artifacts"`
}

type WinGetInstaller struct {
	Architecture         string   `json:"architecture"`
	InstallerType        string   `json:"installer_type"`
	InstallerURL         string   `json:"installer_url"`
	InstallerSHA256      string   `json:"installer_sha256"`
	NestedInstallerType  string   `json:"nested_installer_type"`
	NestedInstallerFiles []string `json:"nested_installer_files"`
}

type WinGetInput struct {
	PackageIdentifier   string            `json:"package_identifier"`
	PackageName         string            `json:"package_name"`
	Publisher           string            `json:"publisher"`
	PackageVersion      string            `json:"package_version"`
	ReleaseTag          string            `json:"release_tag"`
	ChecksumManifestURL string            `json:"checksum_manifest_url"`
	Installers          []WinGetInstaller `json:"installers"`
}

func BuildAssetReferences(repository, version string, targets []Target, checksums map[string]string) ([]AssetReference, error) {
	if err := VerifyChecksumContract(version, targets, checksums); err != nil {
		return nil, err
	}

	references := make([]AssetReference, 0, len(targets))
	for _, target := range targets {
		assetName, err := ReleaseAssetName(version, target)
		if err != nil {
			return nil, err
		}

		assetURL, err := ReleaseAssetURL(repository, version, assetName)
		if err != nil {
			return nil, err
		}

		references = append(references, AssetReference{
			OS:           target.GOOS,
			Architecture: target.GOARCH,
			ArchiveName:  assetName,
			BinaryName:   BinaryFileName(target),
			URL:          assetURL,
			SHA256:       checksums[assetName],
		})
	}

	return references, nil
}

func BuildHomebrewInput(repository, version string, checksums map[string]string) (HomebrewInput, error) {
	versionValue, err := NormalizeVersion(version)
	if err != nil {
		return HomebrewInput{}, err
	}

	tag, err := ReleaseTag(version)
	if err != nil {
		return HomebrewInput{}, err
	}

	checksumManifestName, err := ChecksumManifestName(version)
	if err != nil {
		return HomebrewInput{}, err
	}

	checksumManifestURL, err := ReleaseAssetURL(repository, version, checksumManifestName)
	if err != nil {
		return HomebrewInput{}, err
	}

	homebrewTargets := filterTargets(func(target Target) bool {
		return target.GOOS == "darwin" || target.GOOS == "linux"
	})
	artifacts, err := BuildAssetReferences(repository, version, homebrewTargets, subsetChecksums(version, homebrewTargets, checksums))
	if err != nil {
		return HomebrewInput{}, err
	}

	return HomebrewInput{
		FormulaName:         BinaryName,
		Description:         "Analyze a local photo archive",
		Homepage:            fmt.Sprintf("https://github.com/%s", repository),
		License:             "MIT",
		Version:             versionValue,
		ReleaseTag:          tag,
		ChecksumManifestURL: checksumManifestURL,
		Artifacts:           artifacts,
	}, nil
}

func BuildWinGetInput(repository, version string, checksums map[string]string) (WinGetInput, error) {
	versionValue, err := NormalizeVersion(version)
	if err != nil {
		return WinGetInput{}, err
	}

	tag, err := ReleaseTag(version)
	if err != nil {
		return WinGetInput{}, err
	}

	checksumManifestName, err := ChecksumManifestName(version)
	if err != nil {
		return WinGetInput{}, err
	}

	checksumManifestURL, err := ReleaseAssetURL(repository, version, checksumManifestName)
	if err != nil {
		return WinGetInput{}, err
	}

	windowsTargets := filterTargets(func(target Target) bool {
		return target.GOOS == "windows"
	})
	if err := VerifyChecksumContract(version, windowsTargets, subsetChecksums(version, windowsTargets, checksums)); err != nil {
		return WinGetInput{}, err
	}

	installers := make([]WinGetInstaller, 0, len(windowsTargets))
	for _, target := range windowsTargets {
		assetName, err := ReleaseAssetName(version, target)
		if err != nil {
			return WinGetInput{}, err
		}

		assetURL, err := ReleaseAssetURL(repository, version, assetName)
		if err != nil {
			return WinGetInput{}, err
		}

		installers = append(installers, WinGetInstaller{
			Architecture:         target.WingetArchitecture,
			InstallerType:        "zip",
			InstallerURL:         assetURL,
			InstallerSHA256:      checksums[assetName],
			NestedInstallerType:  "portable",
			NestedInstallerFiles: []string{BinaryFileName(target)},
		})
	}

	return WinGetInput{
		PackageIdentifier:   "Attila.Focalytics",
		PackageName:         "focalytics",
		Publisher:           "Attila",
		PackageVersion:      versionValue,
		ReleaseTag:          tag,
		ChecksumManifestURL: checksumManifestURL,
		Installers:          installers,
	}, nil
}

func filterTargets(include func(Target) bool) []Target {
	targets := DefaultTargets()
	filtered := make([]Target, 0, len(targets))
	for _, target := range targets {
		if include(target) {
			filtered = append(filtered, target)
		}
	}

	return filtered
}

func subsetChecksums(version string, targets []Target, checksums map[string]string) map[string]string {
	subset := make(map[string]string, len(targets))
	for _, target := range targets {
		assetName, err := ReleaseAssetName(version, target)
		if err != nil {
			continue
		}
		if checksum, ok := checksums[assetName]; ok {
			subset[assetName] = checksum
		}
	}

	return subset
}
