package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/attila/focalytics/internal/release"
)

var (
	stdoutWriter io.Writer = os.Stdout
	stderrWriter io.Writer = os.Stderr
	exitWithCode           = os.Exit
)

func main() {
	if len(os.Args) < 2 {
		exitWithError(fmt.Errorf("expected subcommand: asset-name, verify, or metadata"))
	}

	var err error
	switch os.Args[1] {
	case "asset-name":
		err = runAssetName(os.Args[2:])
	case "verify":
		err = runVerify(os.Args[2:])
	case "metadata":
		err = runMetadata(os.Args[2:])
	default:
		err = fmt.Errorf("unknown subcommand %q", os.Args[1])
	}

	if err != nil {
		exitWithError(err)
	}
}

func runAssetName(args []string) error {
	flags := flag.NewFlagSet("asset-name", flag.ContinueOnError)
	flags.SetOutput(stderrWriter)
	version := flags.String("version", "", "release version or tag")
	goos := flags.String("goos", "", "target GOOS")
	goarch := flags.String("goarch", "", "target GOARCH")
	field := flags.String("field", "archive", "archive, binary, base, or checksums")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if *field == "checksums" {
		value, err := release.ChecksumManifestName(*version)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdoutWriter, value)
		return err
	}

	target, err := release.TargetByPlatform(*goos, *goarch)
	if err != nil {
		return err
	}

	switch *field {
	case "archive":
		value, err := release.ReleaseAssetName(*version, target)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdoutWriter, value)
		return err
	case "binary":
		_, err := fmt.Fprintln(stdoutWriter, release.BinaryFileName(target))
		return err
	case "base":
		value, err := release.ReleaseAssetBase(*version, target)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdoutWriter, value)
		return err
	default:
		return fmt.Errorf("unsupported field %q", *field)
	}
}

func runVerify(args []string) error {
	flags := flag.NewFlagSet("verify", flag.ContinueOnError)
	flags.SetOutput(stderrWriter)
	version := flags.String("version", "", "release version or tag")
	checksumPath := flags.String("checksums", "", "path to checksum manifest")
	if err := flags.Parse(args); err != nil {
		return err
	}

	checksums, err := loadChecksums(*checksumPath)
	if err != nil {
		return err
	}

	return release.VerifyChecksumContract(*version, release.DefaultTargets(), checksums)
}

func runMetadata(args []string) error {
	flags := flag.NewFlagSet("metadata", flag.ContinueOnError)
	flags.SetOutput(stderrWriter)
	repository := flags.String("repo", "", "GitHub repository in owner/name form")
	version := flags.String("version", "", "release version or tag")
	checksumPath := flags.String("checksums", "", "path to checksum manifest")
	outDir := flags.String("out", "", "output directory")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if *outDir == "" {
		return fmt.Errorf("out directory is required")
	}

	checksums, err := loadChecksums(*checksumPath)
	if err != nil {
		return err
	}

	homebrewInput, err := release.BuildHomebrewInput(*repository, *version, checksums)
	if err != nil {
		return err
	}

	wingetInput, err := release.BuildWinGetInput(*repository, *version, checksums)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return err
	}

	if err := writeJSON(filepath.Join(*outDir, "homebrew-formula.json"), homebrewInput); err != nil {
		return err
	}

	if err := writeJSON(filepath.Join(*outDir, "winget-manifests.json"), wingetInput); err != nil {
		return err
	}

	return nil
}

func loadChecksums(path string) (map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("checksums path is required")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return release.ParseChecksums(bytes.NewReader(content))
}

func writeJSON(path string, value any) error {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	content = append(content, '\n')
	return os.WriteFile(path, content, 0o644)
}

func exitWithError(err error) {
	_, _ = fmt.Fprintln(stderrWriter, err.Error())
	exitWithCode(1)
}
