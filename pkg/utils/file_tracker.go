package utils

import (
	"errors"
	"io/fs"
	"maps"
)

// FileTracker detects file changes by comparing Adler32 checksums across invocations.
// The first call to HasChanged for a given path stores the checksum and returns false,
// since the initial state is assumed to be already handled (e.g., by feature activation).
type FileTracker struct {
	checksums map[string]string
}

func NewFileTracker(checksums map[string]string) *FileTracker {
	stored := make(map[string]string, len(checksums))
	maps.Copy(stored, checksums)
	return &FileTracker{checksums: stored}
}

// HasChanged computes the current checksum of path, compares it with the stored value,
// and updates the stored value. Returns true only if a previous checksum existed and differs.
// Returns false for missing files (nothing to re-apply) and for first-time checks (no baseline).
func (ft *FileTracker) HasChanged(path string) (bool, error) {
	checksum, err := FileChecksum(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	prev, known := ft.checksums[path]
	ft.checksums[path] = checksum

	return known && prev != checksum, nil
}

// Checksums returns the current checksum map for persistence.
func (ft *FileTracker) Checksums() map[string]string {
	return ft.checksums
}
