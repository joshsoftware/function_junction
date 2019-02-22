package server

import (
	"fmt"
	"net/http"

	"github.com/A9u/function_junction/api"
	"github.com/A9u/function_junction/category"
	"github.com/A9u/function_junction/config"
	"github.com/A9u/function_junction/team"
	"github.com/gorilla/mux"
)

const (
	versionHeader = "Accept"
)

func initRouter(dep dependencies) (router *mux.Router) {
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())
	// TODO: add doc
	// v2 := fmt.Sprintf("application/vnd.%s.v2", config.AppName())

	router = mux.NewRouter()
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Category
	router.HandleFunc("/categories", category.Create(dep.CategoryService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/categories", category.List(dep.CategoryService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/categories/{category_id}", category.FindByID(dep.CategoryService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/categories/{category_id}", category.DeleteByID(dep.CategoryService)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.HandleFunc("/categories", category.Update(dep.CategoryService)).Methods(http.MethodPut).Headers(versionHeader, v1)

	router.HandleFunc("/events/{event_id}/teams", team.Create(dep.TeamService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams", team.List(dep.TeamService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	sh := http.StripPrefix("/docs/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/docs/").Handler(sh)
	return
}

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	api.Success(rw, http.StatusOK, api.Response{Message: "pong"})
}
