package main

import (
	"ddb/common"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/baziotis/golang-btree/btree"
)

func server(wg *sync.WaitGroup, port_index int) {
	defer wg.Done()

	port_str := strconv.Itoa(common.STARTING_PORT + port_index)

	// Make the filename based on the port
	local_keys := btree.GetNewBTree(1, port_str+".db")

	server_addr_str := ":" + port_str
	fmt.Println(server_addr_str)
	udp_addr, err := net.ResolveUDPAddr("udp", server_addr_str)
	common.PanicOnErr(err)
	conn, err := net.ListenUDP("udp", udp_addr)
	common.PanicOnErr(err)

	buffer := make([]byte, 1024)

	defer conn.Close()

	for {
		num_bytes, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic("Could not read from UDP conn")
		}

		components := common.GetComponents(string(buffer[0:num_bytes]))
		command := components[0]

		if command == "EXIT" {
			return
		}

		key := btree.Bytes(components[1])

		var data_to_send_back []byte

		switch command {
		case "INSERT":
			common.Assert(len(components) == 3)
			val := btree.Bytes(components[2])

			local_keys.Insert(key, val)
			data_to_send_back = []byte("SAVED")
		case "GET":
			if found, val := local_keys.Find(key); found {
				data_to_send_back = []byte("FOUND:" + string(val))
			} else {
				data_to_send_back = []byte("NOT_FOUND")
			}
		case "DEL":
			if found := local_keys.Delete(key); found {
				data_to_send_back = []byte("DELETED")
			} else {
				data_to_send_back = []byte("NOT_FOUND")
			}
		default:
			common.Assert(false)
		}

		_, err = conn.WriteToUDP(data_to_send_back, addr)
		common.PanicOnErr(err)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < common.NSERVERS; i++ {
		wg.Add(1)
		go server(&wg, i)
	}
	wg.Wait()
}
