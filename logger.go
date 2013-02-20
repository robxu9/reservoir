package reservoir

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetPrefix("[Reservoir]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	file, err := os.Open("output.log")
	if err != nil {
		if err == os.ErrNotExist {
			file, err = os.Create("output.log")
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	fileBuffer := bufio.NewWriter(file)
	multiWriter := io.MultiWriter(os.Stdout, fileBuffer)
	log.SetOutput(multiWriter)
}

func Errorf(format string, args ...interface{}) {
	log.Printf("[e]"+format, args)
}

// Normal use
func Printf(format string, args ...interface{}) {
	log.Printf("[i] "+format, args)
}

// Prints and panics
func Panicf(format string, args ...interface{}) {
	log.Panicf("[p] "+format, args)
}

// Prints and exits with os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	log.Fatalf("[f] "+format, args)
}
