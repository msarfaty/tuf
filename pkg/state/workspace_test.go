package state

func lessws(x *Workspace, y *Workspace) bool {
	return x.abspath < y.abspath
}

func lesswsf(x *WorkspaceFile, y *WorkspaceFile) bool {
	return x.name < y.name
}
