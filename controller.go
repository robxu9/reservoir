package main

import (
	"encoding/json"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
	gowebhttp "github.com/stretchr/goweb/http"
	"log"
	"net/http"
	"runtime/debug"
)

func Controller_Project() {

	goweb.Map("/", func(c context.Context) error {

		return goweb.API.RespondWithData(c, goweb.DefaultHttpHandler().String())

	})

	// Map main projects view
	goweb.Map("/project", func(c context.Context) error {
		list, err := Model_GetProjectList()

		if err != nil {
			log.Printf("error in controller: %s", err.Error())
			debug.PrintStack()
			return goweb.API.RespondWithError(c, http.StatusInternalServerError, err.Error())
		}

		jsonresp, err := json.MarshalIndent(list, "", "\t")

		if err != nil {
			log.Printf("error in controller: %s", err.Error())
			debug.PrintStack()
			return goweb.API.RespondWithError(c, http.StatusInternalServerError, err.Error())
		}

		return goweb.API.RespondWithData(c, jsonresp)
	})

	goweb.Map("/project/{name}", func(c context.Context) error {
		proj, err := Model_GetProject(c.PathParam("name"))

		if err != nil {
			log.Printf("error in controller: %s", err.Error())
			debug.PrintStack()
			return goweb.API.RespondWithError(c, http.StatusInternalServerError, err.Error())
		}

		jsonresp, err := json.MarshalIndent(proj, "", "\t")

		if err != nil {
			log.Printf("error in controller: %s", err.Error())
			debug.PrintStack()
			return goweb.API.RespondWithError(c, http.StatusInternalServerError, err.Error())
		}

		return goweb.API.RespondWithData(c, jsonresp)
	})

	goweb.Map("/project/{name}", func(c context.Context) error {
		proj, err := Model_GetProject(c.PathParam("name"))

		if err != nil {
			log.Printf("error in controller: %s", err.Error())
			debug.PrintStack()
			return goweb.API.RespondWithError(c, http.StatusInternalServerError, err.Error())
		}

		jsonresp, err := json.MarshalIndent(proj, "", "\t")

		if err != nil {
			log.Printf("error in controller: %s", err.Error())
			debug.PrintStack()
			return goweb.API.RespondWithError(c, http.StatusInternalServerError, err.Error())
		}

		return goweb.API.RespondWithData(c, jsonresp)
	})

}
