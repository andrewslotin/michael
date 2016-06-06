package slack

type Attachment struct {
	AuthorName string `json:"author_name,omitempty"`
	Title      string `json:"title"`
	TitleLink  string `json:"title_link,omitempty"`
	Text       string `json:"text,omitempty"`
}
