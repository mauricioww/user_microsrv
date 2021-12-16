package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	http_gokit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, http_endpoints HttpEndpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware)

	user_router := r.PathPrefix("/user").Subrouter()
	// user_router.Use(authMiddleware)

	user_router.Methods("GET").Path("/{id}").Handler(http_gokit.NewServer(
		http_endpoints.GetUser,
		decodeGetUserRequest,
		encodeResponse,
	))

	user_router.Methods("POST").Handler(http_gokit.NewServer(
		http_endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeResponse,
	))

	user_router.Methods("PUT").Path("/{id}").Handler(http_gokit.NewServer(
		http_endpoints.UpdateUser,
		decodeUpdateUserRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/auth").Handler(http_gokit.NewServer(
		http_endpoints.Authenticate,
		decodeAuthenticateRequest,
		encodeResponse,
	))

	return r
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		header_token := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(header_token) != 2 {
			rw.WriteHeader(http.StatusUnauthorized)
			res := map[string]string{"error": "No Auth Token!"}
			json.NewEncoder(rw).Encode(res)

		} else {
			jwt_token := header_token[1]
			token, err := jwt.Parse(jwt_token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method %v", token.Header["algo"])
				}
				return []byte("this_is_a_secret_shhh"), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println(claims)
				next.ServeHTTP(rw, r)
			} else {
				fmt.Println(err)
				rw.WriteHeader(http.StatusUnauthorized)
				res := map[string]string{"error": "Invalid Token!"}
				json.NewEncoder(rw).Encode(res)
			}

		}
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

func decodeAuthenticateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request AuthenticateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request UpdateUserRequest
	id_param := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_param)

	if err != nil {
		return nil, nil
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.UserId = id
	return request, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request GetUserRequest
	if id_param := mux.Vars(r)["id"]; id_param == "" {
		res := map[string]string{"error": "No Auth Token!"}
		return res, errors.New("No user_id")
	} else {
		if id, err := strconv.Atoi(id_param); err != nil {
			return nil, err
		} else {
			request = GetUserRequest{UserId: id}
		}
	}
	return request, nil
}

func encodeResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(rw).Encode(response)
}
