// Code generated by "stringer -type=responseType -output=response_type_string.go"; DO NOT EDIT

package slack

import "fmt"

const _responseType_name = "ResponseTypeEphemeralResponseTypeInChannel"

var _responseType_index = [...]uint8{0, 21, 42}

func (i responseType) String() string {
	if i >= responseType(len(_responseType_index)-1) {
		return fmt.Sprintf("responseType(%d)", i)
	}
	return _responseType_name[_responseType_index[i]:_responseType_index[i+1]]
}
