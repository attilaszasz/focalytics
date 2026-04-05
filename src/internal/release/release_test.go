package release

import (
	"strings"
	"testing"
)

func TestNormalizeVersionAndReleaseTag(t *testing.T) {
	version, err := NormalizeVersion("refs/tags/v1.2.3")
	if err != nil {
		t.Fatalf("expected normalized version: %v", err)
	}
	if version != "1.2.3" {
		t.Fatalf("unexpected normalized version: %s", version)
	}

	tag, err := ReleaseTag("1.2.3")
	if err != nil {
		t.Fatalf("expected release tag: %v", err)
	}
	if tag != "v1.2.3" {
		t.Fatalf("unexpected release tag: %s", tag)
	}

	if _, err := NormalizeVersion("   "); err == nil {
		t.Fatal("expected empty version error")
	}
}

func TestReleaseAssetNameUsesCanonicalContract(t *testing.T) {
	target, err := TargetByPlatform("windows", "amd64")
	if err != nil {
		t.Fatalf("expected supported target: %v", err)
	}

	assetName, err := ReleaseAssetName("v1.2.3", target)
	if err != nil {
		t.Fatalf("expected asset name: %v", err)
	}

	if assetName != "focalytics_v1.2.3_windows_amd64.zip" {
		t.Fatalf("unexpected asset name: %s", assetName)
	}

	if got := BinaryFileName(target); got != "focalytics.exe" {
		t.Fatalf("unexpected binary name: %s", got)
	}

	checksumName, err := ChecksumManifestName("v1.2.3")
	if err != nil {
		t.Fatalf("expected checksum manifest name: %v", err)
	}
	if checksumName != "focalytics_v1.2.3_checksums.txt" {
		t.Fatalf("unexpected checksum manifest name: %s", checksumName)
	}
}

func TestVerifyChecksumContractDetectsDrift(t *testing.T) {
	targets := DefaultTargets()
	checksums := map[string]string{}
	for index, target := range targets[:len(targets)-1] {
		assetName, err := ReleaseAssetName("1.2.3", target)
		if err != nil {
			t.Fatalf("expected asset name: %v", err)
		}
		checksums[assetName] = strings.Repeat(string(rune('a'+index)), 64)
	}
	checksums["unexpected.txt"] = strings.Repeat("f", 64)

	err := VerifyChecksumContract("1.2.3", targets, checksums)
	if err == nil {
		t.Fatal("expected contract drift error")
	}
	if !strings.Contains(err.Error(), "missing") || !strings.Contains(err.Error(), "unexpected") {
		t.Fatalf("expected missing and unexpected details, got %v", err)
	}
}

func TestTargetByPlatformRejectsUnknownTarget(t *testing.T) {
	if _, err := TargetByPlatform("plan9", "amd64"); err == nil {
		t.Fatal("expected unsupported target error")
	}
}

func TestBuildMetadataUsesReleaseAssetURLs(t *testing.T) {
	checksums := map[string]string{}
	for _, target := range DefaultTargets() {
		assetName, err := ReleaseAssetName("v2.0.0", target)
		if err != nil {
			t.Fatalf("expected asset name: %v", err)
		}
		checksums[assetName] = strings.Repeat("b", 64)
	}

	homebrewInput, err := BuildHomebrewInput("attila/focalytics", "v2.0.0", checksums)
	if err != nil {
		t.Fatalf("expected homebrew input: %v", err)
	}
	if len(homebrewInput.Artifacts) != 4 {
		t.Fatalf("unexpected homebrew artifact count: %d", len(homebrewInput.Artifacts))
	}
	for _, artifact := range homebrewInput.Artifacts {
		if !strings.HasPrefix(artifact.URL, "https://github.com/attila/focalytics/releases/download/v2.0.0/") {
			t.Fatalf("expected release asset URL, got %s", artifact.URL)
		}
	}

	wingetInput, err := BuildWinGetInput("attila/focalytics", "v2.0.0", checksums)
	if err != nil {
		t.Fatalf("expected winget input: %v", err)
	}
	if len(wingetInput.Installers) != 2 {
		t.Fatalf("unexpected winget installer count: %d", len(wingetInput.Installers))
	}
	for _, installer := range wingetInput.Installers {
		if installer.InstallerType != "zip" {
			t.Fatalf("unexpected installer type: %s", installer.InstallerType)
		}
	}
}

func TestBuildMetadataRejectsMissingChecksums(t *testing.T) {
	checksums := map[string]string{}
	for _, target := range DefaultTargets()[:2] {
		assetName, err := ReleaseAssetName("v2.0.0", target)
		if err != nil {
			t.Fatalf("expected asset name: %v", err)
		}
		checksums[assetName] = strings.Repeat("c", 64)
	}

	if _, err := BuildHomebrewInput("attila/focalytics", "v2.0.0", checksums); err == nil {
		t.Fatal("expected homebrew checksum error")
	}
	if _, err := BuildWinGetInput("attila/focalytics", "v2.0.0", checksums); err == nil {
		t.Fatal("expected winget checksum error")
	}
}

func TestParseChecksumsRejectsInvalidLines(t *testing.T) {
	_, err := ParseChecksums(strings.NewReader("not-a-valid-line"))
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseChecksumsRejectsDuplicateEntries(t *testing.T) {
	_, err := ParseChecksums(strings.NewReader(strings.Join([]string{
		strings.Repeat("d", 64) + "  focalytics_v1.2.3_linux_amd64.tar.gz",
		strings.Repeat("e", 64) + "  focalytics_v1.2.3_linux_amd64.tar.gz",
	}, "\n")))
	if err == nil {
		t.Fatal("expected duplicate checksum error")
	}
}

func TestReleaseAssetURLRequiresRepository(t *testing.T) {
	if _, err := ReleaseAssetURL("", "v1.2.3", "file.tar.gz"); err == nil {
		t.Fatal("expected repository error")
	}
}
