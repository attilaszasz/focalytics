package app

import "io"

type ScanRequest struct {
	ArchiveRoot string
	Interactive bool
	Stdout      io.Writer
	Stderr      io.Writer
}
