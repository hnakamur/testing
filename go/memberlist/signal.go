package main

import (
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler holds some functions which are called when a signal is
// received.
type SignalHandler struct {
	// NotifySIGHUP is called when a SIGHUP signal is received. SIGHUP signal is
	// used to that the user's terminal is disconnect. It also used to re-reading
	// configuration files.
	NotifySIGHUP func() (exit bool, code int)

	// NotifySIGINT is called when a SIGINT signal is received. SIGINT signal is
	// used to interrupt the process. It is typically sent when the user presses
	// Ctrl-C.
	NotifySIGINT func() (exit bool, code int)

	// NotifySIGTERM is called when a SIGTERM signal is received. SIGTERM is sent
	// to request terminating the process.
	NotifySIGTERM func() (exit bool, code int)

	// NotifySIGQUIT is called when a SIGQUIT signal is received. SIGQUIT is sent
	// to quit the process and request performing a core dump. It is typically
	// sent when the user presses Ctrl-\ or Ctrl-Break.
	NotifySIGQUIT func() (exit bool, code int)
}

// Wait blocks to receive signals. When a signal is received, it calls the
// handler function corresponding to the received signal. If the handler
// function returns true with exit value, it calls os.Exit with returned code.
func (h *SignalHandler) Wait() {
	schan := make(chan os.Signal, 1)
	signal.Notify(schan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		exit := false
		code := 0
		switch <-schan {
		case syscall.SIGHUP:
			if h.NotifySIGHUP != nil {
				exit, code = h.NotifySIGHUP()
			}
		case syscall.SIGINT:
			if h.NotifySIGINT != nil {
				exit, code = h.NotifySIGINT()
			} else {
				exit, code = true, 2
			}
		case syscall.SIGTERM:
			if h.NotifySIGTERM != nil {
				exit, code = h.NotifySIGTERM()
			} else {
				exit, code = true, 2
			}
		case syscall.SIGQUIT:
			if h.NotifySIGQUIT != nil {
				exit, code = h.NotifySIGQUIT()
			} else {
				exit, code = true, 2
			}
		}
		if exit {
			os.Exit(code)
		}
	}
}
