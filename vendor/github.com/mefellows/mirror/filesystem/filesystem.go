package filesystem

import (
	"os"
	"time"
)

// Generic File System abstraction
type FileSystem interface {
	Dir(string) ([]File, error)                           // Read the contents of a directory
	Read(File) ([]byte, error)                            // Read a File
	ReadFile(file string) (File, error)                   // Read a file and return a File
	Write(file File, data []byte, perm os.FileMode) error // Write a File
	FileTree(root File) *FileTree                         // Returns a FileTree structure of Files representing the FileSystem hierarchy
	FileMap(root File) FileMap                            // Returns a FileMap structure of Files representing a flattened FileSystem hierarchy
	MkDir(file File) error
	Delete(file string) error // Delete a file on the FileSystem
}

type FileMap map[string]File

// Simple File abstraction (based on os.FileInfo)
//
// All local and remote files will be represented as a File.
// It is up to the specific FileSystem implementation to uphold this
//
type File struct {
	FileName    string      // base name of the file
	FilePath    string      // Full path to file, including filename
	FileSize    int64       // length in bytes for regular files; system-dependent for others
	FileModTime time.Time   // modification time
	FileMode    os.FileMode // File details including perms
}

func (f File) Name() string {
	return f.FileName
}

func (f File) Path() string {
	return f.FilePath
}

func (f File) Size() int64 {
	return f.FileSize
}

func (f File) ModTime() time.Time {
	return f.FileModTime
}

func (f File) IsDir() bool {
	return f.Mode().IsDir()
}

func (f File) Mode() os.FileMode {
	return f.FileMode
}

func (f File) Sys() interface{} {
	return nil
}
