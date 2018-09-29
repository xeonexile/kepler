package kepler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Await holds current goroutine until SIGINT or SIGTERM
func Await() {
	//Initialize system signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Waiting for SIGINT or SIGTERM\n")
	sig := <-sigchan
	log.Printf("Caught signal %v: terminating\n", sig)
}
