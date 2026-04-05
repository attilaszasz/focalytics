package render

import (
	"fmt"
	"path/filepath"
	"time"
)

func defaultReportFilename(generatedAt time.Time) string {
	return fmt.Sprintf("focalytics_report_%s.html", generatedAt.Format("20060102_1504"))
}

func (s Service) reportPath(generatedAt time.Time) (string, error) {
	workingDirectory, err := s.getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	return filepath.Join(workingDirectory, defaultReportFilename(generatedAt)), nil
}
