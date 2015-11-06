package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

type Config struct {
	Name          string
	BindAddr      string
	AdvertiseAddr string
}

type Agent struct {
	memberlist *memberlist.Memberlist
	broadcasts *memberlist.TransmitLimitedQueue
	config     *memberlist.Config
}

func NewAgent(c *Config) (*Agent, error) {
	var err error

	mlConfig := memberlist.DefaultLANConfig()
	mlConfig.Name = c.Name
	mlConfig.BindAddr, mlConfig.BindPort, err = parseHostPort(c.BindAddr)
	if err != nil {
		return nil, err
	}
	if c.AdvertiseAddr != "" {
		mlConfig.AdvertiseAddr, mlConfig.AdvertisePort, err = parseHostPort(c.AdvertiseAddr)
		if err != nil {
			return nil, err
		}
	} else {
		mlConfig.AdvertiseAddr = mlConfig.BindAddr
		mlConfig.AdvertisePort = mlConfig.BindPort
	}

	agent := &Agent{}
	agent.config = mlConfig

	mlConfig.Delegate = &Delegate{agent}
	ml, err := memberlist.Create(mlConfig)
	if err != nil {
		log.Fatalf("create memberlist: %s", err.Error())
	}

	agent.memberlist = ml
	agent.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return agent.memberlist.NumMembers()
		},
		RetransmitMult: mlConfig.RetransmitMult,
	}

	return agent, nil
}

func (a *Agent) QueueBroadcast(msg []byte) {
	a.broadcasts.QueueBroadcast(&Broadcast{msg})
}

type Delegate struct {
	agent *Agent
}

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (d *Delegate) NodeMeta(limit int) []byte {
	return []byte(fmt.Sprintf("meta: %s", d.agent.config.Name))
}

// NotifyMsg is called when a user-data message is received.
// Care should be taken that this method does not block, since doing
// so would block the entire UDP packet receive loop. Additionally, the byte
// slice may be modified after the call returns, so it should be copied if needed.
func (d *Delegate) NotifyMsg(msg []byte) {
	log.Printf("notify.msg: %s\n", msg)
}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
// The total byte size of the resulting data to send must not exceed
// the limit.
func (d *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return d.agent.broadcasts.GetBroadcasts(overhead, limit)
}

// LocalState is used for a TCP Push/Pull. This is sent to
// the remote side in addition to the membership information. Any
// data can be sent here. See MergeRemoteState as well. The `join`
// boolean indicates this is for a join instead of a push/pull.
func (d *Delegate) LocalState(join bool) []byte {
	log.Printf("LocalState join: %v\n", join)
	return []byte(fmt.Sprintf("localState: %s", string(d.agent.config.Name)))
}

// MergeRemoteState is invoked after a TCP Push/Pull. This is the
// state received from the remote side and is the result of the
// remote side's LocalState call. The 'join'
// boolean indicates this is for a join instead of a push/pull.
func (d *Delegate) MergeRemoteState(buf []byte, join bool) {
	log.Printf("MergeRemoteState buf: %s, join: %v", string(buf), join)
}

// Broadcast implements broadcast message of memberlist.
type Broadcast struct {
	msg []byte
}

// Invalidates
func (b *Broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

// Message returns data of the message.
func (b *Broadcast) Message() []byte {
	return b.msg
}

// Finished
func (b *Broadcast) Finished() {
}

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

func main() {
	var err error

	defaultName, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname: %s", err)
	}

	name := flag.String("name", defaultName, "node `name` to distinguish in a cluster")
	bindAddr := flag.String("bind", "0.0.0.0:7946", "`address` to communicate between memberlist nodes in a cluster")
	advertiseAddr := flag.String("advertise", "", "`address` to advertise to cluster")
	joinAddrs := flag.String("join", "", "`addresses` of existing nodes to join the cluster")
	flag.Parse()

	nodes := []string{}
	for _, addr := range strings.Split(*joinAddrs, ",") {
		if a := strings.TrimSpace(addr); a != "" {
			nodes = append(nodes, a)
		}
	}

	agent, err := NewAgent(&Config{
		Name:          *name,
		BindAddr:      *bindAddr,
		AdvertiseAddr: *advertiseAddr,
	})
	if err != nil {
		log.Fatalf("agent: %s", err)
	}

	if len(nodes) > 0 {
		_, err := agent.memberlist.Join(nodes)
		if err != nil {
			log.Fatalf("join memberlist: %s", err.Error())
		}
	}

	go func() {
		i := 0
		for {
			agent.QueueBroadcast([]byte(fmt.Sprintf("%d: hello from %s", i, agent.config.Name)))
			i++
			time.Sleep(1 * time.Second)
		}
	}()

	shutdown := func() (exit bool, code int) {
		log.Println("shutting down...")
		if err := agent.memberlist.Leave(1 * time.Second); err != nil {
			log.Fatalf("leave: %s", err.Error())
		}
		return true, 0
	}
	handler := &SignalHandler{
		NotifySIGINT:  shutdown,
		NotifySIGQUIT: shutdown,
	}
	handler.Wait()
}

func parseHostPort(s string) (host string, port int, err error) {
	host, ps, err := net.SplitHostPort(s)
	if err != nil {
		return
	}
	port, err = strconv.Atoi(ps)
	if err != nil {
		return
	}
	return
}
