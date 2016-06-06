package slack

type Attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link,omitempty"`
}
