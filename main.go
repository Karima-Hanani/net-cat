package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MaxUsers    = 10
	DefaultPort = "8989"
	WelcomeLogo = ` 
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
`
)

type ChatServer struct {
	users     map[string]net.Conn
	mutex     sync.RWMutex
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

func (cs *ChatServer) GetHistory() []string {
	cs.historyMu.RLock()
	defer cs.historyMu.RUnlock()
	h := make([]string, len(cs.history))
	copy(h, cs.history)
	return h
}

func (cs *ChatServer) AddUser(username string, conn net.Conn) bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	if _, exists := cs.users[username]; exists {
		return false
	}
	cs.users[username] = conn
	return true
}

func (cs *ChatServer) RemoveUser(username string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.users, username)
}

func (cs *ChatServer) UserCount() int {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	return len(cs.users)
}

func (cs *ChatServer) GetUserList() []string {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	users := make([]string, 0, len(cs.users))
	for username := range cs.users {
		users = append(users, username)
	}
	return users
}

func (cs *ChatServer) Broadcast(message string, sender string) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	for username, conn := range cs.users {
		if username != sender {
			conn.Write([]byte(message))
		}
	}
}

func (cs *ChatServer) BroadcastWithPrompt(message string, sender string) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	for username, conn := range cs.users {
		if username != sender {
			now := time.Now().Format("2006-01-02 15:04:05")
			notification := fmt.Sprintf("\r\n%s[%s][%s]: ", message, now, username)
			conn.Write([]byte(notification))
		}
	}
}

func SanitizeText(s string) string {
	s = strings.ToValidUTF8(s, "�")

	var result strings.Builder
	for _, r := range s {
		if r >= 32 && r != 127 {
			result.WriteRune(r)
		} else if r == '\t' || r == '\n' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (cs *ChatServer) HandleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	conn.Write([]byte(WelcomeLogo))
	conn.Write([]byte("Welcome to TCP-Chat!\n"))

	var username string
	for {
		conn.Write([]byte("ENTER USERNAME: "))
		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		username = strings.TrimSpace(SanitizeText(input))
		if username == "" {
			conn.Write([]byte("Username cannot be empty. Try again.\n"))
			continue
		}

		if cs.AddUser(username, conn) {
			break
		}
		conn.Write([]byte("Username already taken. Try another.\n"))
	}

	fmt.Printf("New client: %s (Total: %d)\n", username, cs.UserCount())

	for _, msg := range cs.GetHistory() {
		conn.Write([]byte(msg))
	}

	joinMsg := fmt.Sprintf("%s joined the chat\n", username)
	cs.BroadcastWithPrompt(joinMsg, username)
	cs.AddToHistory(joinMsg)

	for {
		now := time.Now().Format("2006-01-02 15:04:05")
		prompt := fmt.Sprintf("[%s][%s]: ", now, username)
		conn.Write([]byte(prompt))

		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if message == "/users" {
			users := cs.GetUserList()
			conn.Write([]byte(fmt.Sprintf("Online users: %s\n", strings.Join(users, ", "))))
			continue
		}

		if message == "/quit" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}

		fullMsg := fmt.Sprintf("[%s][%s]: %s\n", now, username, message)
		fmt.Print(fullMsg)
		cs.BroadcastWithPrompt(fullMsg, username)
		cs.AddToHistory(fullMsg)
	}

	cs.RemoveUser(username)

	leaveMsg := fmt.Sprintf("%s left the chat\n", username)
	cs.BroadcastWithPrompt(leaveMsg, username)
	cs.AddToHistory(leaveMsg)
	fmt.Printf("%s disconnected (Total: %d)\n", username, cs.UserCount())
}

func main() {
	port := DefaultPort

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	if len(os.Args) == 2 {
		portNum, err := strconv.Atoi(os.Args[1])
		if err != nil || portNum < 1 || portNum > 65535 {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		port = os.Args[1]
	}

	server := NewChatServer()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on :%s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		if server.UserCount() >= MaxUsers {
			conn.Write([]byte("Server is full. Try again later.\n"))
			conn.Close()
			continue
		}

		go server.HandleClient(conn)
	}
}