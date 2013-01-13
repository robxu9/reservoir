package main

import (
	"fmt"
	"github.com/stretchrcom/goweb/goweb"
)

func main() {
	goweb.MapFunc("/test/{func}/respond", func(c *goweb.Context) {
		fmt.Fprintf(c.ResponseWriter, "You called me with %s!", c.PathParams["func"])
	})

	goweb.ListenAndServe(":8080")
}
