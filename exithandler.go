package reservoir

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var exitHandler ExitHandler = &DefaultHandler{}
var exitTasks map[string]func() = make(map[string]func())
var exitIdRS *RS = NewAlphaNumericRS()
var Reservoir_Shutdown chan bool = make(chan bool)

func AddExitTask(task func()) string {
	for {
		idstring := exitIdRS.NewRandomString(8)
		if AddDefinedExitTask(idstring, task) {
			return idstring
		}
	}
	return ""
}

func AddDefinedExitTask(name string, task func()) bool {
	_, ok := exitTasks[name]
	if !ok {
		exitTasks[name] = task
	}
	return !ok
}

func RmExitTask(id string) {
	delete(exitTasks, id)
}

type ExitHandler interface {
	OnExit()
}

type DefaultHandler struct {
}

func (self *DefaultHandler) OnExit() {
	for key, value := range exitTasks {
		log.Printf("Running Exit Task %s...", key)
		value()
	}
}

func init() {
	go signalCatcher()
}

func signalCatcher() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)
	var type string
	select {
		case signal := <-ch:
			type = signal.String()
		case <-Reservoir_Shutdown:
			type = "normal call to shutdown"
	}
	log.Printf("received \"%s\", exiting.", signal.String())
	exitHandler.OnExit()
	os.Exit(0)
}

func SetExitHandler(handler ExitHandler) {
	exitHandler = handler
}
