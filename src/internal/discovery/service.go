package discovery

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/attila/focalytics/internal/progress"
)

var supportedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".tif":  {},
	".tiff": {},
	".crw":  {},
	".dng":  {},
	".arw":  {},
	".cr2":  {},
	".cr3":  {},
	".nef":  {},
	".orf":  {},
	".raf":  {},
	".rw2":  {},
}

type Service struct {
	readDir func(string) ([]os.DirEntry, error)
	now     func() time.Time
}

func NewService() Service {
	return Service{
		readDir: os.ReadDir,
		now:     time.Now,
	}
}

func (s Service) Discover(root string, sink progress.Sink, stdout io.Writer) (Result, error) {
	start := s.now()
	result := Result{}
	if sink == nil {
		sink = progress.NoopSink{}
	}

	err := s.walk(root, root, &result, sink, stdout, start)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func (s Service) walk(root, current string, result *Result, sink progress.Sink, stdout io.Writer, start time.Time) error {
	entries, err := s.readDir(current)
	if err != nil {
		return err
	}
	result.DirectoriesSeen++
	s.publishProgress(sink, current, *result, start, "scanning directory")

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		entryPath := filepath.Join(current, entry.Name())
		info, err := entry.Info()
		if err != nil {
			s.recordWarning(result, sink, entryPath, fmt.Sprintf("unable to inspect entry: %v", err), start)
			continue
		}

		if info.Mode()&os.ModeSymlink != 0 {
			s.recordWarning(result, sink, entryPath, "skipping symlink entry", start)
			continue
		}

		if info.IsDir() {
			if err := s.walk(root, entryPath, result, sink, stdout, start); err != nil {
				s.recordWarning(result, sink, entryPath, fmt.Sprintf("unable to read child directory: %v", err), start)
				continue
			}
			continue
		}

		result.FilesSeen++
		candidate, ok := buildCandidate(root, entryPath)
		if !ok {
			s.publishProgress(sink, entryPath, *result, start, "skipping unsupported file")
			continue
		}

		result.Candidates = append(result.Candidates, candidate)
		result.CandidatesFound++
		s.publishProgress(sink, entryPath, *result, start, "candidate discovered")
	}

	return nil
}

func buildCandidate(root, path string) (Candidate, bool) {
	extension := strings.ToLower(filepath.Ext(path))
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return Candidate{}, false
	}
	relative = filepath.ToSlash(relative)
	if extension == ".xmp" {
		return Candidate{Kind: CandidateKindSidecar, Path: path, RelativePath: relative}, true
	}
	if _, ok := supportedImageExtensions[extension]; ok {
		return Candidate{Kind: CandidateKindImage, Path: path, RelativePath: relative}, true
	}

	return Candidate{}, false
}

func (s Service) publishProgress(sink progress.Sink, currentPath string, result Result, start time.Time, message string) {
	elapsed := s.now().Sub(start).Seconds()
	throughput := 0.0
	if elapsed > 0 {
		throughput = float64(result.FilesSeen) / elapsed
	}

	_ = sink.Publish(progress.Event{
		Stage:               "discovery",
		Kind:                progress.EventKindStatus,
		Message:             message,
		CurrentPath:         currentPath,
		FilesSeen:           result.FilesSeen,
		Warnings:            result.WarningsObserved,
		CandidatesFound:     result.CandidatesFound,
		DirectoriesSeen:     result.DirectoriesSeen,
		ThroughputPerSecond: throughput,
	})
}

func (s Service) recordWarning(result *Result, sink progress.Sink, path, message string, start time.Time) {
	result.Warnings = append(result.Warnings, Warning{Path: path, Message: message})
	result.WarningsObserved++
	elapsed := s.now().Sub(start).Seconds()
	throughput := 0.0
	if elapsed > 0 {
		throughput = float64(result.FilesSeen) / elapsed
	}
	_ = sink.Publish(progress.Event{
		Stage:               "discovery",
		Kind:                progress.EventKindWarning,
		Message:             message,
		CurrentPath:         path,
		FilesSeen:           result.FilesSeen,
		Warnings:            result.WarningsObserved,
		CandidatesFound:     result.CandidatesFound,
		DirectoriesSeen:     result.DirectoriesSeen,
		ThroughputPerSecond: throughput,
	})
}
