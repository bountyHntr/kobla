package chain

import (
	"fmt"
	"io"
	"kobla/blockchain/core/pb"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	defaultAddress = "localhost:8888"
	netProtocol    = "tcp"
)

//go:generate stringer -type=Command
type Command int

const (
	commandSync Command = iota
	commandGetBlock
	commandSendBlock
	commandSendTx
)

type communicationManager struct {
	mu         sync.RWMutex
	address    string
	knownNodes map[string]struct{}
	bc         *Blockchain
}

func newCommunicationManager(address, syncNode string, bc *Blockchain) *communicationManager {
	knownNodes := map[string]struct{}{
		syncNode: {},
	}

	if address == "" {
		address = defaultAddress
	}

	return &communicationManager{
		address:    address,
		knownNodes: knownNodes,
		bc:         bc,
	}
}

func (cm *communicationManager) listen() error {
	ln, err := net.Listen(netProtocol, defaultAddress)
	if err != nil {
		return err
	}

	go func() {
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.WithError(err).Panic("accept connection")
			}

			go cm.handleConnection(conn)
		}
	}()

	return nil
}

func (cm *communicationManager) handleConnection(conn net.Conn) {
	defer conn.Close()

	remote := conn.RemoteAddr().String()
	cm.addNewNode(remote)

	request, err := io.ReadAll(conn)
	if err != nil {
		log.WithError(err).
			WithField("node", remote).
			Error("read request")
	}

	command, request := parseCommand(request)

	switch command {
	case commandSync:
		err = cm.handleSync(request)
	case commandGetBlock:
	case commandSendBlock:
	case commandSendTx:
	}

	if err != nil {
		log.WithError(err).
			WithField("command", command.String()).
			Error("execute command")
	}
}

func (cm *communicationManager) sendSync() error {
	syncNode := cm.randomNode()
	conn, err := net.Dial(netProtocol, syncNode)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", syncNode, err)
	}

	request := pb.ChainStatus{
		Height:      cm.bc.lastBlock().Number,
		AddressFrom: cm.address,
	}

	data, err := proto.Marshal(&request)
	if err != nil {

	}
}

func (cm *communicationManager) handleSync(request []byte) error {
	var status pb.ChainStatus
	if err := proto.Unmarshal(request, &status); err != nil {
		return fmt.Errorf("unmarshal request: %w", err)
	}

	return nil
}

func (cm *communicationManager) addNewNode(node string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.knownNodes[node]; ok {
		return
	}

	cm.knownNodes[node] = struct{}{}
}

func (cm *communicationManager) randomNode() (node string) {
	cm.mu.RLock()
	for n := range cm.knownNodes {
		node = n
		break
	}
	cm.mu.RUnlock()

	return
}

func parseCommand(request []byte) (Command, []byte) {
	return Command(request[0]), request[1:]
}
