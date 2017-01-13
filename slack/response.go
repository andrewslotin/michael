package slack

type Response struct {
	Message
	ResponseType responseType `json:"response_type,omitempty"`
}

func NewEphemeralResponse(text string) *Response {
	return &Response{
		ResponseType: ResponseTypeEphemeral,
		Message:      Message{Text: text},
	}
}

func NewInChannelResponse(text string) *Response {
	return &Response{
		ResponseType: ResponseTypeInChannel,
		Message:      Message{Text: text},
	}
}
