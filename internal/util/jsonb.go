package util

import "encoding/json"

// ToJSONB trả về string JSON để Postgres tự cast sang jsonb
func ToJSONB(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	return string(b)
}
