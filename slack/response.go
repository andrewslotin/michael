package slack

type Response struct {
	ResponseType responseType `json:"response_type,omitempty"`
	Text         string       `json:"text"`
}

func NewEphemeralResponse(text string) Response {
	return Response{
		ResponseType: ResponseTypeEphemeral,
		Text:         text,
	}
}

func NewInChannelResponse(text string) Response {
	return Response{
		ResponseType: ResponseTypeInChannel,
		Text:         text,
	}
}
