package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var OnlineUsers = make(map[string]net.Conn)

func Broadcast(){

}

func Helper(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)

	conn.Write([]byte("ENTER USERNAME: "))
	username, err := reader.ReadString('\n')
	username = strings.Trim(username,"\n")
	OnlineUsers[username] = conn 
	if err != nil {
		fmt.Println("client disconnected before sending username")
		return
	}

	fmt.Println("New client connected:", username)

	for {
		str := fmt.Sprint("[",username,"]:")

		conn.Write([]byte(str))
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(username, "disconnected")
			return
		}
		
		fmt.Print(username, ":", message)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080") 
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server running on port 8080")
	for {
		if len(OnlineUsers) > 10 {
			continue
		}
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go Helper(conn)
	}
}
