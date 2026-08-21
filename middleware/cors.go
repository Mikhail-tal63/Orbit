package middleware


func IsAllowedOrigin(origin string) bool {
	_, ok := allowedOrigins[origin]
	return ok
}

var allowedOrigins = map[string]struct{}{
	"http://localhost:5173": {},
	"http://localhost:8081": {},
	"http://localhost:8082": {},
	"http://localhost:3000": {},
	"http://127.0.0.1:5173": {},
	"http://127.0.0.1:8081": {},
	"http://127.0.0.1:8082": {},
	"http://127.0.0.1:3000": {},
}
