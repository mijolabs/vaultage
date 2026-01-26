package backup

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// ArchiveEntry represents a file to be added to the tar archive.
type ArchiveEntry struct {
	// Name is the path within the archive
	Name string
	// Data contains the file contents (for in-memory files)
	Data []byte
	// Path is the filesystem path (for on-disk files, mutually exclusive with Data)
	Path string
	// Mode is the file permission mode
	Mode fs.FileMode
}

// CreateArchive writes a tar archive to w containing the provided entries.
// Entries can be either in-memory (Data set) or from disk (Path set).
func CreateArchive(w io.Writer, entries []ArchiveEntry) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	for _, entry := range entries {
		if err := writeEntry(tw, entry); err != nil {
			return fmt.Errorf("writing %s: %w", entry.Name, err)
		}
	}

	return nil
}

// writeEntry writes a single entry to the tar archive.
func writeEntry(tw *tar.Writer, entry ArchiveEntry) error {
	if entry.Data != nil {
		return writeMemoryEntry(tw, entry)
	}
	return writeDiskEntry(tw, entry)
}

// writeMemoryEntry writes an in-memory file to the tar archive.
func writeMemoryEntry(tw *tar.Writer, entry ArchiveEntry) error {
	mode := entry.Mode
	if mode == 0 {
		mode = 0644
	}

	header := &tar.Header{
		Name:    entry.Name,
		Size:    int64(len(entry.Data)),
		Mode:    int64(mode),
		ModTime: time.Now(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	if _, err := tw.Write(entry.Data); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}

	return nil
}

// writeDiskEntry writes a file from disk to the tar archive.
// If the path is a directory, it recursively adds all files within it.
func writeDiskEntry(tw *tar.Writer, entry ArchiveEntry) error {
	info, err := os.Stat(entry.Path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", entry.Path, err)
	}

	if info.IsDir() {
		return writeDirEntry(tw, entry)
	}

	return writeFileEntry(tw, entry, info)
}

// writeFileEntry writes a regular file from disk to the tar archive.
func writeFileEntry(tw *tar.Writer, entry ArchiveEntry, info fs.FileInfo) error {
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("creating header: %w", err)
	}
	header.Name = entry.Name

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	file, err := os.Open(entry.Path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(tw, file); err != nil {
		return fmt.Errorf("copying data: %w", err)
	}

	return nil
}

// writeDirEntry recursively writes a directory and its contents to the tar archive.
func writeDirEntry(tw *tar.Writer, entry ArchiveEntry) error {
	return filepath.WalkDir(entry.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path within the archive
		relPath, err := filepath.Rel(entry.Path, path)
		if err != nil {
			return fmt.Errorf("calculating relative path: %w", err)
		}

		archivePath := entry.Name
		if relPath != "." {
			archivePath = filepath.Join(entry.Name, relPath)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("getting file info: %w", err)
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("creating header: %w", err)
		}
		header.Name = archivePath

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("writing header: %w", err)
		}

		// If it's a directory, we only need the header
		if d.IsDir() {
			return nil
		}

		// Copy file contents
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer file.Close()

		if _, err := io.Copy(tw, file); err != nil {
			return fmt.Errorf("copying data: %w", err)
		}

		return nil
	})
}
