package github

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	URL    string `json:"html_url"`
	Author struct {
		Name string `json:"login"`
	} `json:"user"`
}
