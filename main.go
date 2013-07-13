package main

import (
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"
)

const (
	MAJOR_VERSION = 0
	MINOR_VERSION = 1
	PATCH_VERSION = 0
)

func main() {

	log.Printf("Reservoir %s.%s.%s | A GOod Build Server\n", MAJOR_VERSION, MINOR_VERSION, PATCH_VERSION)
	log.Print("")

	var serverConfig map[string]map[string]interface{}

	err := Config_GetConfig("server", &serverConfig)
	if err != nil {
		panic(err)
	}

	useFcgi := serverConfig[Config_Environment]["fcgi"].(bool)
	addr := serverConfig[Config_Environment]["listenon"].(string)
	runtime.GOMAXPROCS(serverConfig[Config_Environment]["gomaxprocs"].(int))

	log.Printf("Starting listener at %s...", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	handler := &AuthHandler{goweb.DefaultHttpHandler()}

	if useFcgi {
		log.Print("Starting in FCGI mode.")
		err = fcgi.Serve(listener, handler)
	} else {
		log.Print("Starting in HTTP mode.")
		err = http.Serve(listener, handler)
	}

	if err != nil {
		panic(err)
	}

	log.Print("Starting signal handler...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {

			// sig is a ^C, handle it

			// stop the HTTP server
			log.Print("Stopping the server...")
			listener.Close()

			log.Print("Cleaning up remaining tasks...")

			// FIXME cleanup
			os.Exit(0)
		}
	}()

	// begin the server
	log.Fatalf("Error in Serve: %s", s.Serve(listener))

}
