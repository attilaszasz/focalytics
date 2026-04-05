package release

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

func ParseChecksums(reader io.Reader) (map[string]string, error) {
	checksums := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid checksum line %d: %q", lineNumber, line)
		}

		assetName := fields[1]
		if _, exists := checksums[assetName]; exists {
			return nil, fmt.Errorf("duplicate checksum entry for %s", assetName)
		}

		checksums[assetName] = fields[0]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return checksums, nil
}

func VerifyChecksumContract(version string, targets []Target, checksums map[string]string) error {
	if len(checksums) == 0 {
		return fmt.Errorf("checksum manifest is empty")
	}

	expected := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		assetName, err := ReleaseAssetName(version, target)
		if err != nil {
			return err
		}
		expected[assetName] = struct{}{}
	}

	missing := make([]string, 0)
	for assetName := range expected {
		if _, ok := checksums[assetName]; !ok {
			missing = append(missing, assetName)
		}
	}

	unexpected := make([]string, 0)
	for assetName := range checksums {
		if _, ok := expected[assetName]; !ok {
			unexpected = append(unexpected, assetName)
		}
	}

	sort.Strings(missing)
	sort.Strings(unexpected)

	if len(missing) == 0 && len(unexpected) == 0 {
		return nil
	}

	parts := make([]string, 0, 2)
	if len(missing) > 0 {
		parts = append(parts, fmt.Sprintf("missing: %s", strings.Join(missing, ", ")))
	}
	if len(unexpected) > 0 {
		parts = append(parts, fmt.Sprintf("unexpected: %s", strings.Join(unexpected, ", ")))
	}

	return fmt.Errorf("checksum contract drift detected (%s)", strings.Join(parts, "; "))
}
