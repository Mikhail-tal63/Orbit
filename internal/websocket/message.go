package websocket

import "encoding/json"

func MarshalEvent(eventType string, field map[string]any) ([]byte, error) {
	out := make(map[string]any, len(field)+1)
	for i, v := range field {
		out[i] = v
	}
	out["type"] = eventType
	return json.Marshal(out)
}
