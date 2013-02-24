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

func AddExitTask(task func()) string {
	idstring := exitIdRS.NewRandomString(8)
	exitTasks[idstring] = task
	return idstring
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
}

func init() {
	go signalCatcher()
}

func signalCatcher() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT)
	for signal := range ch {
		log.Printf("received %s, exiting.", signal.String())
		exitHandler.OnExit()
		os.Exit(0)
	}
}

func SetExitHandler(handler ExitHandler) {
	exitHandler = handler
}
