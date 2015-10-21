package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
)

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

	mlConf := memberlist.DefaultLANConfig()
	mlConf.Name = *name
	mlConf.BindAddr, mlConf.BindPort, err = parseHostPort(*bindAddr)
	if err != nil {
		log.Fatalf("bind: %s", err.Error())
	}

	if *advertiseAddr != "" {
		mlConf.AdvertiseAddr, mlConf.AdvertisePort, err = parseHostPort(*advertiseAddr)
		if err != nil {
			log.Fatalf("bind: %s", err.Error())
		}
	} else {
		mlConf.AdvertiseAddr = mlConf.BindAddr
		mlConf.AdvertisePort = mlConf.BindPort
	}

	nodes := []string{}
	for _, addr := range strings.Split(*joinAddrs, ",") {
		if a := strings.TrimSpace(addr); a != "" {
			nodes = append(nodes, a)
		}
	}

	log.Printf("node name: %s\n", mlConf.Name)
	log.Printf("binding address: %s:%d\n", mlConf.BindAddr, mlConf.BindPort)
	log.Printf("advertising address: %s:%d\n", mlConf.AdvertiseAddr, mlConf.AdvertisePort)
	log.Printf("addresses joining to: %s\n", strings.Join(nodes, ", "))

	ml, err := memberlist.Create(mlConf)
	if err != nil {
		log.Fatalf("create memberlist: %s", err.Error())
	}

	if len(nodes) > 0 {
		_, err := ml.Join(nodes)
		if err != nil {
			log.Fatalf("join memberlist: %s", err.Error())
		}
	}

	go func() {
		for {
			for _, m := range ml.Members() {
				log.Printf("member: %s %s:%d\n", m.Name, m.Addr, m.Port)
			}
			log.Println()
			time.Sleep(1 * time.Second)
		}
	}()

	shutdown := func() (exit bool, code int) {
		log.Println("shutting down...")
		if err := ml.Leave(1 * time.Second); err != nil {
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
