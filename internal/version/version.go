package version

// Version is the application version string.
// It can be overridden at compile time via linker flags:
// go build -ldflags "-X github.com/Chintanpatel24/Matt/internal/version.Version=v1.1.0"
var Version = "v2.0.0"
