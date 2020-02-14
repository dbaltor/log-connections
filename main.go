/*
* Listen for connections in a given port number and log to files all content received.
* Files named as connection_connIpAddress:connPortNumber_YYYY-mm-dd hh:mm:ss.log
* Default port number: 8080
*
* Command line:
* ./LogConn [<port number to listen to>]
*
* Author: Denis Baltor
* Date : 14.02.2020
 */
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

const PORT = 8080

func main() {
	argsWithProg := os.Args
	port := PORT
	if len(argsWithProg) == 1 {
		fmt.Println("No parameter provided. Using default port number: " + strconv.Itoa(PORT))
	} else {
		i, err := strconv.Atoi(argsWithProg[1])
		if err != nil {
			// handle error
			fmt.Println("Parameter provided is not Integer. Using default port number: " + strconv.Itoa(PORT))
			i = PORT
		}
		port = i
	}
	fmt.Println("Port: ", port)

	server, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		go saveConnToFile(connection)
	}
}

func saveConnToFile(connection net.Conn) {
	fmt.Println("A client has connected!")
	defer connection.Close()
	fmt.Println(connection.RemoteAddr().String())

	file, err := os.Create("connection_" + connection.RemoteAddr().String() + "_" + time.Now().Format("2006-01-02 15:04:05") + ".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 256)

	for {
		n, err := connection.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				return
			}
			break
		}

		_, err = file.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Connection closed: ", connection.RemoteAddr())
}
