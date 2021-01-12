package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"gitlab.com/hydra/forum-api/api/logger"
	"gitlab.com/hydra/forum-api/api/routers"
)

// Run would start server for the api
// @Version 1.0.0
// @Title Forum Backend API
// @Description Forum API implemented using golang
// @ContactName Naufal
// @ContactEmail primayoriko@gmail.com
// @ContactURL https://primayoriko.github.io
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License
// @Server http:/localhost:8008 LocalServer
// @Security AuthorizationHeader read write
// @SecurityScheme Authorization http bearer Input your token
func Run() {
	logger.NewLogger()
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./docs/swaggerui/"))
	r.PathPrefix("/docs").
		Subrouter().
		StrictSlash(true).
		// Path("/").
		Queries().
		Handler(http.StripPrefix("/docs/", fs))

	routers.AddUserRoutes(r)
	routers.AddThreadRoutes(r)
	routers.AddPostRoutes(r)

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
