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
