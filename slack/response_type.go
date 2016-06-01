package slack

import "fmt"

type responseType uint8

//go:generate stringer -type=responseType -output=response_type_string.go
const (
	ResponseTypeEphemeral responseType = iota
	ResponseTypeInChannel
)

func (t responseType) MarshalJSON() ([]byte, error) {
	switch t {
	case ResponseTypeEphemeral:
		return []byte(`"ephemeral"`), nil
	case ResponseTypeInChannel:
		return []byte(`"in_channel"`), nil
	default:
		return nil, fmt.Errorf("invalid response type %v", t)
	}
}
