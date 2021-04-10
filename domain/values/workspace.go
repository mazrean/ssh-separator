package values

type (
	WorkspaceName   string
	WorkspaceStatus int
)

const (
	StatusDown WorkspaceStatus = iota
	StatusUp   WorkspaceStatus = iota
)

func NewWorkspaceName(name string) WorkspaceName {
	return WorkspaceName(name)
}
