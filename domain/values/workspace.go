package values

type (
	WorkspaceID     string
	WorkspaceName   string
	WorkspaceStatus int
)

const (
	StatusDown WorkspaceStatus = iota
	StatusUp   WorkspaceStatus = iota
)

func NewWorkspaceID(id string) WorkspaceID {
	return WorkspaceID(id)
}

func NewWorkspaceName(name string) WorkspaceName {
	return WorkspaceName(name)
}
