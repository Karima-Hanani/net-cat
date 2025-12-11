package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type ChatServer struct {
	users     map[string]net.Conn
	mutix     sync.RWMutex
	history   []string
	historyMu sync.RWMutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		users:   make(map[string]net.Conn),
		history: make([]string, 0),
	}
}

func (cs *ChatServer) AddToHistory(message string) {
	cs.historyMu.Lock()
	defer cs.historyMu.Unlock()
	cs.history = append(cs.history, message)
}

// ///////////////////////////////////
func (cs *ChatServer) GetHistory() []string {
	cs.historyMu.RLock()
	defer cs.historyMu.RUnlock()
	HistoryCopy := make([]string, len(cs.history))
	copy(HistoryCopy, cs.history)
	return HistoryCopy
}

func (cs *ChatServer) UserNotExists(username string, conn net.Conn) bool {
	cs.mutix.Lock()
	defer cs.mutix.Unlock()
	if _, exists := cs.users[username]; exists {
		return false
	}
	cs.users[username] = conn
	return true
}

func (cs *ChatServer) RemoveUser(username string) {
	cs.mutix.Lock()
	defer cs.mutix.Unlock()
	delete(cs.users, username)
}

func (cs *ChatServer) UserCount() int {
	cs.mutix.RLock()
	defer cs.mutix.RUnlock()
	return len(cs.users)
}

// //////////////////////////////////////
func (cs *ChatServer) GetUsersList() []string {
	cs.mutix.RLock()
	defer cs.mutix.RUnlock()
	users := make([]string, len(cs.users))
	for username := range cs.users {
		users = append(users, username)
	}
	return users
}

func (cs *ChatServer) Broadcast(message, sender string) {
	cs.mutix.Lock()
	defer cs.mutix.Unlock()
	for username, conn := range cs.users {
		if username != sender {
			conn.Write([]byte(message))
		}
	}
}

func (cs *ChatServer) HandleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte(` 
            _nnnn_
           dGGGGMMb
          @p~qp~~qMb
          M|@||@) M|
          @,----.JM|
         JS^\__/  qKL
        dZP        qKRb
       dZP          qKKb
      fZP            SMMb
      HZM            MMMM
      FqM            MMMM
    __| ".        |\dS"qML
    |    '.       | '' \Zq
   _)      \.___.,|     .'
   \____   )MMMMMP|   .'
        '-'       '--'
   `))

	var username string

	for {
		conn.Write([]byte("[ENTER YOUR NAME]:"))
		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		username = strings.TrimSpace(input)
		if username == "" {
			conn.Write([]byte("Username cannot be empty. Try again.\n"))
			continue 
		}
		if cs.UserNotExists(username,conn) {
			break
		}
		conn.Write([]byte("Username already taken. Try another.\n"))
	}
	fmt.Printf("New client : %s (Total : %d )",username, cs.UserCount())
}
