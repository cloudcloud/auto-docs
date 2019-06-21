package autodocs

// Config manages all configuration attributes for auto-docs.
type Config struct {
	// Git captures details about working with git.
	Git Git

	// Listen contains the host:port for listening on HTTP
	// requests incoming.
	Listen string

	// Name is the identifier and brand for auto-docs.
	Name string
}

// Git is a structure to capture information about working with
// a git repository.
type Git struct {
	// Branch is used for checking against a particular branch.
	// If not specified, master is used.
	Branch string

	// LocalPath contains a local location that is used to store
	// and interact with the repository.
	LocalPath string

	// Password is used for the username authentication.
	Password string

	// SSHKey captures a key for use with SSH authentication.
	SSHKey string

	// Timeout specifies length of time to wait before terminating
	// checking for an update.
	Timeout int

	// URI for the repository remote.
	URI string

	// Username captures the user to connect as.
	Username string

	// Period provides a way to specify the number of seconds
	// between polling git for new commits.
	Period int
}

//
