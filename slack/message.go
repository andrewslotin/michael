package slack

type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}
