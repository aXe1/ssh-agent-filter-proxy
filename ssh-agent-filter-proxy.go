package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
)

type AgentProxy struct {
	parentAgent agent.Agent
}

func (proxy AgentProxy) Add(key agent.AddedKey) error {
	return proxy.parentAgent.Add(key)
}

func (proxy AgentProxy) List() ([]*agent.Key, error) {
	keysList, _ := proxy.parentAgent.List()

	var filteredKeysList []*agent.Key
	for _, key := range keysList {
		log.Print(key.Comment)
		if key.Comment == os.Args[2] {
			filteredKeysList = append(filteredKeysList, key)
		}
	}
	log.Print(len(filteredKeysList))

	return filteredKeysList, nil
}

func (proxy AgentProxy) Lock(passphrase []byte) error {
	return proxy.parentAgent.Lock(passphrase)
}

func (proxy AgentProxy) Remove(key ssh.PublicKey) error {
	return proxy.parentAgent.Remove(key)
}

func (proxy AgentProxy) RemoveAll() error {
	return proxy.parentAgent.RemoveAll()
}

func (proxy AgentProxy) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return proxy.parentAgent.Sign(key, data)
}

func (proxy AgentProxy) Signers() ([]ssh.Signer, error) {
	return proxy.parentAgent.Signers()
}

func (proxy AgentProxy) Unlock(passphrase []byte) error {
	return proxy.parentAgent.Unlock(passphrase)
}

const (
	CONN_HOST = "127.0.0.1"
	CONN_TYPE = "tcp"
)

func main() {
	connPort := os.Args[1]
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + connPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn2, err := net.Dial("unix", socket)
	if err != nil {
		dat, _ := ioutil.ReadFile(socket)
		log.Print(string(dat))
		re := regexp.MustCompile(`\A!<socket\s*>(?P<Port>\d+)\s+`)
		portStr := re.FindStringSubmatch(string(dat))[1]
		log.Print(portStr)
		log.Print("127.0.0.1:" + portStr)
		conn2, _ = net.Dial("tcp", "127.0.0.1:"+portStr)
	}

	agentClient := agent.NewClient(conn2)
	proxy := AgentProxy{agentClient}
	keyList, _ := proxy.List()
	fmt.Printf("len=%d cap=%d %v\n", len(keyList), cap(keyList), keyList)

	for _, element := range keyList {
		log.Print(element.Comment)
	}

	buf := make([]byte, 16)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("%v\n", buf)
	conn.Write(buf)
	buf = make([]byte, 12)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("%v\n", buf)
	conn.Write(buf)

	agent.ServeAgent(proxy, conn)

	conn.Close()
}
