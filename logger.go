package reservoir

import (
	"bufio"
	"io"
	"log"
	"os"
)

func init() {
	log.SetPrefix("Reservoir: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	file, err := os.Open("output.log")
	if err != nil {
		if err == os.ErrNotExist || err == os.ErrInvalid {
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
