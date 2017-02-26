package filesystem

import (
	"os"
	"time"
)

type MockFileSystem struct {
	ReadBytes []byte
	ReadError error

	WriteError error

	DirFiles []File
	DirError error

	FileTreeTree  *FileTree
	FileMapMap    FileMap
	MockFile      File
	MockFileError error
}

func (fs MockFileSystem) Dir(string) ([]File, error) {
	return fs.DirFiles, fs.DirError

}

func (fs MockFileSystem) Read(File) ([]byte, error) {
	return fs.ReadBytes, fs.ReadError
}

func (fs MockFileSystem) MkDir(file File) error {
	return fs.DirError
}

func (fs MockFileSystem) ReadFile(file string) (File, error) {
	return fs.MockFile, fs.MockFileError
}

func (fs MockFileSystem) Delete(string) error {
	return fs.MockFileError
}

func (fs MockFileSystem) Write(File, []byte, os.FileMode) error {
	return fs.WriteError
}

func (fs MockFileSystem) FileMap(root File) FileMap {
	return fs.FileMapMap
}

func (fs MockFileSystem) FileTree(root File) *FileTree {
	return fs.FileTreeTree
}

type MockFile struct {
	MockName     string    // base name of the file
	MockPath     string    // base name of the file
	MockSize     int64     // length in bytes for regular files; system-dependent for others
	MockModTime  time.Time // modification time
	MockIsDir    bool      // abbreviation for Mode().IsDir()
	MockFileMode os.FileMode
	MockFileSys  interface{}
}

func (f *MockFile) Path() string {
	return f.MockPath
}

func (f *MockFile) Name() string {
	return f.MockName
}

func (f *MockFile) Size() int64 {
	return f.MockSize
}

func (f *MockFile) ModTime() time.Time {
	return f.MockModTime
}

func (f *MockFile) IsDir() bool {
	return f.MockIsDir
}

func (f *MockFile) Mode() os.FileMode {
	return f.MockFileMode
}

func (f *MockFile) Sys() interface{} {
	return f.MockFileSys
}
