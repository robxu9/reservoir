package reservoir

import (
	"errors"
	properties "github.com/dmotylev/goproperties"
	"github.com/ogier/pflag"
	"net"
	"os"
	"path/filepath"
)

const (
	MAJOR_VERSION = 1
	MINOR_VERSION = 1
	PATCH_VERSION = 0

	CONFIG_WORKER_DIR = "workers"
)

// Flags
var helpFlag *bool = pflag.Bool("help", false, "Show help.")
var versionFlag *bool = pflag.Bool("version", false, "Shows the current version of Reservoir.")
var verboseFlag *bool = pflag.Bool("verbose", false, "Be very verbose.")
var workerDirFlag *string = pflag.String("workerdir", CONFIG_WORKER_DIR, "Look for worker configuration files in a different directory.")

func main() {
	pflag.Parse()

	if helpFlag {
		pflag.PrintDefaults()
		return
	}

	if versionFlag {
		Printf("This is version %d.%d-%d.", MAJOR_VERSION, MINOR_VERSION, PATCH_VERSION)
		return
	}

	Printf("Starting up version %d.%d-%d...\n", MAJOR_VERSION, MINOR_VERSION, PATCH_VERSION)

	// Read configuration
	Printf("Reading worker configuration:\n")
	fileWalkErr := filepath.Walk(workerDirFlag, visit)

	if fileWalkErr != nil {
		Panicf("There seems to have been a problem reading configuration: %s\n", fileWalkErr)
	}

	Printf("Starting up scheduler...\n")
	Scheduler_Run()
	if SchedulerStatus == 0 {
		Panicf("Scheduler failed to start!\n")
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() { // Ignore directories
		return nil
	}

	Printf("\tLoading configuration %s...", path)
	props, err := properties.Load(path)

	if err != nil {
		return err
	}

	workername := props.GetString("name", "")
	workerhost := props.GetString("host", "")
	workerprocesses := props.GetUint("subworkers", 2)

	if workername == "" || workerhost == "" || workerprocesses == 0 {
		return errors.New("Cannot initialise worker: \"name\" and/or \"host\" is empty, or \"subworkers\" <= 0!")
	}

	for i := 0; i < workerprocesses; i++ {
		tcpAddr, err := net.ResolveTCPAddr("tcp", workerhost)
		if err != nil {
			return err
		}
		tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return err
		}
		Scheduler_QueueWorker(&Worker{
			workername,
			i,
			workerhost,
			tcpConn,
		})
	}

	return nil
}
