//go:build integration

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/attila/focalytics/internal/app"
)

func TestExecuteIntegrationWithLocalDirectory(t *testing.T) {
	archiveRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(archiveRoot, "nested"), 0o755); err != nil {
		t.Fatalf("create nested dir: %v", err)
	}
	for path := range map[string]struct{}{
		filepath.Join(archiveRoot, "nested", "frame.jpg"):  {},
		filepath.Join(archiveRoot, "nested", "frame.xmp"):  {},
		filepath.Join(archiveRoot, "nested", "ignore.txt"): {},
		filepath.Join(archiveRoot, "cover.jpeg"):           {},
	} {
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	workDir := t.TempDir()
	originalWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("chdir workdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWorkingDirectory)
	}()

	exitCode := Execute([]string{archiveRoot}, bytes.NewBuffer(nil), stdout, stderr)
	if exitCode != app.DefaultExitPolicy().Success {
		t.Fatalf("expected success exit code, got %d", exitCode)
	}
	output := stdout.String()
	if strings.Contains(output, "cover.jpeg") || strings.Contains(output, "frame.jpg") || strings.Contains(output, "frame.xmp") {
		t.Fatalf("expected candidate output to be suppressed, got %q", output)
	}
	reportFiles, err := filepath.Glob(filepath.Join(workDir, "focalytics_report_*.html"))
	if err != nil {
		t.Fatalf("glob report files: %v", err)
	}
	if len(reportFiles) != 1 {
		t.Fatalf("expected exactly one report file, got %v", reportFiles)
	}
	sort.Strings(reportFiles)
	if !strings.Contains(output, reportFiles[0]) {
		t.Fatalf("expected stdout to mention report path %q, got %q", reportFiles[0], output)
	}
	reportContent, err := os.ReadFile(reportFiles[0])
	if err != nil {
		t.Fatalf("read report file: %v", err)
	}
	if !strings.Contains(string(reportContent), "focalytics report") {
		t.Fatalf("expected generated report content, got %q", string(reportContent))
	}
	stderrOutput := stderr.String()
	if strings.Contains(stderrOutput, "throughput=") {
		t.Fatalf("expected non-interactive stderr to stay quiet, got %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "embedded metadata unavailable") {
		t.Fatalf("expected metadata warning on stderr, got %q", stderrOutput)
	}
	if strings.Contains(stderrOutput, "run complete") {
		t.Fatalf("expected lifecycle chatter to be suppressed, got %q", stderrOutput)
	}
}

func TestExecuteIntegrationWithPhoneFilter(t *testing.T) {
	archiveRoot := t.TempDir()
	fixtures := map[string]string{
		filepath.Join(archiveRoot, "phone.jpg"):  "not-a-real-jpeg",
		filepath.Join(archiveRoot, "camera.jpg"): "not-a-real-jpeg",
		filepath.Join(archiveRoot, "phone.xmp"):  `<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" xmlns:aux="http://ns.adobe.com/exif/1.0/aux/" exif:FocalLength="6.9" exif:FNumber="1.8" exif:ExposureTime="1/120" exif:ISOSpeedRatings="50" tiff:Model="iPhone 15 Pro" aux:Lens="iPhone Lens"/></rdf:RDF></x:xmpmeta>`,
		filepath.Join(archiveRoot, "camera.xmp"): `<?xpacket begin=""?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:tiff="http://ns.adobe.com/tiff/1.0/" xmlns:aux="http://ns.adobe.com/exif/1.0/aux/" exif:FocalLength="50" exif:FocalLengthIn35mmFormat="50" exif:FNumber="2.8" exif:ExposureTime="1/125" exif:ISOSpeedRatings="200" tiff:Model="Canon EOS 5D Mark IV" aux:Lens="EF50mm"/></rdf:RDF></x:xmpmeta>`,
	}
	for path, content := range fixtures {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write fixture %q: %v", path, err)
		}
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	workDir := t.TempDir()
	originalWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("chdir workdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWorkingDirectory)
	}()

	exitCode := Execute([]string{"--ignore-phone-photos", archiveRoot}, bytes.NewBuffer(nil), stdout, stderr)
	if exitCode != app.DefaultExitPolicy().Success {
		t.Fatalf("expected success exit code, got %d", exitCode)
	}
	reportFiles, err := filepath.Glob(filepath.Join(workDir, "focalytics_report_*.html"))
	if err != nil {
		t.Fatalf("glob report files: %v", err)
	}
	if len(reportFiles) != 1 {
		t.Fatalf("expected exactly one report file, got %v", reportFiles)
	}
	if strings.Contains(stdout.String(), "Phone filter active") {
		t.Fatalf("expected filter summary to stay off stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Phone filter active: excluded 1 phone-made photo") {
		t.Fatalf("expected filter summary on stderr, got %q", stderr.String())
	}
	reportContent, err := os.ReadFile(reportFiles[0])
	if err != nil {
		t.Fatalf("read report file: %v", err)
	}
	html := string(reportContent)
	if !strings.Contains(html, "Phone filter active for gear and technical insights: excluded 1 phone-made photo.") {
		t.Fatalf("expected overview filter note, got %q", html)
	}
	if !strings.Contains(html, "Phone filter active: excluded 1 phone-made photo from this section.") {
		t.Fatalf("expected section filter note, got %q", html)
	}
	if strings.Contains(html, "iPhone 15 Pro") {
		t.Fatalf("expected phone camera to be excluded from affected analytics, got %q", html)
	}
	if !strings.Contains(html, ">2</div>") {
		t.Fatalf("expected total photos to reflect the full archive, got %q", html)
	}
}
