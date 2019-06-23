package autodocs

// Page is a structure to store an individual page that has been
// processed and ready to be displayed.
type Page struct {
	// Name is the displayable identifier for the page.
	Name string `json:"name"`

	// Content contains the parsed content for this page.
	Content string `json:"content"`
}
