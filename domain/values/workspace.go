package values

type (
	WorkspaceName string
)

func NewWorkspaceName(name string) WorkspaceName {
	return WorkspaceName(name)
}
