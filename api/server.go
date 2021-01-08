package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/hydra/forum-api/api/routers"
)

// Run would start server for the api
func Run() {
	// if err := migrations.MigrateModels(); err != nil {
	// 	log.Fatalf("Error model migration occured, %v", err)
	// } else {
	// 	fmt.Println("Models migrated")
	// }

	r := mux.NewRouter()

	routers.AddUserRoutes(r)
	routers.AddThreadRoutes(r)
	routers.AddPostRoutes(r)
	// http.Handle("/", r)

	writeTimeout, _ := strconv.Atoi(os.Getenv("WRITE_TIMEOUT_SEC"))
	readTimeout, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT_SEC"))

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("localhost:%s", os.Getenv("API_PORT")),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	fmt.Printf("Start server at localhost:%s!\n", os.Getenv("API_PORT"))
	log.Fatalln(server.ListenAndServe())
}

/*

// Adapter is an alias so I dont have to type so much.
type Adapter func(http.Handler) http.Handler

// Adapt takes Handler funcs and chains them to the main handler.
func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
    // The loop is reversed so the adapters/middleware gets executed in the same
    // order as provided in the array.
    for i := len(adapters); i > 0; i-- {
        handler = adapters[i-1](handler)
    }
    return handler
}

// RefreshToken is the main handler.
func RefreshToken(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("hello world"))
}

// ValidateRefreshToken is the middleware.
func ValidateRefreshToken(hKey string) Adapter {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            // Check if a header key exists and has a value
            if value := req.Header.Get(hKey); value == "" {
                res.WriteHeader(http.StatusForbidden)
                res.Write([]byte("invalid request token"))
                return
            }

            // Serve the next handler
            next.ServeHTTP(res, req)
        })
    }
}

// MethodLogger logs the method of the request.
func MethodLogger() Adapter {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
            log.Printf("method=%s uri=%s\n", req.Method, req.RequestURI)
            next.ServeHTTP(res, req)
        })
    }
}

func main() {
    sm := mux.NewRouter()
    getR := sm.Methods(http.MethodGet).Subrouter()
    getR.HandleFunc("/refresh-token", Adapt(
        http.HandlerFunc(RefreshToken),
        MethodLogger(),
        ValidateRefreshToken("Vikee-Request-Token"),
    ).ServeHTTP)

    srv := &http.Server{
        Handler:      sm,
        Addr:         "localhost:8888",
        WriteTimeout: 30 * time.Second,
        ReadTimeout:  30 * time.Second,
    }
    log.Fatalln(srv.ListenAndServe())
}
*/
