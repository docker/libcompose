package options

// Build holds options of compose build.
type Build struct {
	NoCache bool
}

// Delete holds options of compose rm.
type Delete struct {
	RemoveVolume bool
}

// Down holds options of compose down.
type Down struct {
	RemoveVolume bool
}

// Create holds options of compose create.
type Create struct {
	NoRecreate    bool
	ForceRecreate bool
	NoBuild       bool
	// ForceBuild bool
}

// Up holds options of compose up.
type Up struct {
	Create
}
