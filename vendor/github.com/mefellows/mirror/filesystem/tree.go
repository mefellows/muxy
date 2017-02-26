package filesystem

import (
	"strings"
)

// A FileTree is a two-way linked Tree data-structure the represents
// a hierarchy of files and directories.
// At any point in the hierarchy one can navigate freely between nodes,
// and the ordering is preserved

// Default Impl
type FileTree struct {
	StdFile       File
	StdParentNode *FileTree
	StdChildNodes []*FileTree
}

func (fs *FileTree) ParentNode() *FileTree   { return fs.StdParentNode }
func (fs *FileTree) ChildNodes() []*FileTree { return fs.StdChildNodes }
func (fs *FileTree) File() File              { return fs.StdFile }

// Convert a FileTree to an Ordered ListMap
func FileTreeToMap(tree FileTree, base string) (map[string]File, error) {

	fileMap := map[string]File{}

	treeFunc := func(tree FileTree) (FileTree, error) {
		if _, present := fileMap[strings.Replace(tree.File().Path(), "\\", "/", -1)]; !present {
			path := strings.TrimPrefix(tree.File().Path(), base)
			fileMap[path] = tree.File()
		}
		return tree, nil
	}

	err := FileTreeWalk(tree, treeFunc)

	return fileMap, err
}

// Compare two file trees given a comparison function that returns true if two files are 'identical' by
// their own definition.
//
// Best we can do here is O(n) - we need to traverse 'src' and then compare 'target'
/*
func FileTreeDiff(src FileTree, target FileTree, comparators ...func(left File, right File) bool) (diff []File, err error) {
	// Prep our two trees into lists concurrently
	var leftMap map[string]File
	var rightMap map[string]File
	var done sync.WaitGroup
	done.Add(2)
	go func() {
		leftMap, err = FileTreeToMap(src)
		done.Done()
	}()
	go func() {
		rightMap, err = FileTreeToMap(target)
		done.Done()
	}()
	done.Wait()

	return FileMapDiff(leftMap, rightMap, comparators...)
}
*/

func FileMapDiff(src map[string]File, target map[string]File, comparators ...func(left File, right File) bool) (diff []File, err error) {
	// Iterate over the src list, comparing each item to the corresponding
	// match in the target Map
	diff = make([]File, 0)
	for filename, file := range src {
		rightFile := target[filename]
		// All comparators need to agree they are NOT different (false)
		for _, c := range comparators {
			if !c(file, rightFile) {
				diff = append(diff, file)
				break
			}
		}
	}

	return diff, nil
}

// Recursively walk a FileTree and run a self-type function on each node.
// Walker function is able to mutate the FileTree.
//
// Navigates the tree in a top left to bottom right fashion
func FileTreeWalk(tree FileTree, treeFunc func(tree FileTree) (FileTree, error)) error {
	if len(tree.ChildNodes()) > 0 {
		for _, node := range tree.ChildNodes() {

			// Mutate the tree and return any errors
			node, err := treeFunc(*node)
			if err != nil {
				return err
			}
			FileTreeWalk(node, treeFunc)
		}
	}
	return nil
}
