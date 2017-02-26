package filesystem

// A function that returns true iff the src and dest files are the same
// based on their definition
type FileComparator func(src File, dest File) bool

// Compares the last modified time of the File
var ModifiedComparator = func(src File, dest File) bool {
	if src.ModTime().After(dest.ModTime()) {
		return false
	}
	return true
}
