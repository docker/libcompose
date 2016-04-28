package types

// BuildOptions holds options of compose build.
type BuildOptions struct {
	NoCache bool
}

// DeleteOptions holds options of compose rm.
type DeleteOptions struct {
	RemoveVolume bool
}

// DownOptions holds options of compose down.
type DownOptions struct {
	RemoveVolume bool
}

// CreateOptions holds options of compose create.
type CreateOptions struct {
	NoRecreate    bool
	ForceRecreate bool
	NoBuild       bool
	// ForceBuild bool
}

// UpOptions holds options of compose up.
type UpOptions struct {
	CreateOptions
}
