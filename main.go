package main

import (
	"net"
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

/////////////////////////////////////
func (cs *ChatServer) GetHistory() []string {
	cs.historyMu.RLock()
	defer cs.historyMu.RUnlock()
	HistoryCopy := make([]string, len(cs.history))
	copy(HistoryCopy, cs.history)
	return HistoryCopy
}

func (cs *ChatServer) AddUser(username string, conn net.Conn) bool {
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

////////////////////////////////////////
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


