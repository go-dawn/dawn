package schema

import (
	"errors"
	"time"
)

const (
	JsonTimeFormat = "2006-01-02 15:04:05"
)

// JsonTime uses custom format to marshal
// time to json
type JsonTime time.Time

// MarshalJSON marshals JsonTime into json
func (jt JsonTime) MarshalJSON() ([]byte, error) {
	t := time.Time(jt)

	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("JsonTime.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(JsonTimeFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, JsonTimeFormat)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON unmarshal data to JsonTime
func (jt *JsonTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	t, err := time.Parse(`"`+JsonTimeFormat+`"`, string(data))

	*jt = JsonTime(t)

	return err
}

// String returns custom formatted time string
func (jt JsonTime) String() string {
	return time.Time(jt).Format(JsonTimeFormat)
}
