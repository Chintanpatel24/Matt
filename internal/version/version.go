package version

// Version is the application version string.
// It is automatically injected at build time via -ldflags:
// go build -ldflags "-X github.com/Chintanpatel24/Matt/internal/version.Version=v1.1.0"
// Or updated automatically by the GitHub release workflow when a release tag is pushed.
var Version = "v1.1.1"
