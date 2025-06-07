package state

func lessws(x *Workspace, y *Workspace) bool {
	return x.Abspath < y.Abspath
}

func lesswsf(x *WorkspaceFile, y *WorkspaceFile) bool {
	return x.Name < y.Name
}
