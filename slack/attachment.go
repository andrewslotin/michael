package slack

import "encoding/json"

type Attachment struct {
	AuthorName string
	Title      string
	TitleLink  string
	Text       string
	Markdown   bool
}

type internalAttachment struct {
	AuthorName *string  `json:"author_name,omitempty"`
	Title      *string  `json:"title"`
	TitleLink  *string  `json:"title_link,omitempty"`
	Text       *string  `json:"text,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in"`
}

func (a Attachment) MarshalJSON() ([]byte, error) {
	v := internalAttachment{
		AuthorName: &a.AuthorName,
		Title:      &a.Title,
		TitleLink:  &a.TitleLink,
		Text:       &a.Text,
	}

	if a.Markdown {
		v.MarkdownIn = []string{"text"}
	}

	return json.Marshal(v)
}

func (a *Attachment) UnmarshalJSON(data []byte) error {
	if a == nil {
		*a = Attachment{}
	}

	v := internalAttachment{
		AuthorName: &a.AuthorName,
		Title:      &a.Title,
		TitleLink:  &a.TitleLink,
		Text:       &a.Text,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, field := range v.MarkdownIn {
		if field == "text" {
			a.Markdown = true
			break
		}
	}

	return nil
}
