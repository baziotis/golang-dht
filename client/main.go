package main

/// TODO: Make a distributed hash table. Dumb idea first:
/// Make a single client which receives keys and distributes
/// them into servers, based on the modulo of the key.
/// - Then, consistent hashing
/// - Then, Chord

import (
	"bufio"
	"ddb/common"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	connections := make([]*net.UDPConn, common.NSERVERS)

	// Initiate connections
	some_connection_failed := false
	server_addr_str := "127.0.0.1" + ":"
	for i := 0; i < common.NSERVERS; i++ {
		udp_addr, err := net.ResolveUDPAddr("udp", server_addr_str+
			strconv.Itoa(common.STARTING_PORT+i))
		some_connection_failed = err != nil
		conn, err := net.DialUDP("udp", nil, udp_addr)
		connections[i] = conn
		some_connection_failed = err != nil
		// Note: Don't defer inside a loop
	}

	defer func() {
		for i := 0; i < common.NSERVERS; i++ {
			connections[i].Write([]byte("EXIT"))
			connections[i].Close()
		}
	}()

	if some_connection_failed {
		fmt.Println("Some connection failed")
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, err := reader.ReadString('\n')
		common.PanicOnErr(err)

		remove_nl := text[:len(text)-1]

		components := common.GetComponents(remove_nl)

		command := components[0]

		if command == "EXIT" {
			return
		}

		key_str := components[1]
		key, err := strconv.Atoi(key_str)
		common.PanicOnErr(err)

		conn := connections[key%common.NSERVERS]

		to_send := []byte(remove_nl)
		_, err = conn.Write(to_send)
		common.PanicOnErr(err)

		buffer := make([]byte, 1024)
		num_bytes, _, err := conn.ReadFromUDP(buffer)
		common.PanicOnErr(err)

		str_buffer := string(buffer[0:num_bytes])
		fmt.Printf("Reply: %s\n", str_buffer)
	}
}
