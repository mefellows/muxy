package filesystem

type MockFileTree struct {
	ParentNodeFileTree *FileTree
	ChildNodesArray    []*FileTree
	FileFile           File
}

func (t *MockFileTree) ParentNode() *FileTree {
	return t.ParentNodeFileTree
}

func (t *MockFileTree) ChildNodes() []*FileTree {
	return t.ChildNodesArray
}

func (t *MockFileTree) File() File {
	return t.FileFile
}
