package discovery

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/attila/focalytics/internal/progress"
)

type recordingSink struct {
	events []progress.Event
}

func (r *recordingSink) Publish(event progress.Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestDiscoverReturnsDeterministicCandidates(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "b", "later.xmp"))
	writeFile(t, filepath.Join(root, "b", "later.jpg"))
	writeFile(t, filepath.Join(root, "a", "earlier.jpeg"))
	writeFile(t, filepath.Join(root, "a", "skip.txt"))
	writeFile(t, filepath.Join(root, "a", "nested", "raw.CR3"))

	stdout := &bytes.Buffer{}
	sink := &recordingSink{}
	service := NewService()
	service.now = fixedClock()

	result, err := service.Discover(root, sink, stdout)
	if err != nil {
		t.Fatalf("expected discovery success: %v", err)
	}

	got := make([]string, 0, len(result.Candidates))
	for _, candidate := range result.Candidates {
		got = append(got, string(candidate.Kind)+":"+candidate.RelativePath)
	}
	want := []string{
		"image:a/earlier.jpeg",
		"image:a/nested/raw.CR3",
		"image:b/later.jpg",
		"sidecar:b/later.xmp",
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected candidates: got %v want %v", got, want)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no candidate stream output, got %q", stdout.String())
	}
	if result.WarningsObserved != 0 {
		t.Fatalf("expected no warnings, got %d", result.WarningsObserved)
	}
	if len(sink.events) == 0 {
		t.Fatal("expected progress events")
	}
}

func TestDiscoverSkipsSymlinkEntries(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "photo.jpg"))
	if err := os.Symlink(filepath.Join(root, "photo.jpg"), filepath.Join(root, "photo-link.jpg")); err != nil {
		t.Skipf("symlink unsupported on this system: %v", err)
	}

	service := NewService()
	service.now = fixedClock()
	result, err := service.Discover(root, &recordingSink{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("expected discovery success: %v", err)
	}
	if len(result.Candidates) != 1 {
		t.Fatalf("expected one real candidate, got %d", len(result.Candidates))
	}
	if result.WarningsObserved == 0 {
		t.Fatal("expected symlink warning")
	}
}

func TestDiscoverWarnsOnUnreadableChildAndContinues(t *testing.T) {
	service := NewService()
	service.now = fixedClock()
	service.readDir = func(path string) ([]os.DirEntry, error) {
		switch filepath.ToSlash(path) {
		case "/root":
			return []os.DirEntry{fakeDirEntry{name: "good.jpg"}, fakeDirEntry{name: "locked", dir: true}}, nil
		case "/root/locked":
			return nil, errors.New("permission denied")
		default:
			return nil, errors.New("unexpected path")
		}
	}

	stdout := &bytes.Buffer{}
	sink := &recordingSink{}
	result, err := service.Discover("/root", sink, stdout)
	if err != nil {
		t.Fatalf("expected child warning, not fatal error: %v", err)
	}
	if len(result.Candidates) != 1 {
		t.Fatalf("expected one candidate, got %d", len(result.Candidates))
	}
	if result.WarningsObserved != 1 {
		t.Fatalf("expected one warning, got %d", result.WarningsObserved)
	}
	if len(sink.events) == 0 || sink.events[len(sink.events)-1].Kind != progress.EventKindWarning {
		t.Fatalf("expected warning event, got %+v", sink.events)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no candidate output, got %q", stdout.String())
	}
}

func TestDiscoverReturnsRootReadError(t *testing.T) {
	service := NewService()
	service.readDir = func(string) ([]os.DirEntry, error) {
		return nil, errors.New("missing")
	}

	_, err := service.Discover("/root", &recordingSink{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected root read error")
	}
}

type fakeDirEntry struct {
	name string
	dir  bool
	mode os.FileMode
}

func (f fakeDirEntry) Name() string      { return f.name }
func (f fakeDirEntry) IsDir() bool       { return f.dir }
func (f fakeDirEntry) Type() os.FileMode { return f.mode }
func (f fakeDirEntry) Info() (os.FileInfo, error) {
	return fakeFileInfo(f), nil
}

type fakeFileInfo struct {
	name string
	dir  bool
	mode os.FileMode
}

func (f fakeFileInfo) Name() string { return f.name }
func (f fakeFileInfo) Size() int64  { return 0 }
func (f fakeFileInfo) Mode() os.FileMode {
	if f.mode != 0 {
		return f.mode
	}
	if f.dir {
		return os.ModeDir | 0o755
	}
	return 0o644
}
func (f fakeFileInfo) ModTime() time.Time { return time.Unix(1, 0) }
func (f fakeFileInfo) IsDir() bool        { return f.dir }
func (f fakeFileInfo) Sys() any           { return nil }

func fixedClock() func() time.Time {
	current := time.Unix(10, 0)
	return func() time.Time {
		current = current.Add(time.Second)
		return current
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
