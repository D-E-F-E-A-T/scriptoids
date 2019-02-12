package scriptoids

// A Package represents a script/binary that's been installed by Scriptoids.
type Package struct {
	Name        string
	Description string
	Version     string
	EntryPoint  string
}
