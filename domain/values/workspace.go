package values

type (
	WorkspaceID     string
	WorkspaceName   string
	WorkspaceStatus int
)

const (
	// StatusDown the status of a workspace when it is down
	StatusDown WorkspaceStatus = iota
	// StatusUp the status of a workspace when it is up
	StatusUp WorkspaceStatus = iota
)

func NewWorkspaceID(id string) WorkspaceID {
	return WorkspaceID(id)
}

func NewWorkspaceName(name string) WorkspaceName {
	return WorkspaceName(name)
}
