package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// Основная функция
func main() {
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.DialTimeout("tcp", address, *timeoutFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка соединения: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Присоеденино к %s\n", address)

	exitChan := make(chan struct{})

	// Чтение из соединения и вывод в STDOUT
	go readFromConnection(conn, exitChan)

	// Чтение из STDIN и запись в соединение
	go writeToConnection(conn, exitChan)

	<-exitChan
	fmt.Println("Завершение")
}

// Функция для чтения по соединению
func readFromConnection(conn net.Conn, exitChan chan struct{}) {
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Сервер закрыл соединение")
			} else {
				fmt.Fprintf(os.Stderr, "Ошибка при чтении: %v\n", err)
			}
			close(exitChan)
			return
		}
		fmt.Print(string(data))
	}
}

func writeToConnection(conn net.Conn, exitChan chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Закртие")
				conn.Close()
			}
			close(exitChan)
			return
		}
		_, err = conn.Write(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка соединения: %v\n", err)
			close(exitChan)
			return
		}
	}
}
