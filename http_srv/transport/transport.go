package transport

import (
	"context"
	"encoding/json"
	"net/http"

	http_gokit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, http_endpoints HttpEndpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware)

	r.Methods("POST").Path("/user").Handler(http_gokit.NewServer(
		http_endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	))

	return r
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func encodeCreateUserResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(rw).Encode(response)
}
