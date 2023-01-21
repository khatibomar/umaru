package internals

type GithubRepos []struct {
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
}

type ViewData struct {
	Data any
	Year string
}
