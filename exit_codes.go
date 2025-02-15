package anchor

const (
	// OK signals the Anchor shutdown with no errors
	// after it was interrupted by the Wire
	OK = 0
	// Interrupted signals the Anchor failed to setup in the SetupTimeout provided.
	Interrupted = 1
	// SetupFailed signals the Anchor received an error during Setup.
	SetupFailed = 3
	// Internal signals the Anchor shutdown due to a Component returning an error.
	Internal = 4
)
