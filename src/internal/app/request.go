package app

import "io"

type ScanRequest struct {
	ArchiveRoot       string
	IgnorePhonePhotos bool
	Interactive       bool
	Stdout            io.Writer
	Stderr            io.Writer
}
