package options

// Build holds options of compose build.
type Build struct {
	NoCache     bool
	ForceRemove bool
	Pull        bool
}

// Delete holds options of compose rm.
type Delete struct {
	RemoveVolume         bool
	BeforeDeleteCallback func([]string) bool
}

// Down holds options of compose down.
type Down struct {
	RemoveVolume  bool
	RemoveImages  ImageType
	RemoveOrphans bool
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
	Log bool
}

// ImageType defines the type of image (local, all)
type ImageType string

// Valid indicates whether the image type is valid.
func (i ImageType) Valid() bool {
	switch string(i) {
	case "", "local", "all":
		return true
	default:
		return false
	}
}
