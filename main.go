package main

import (
	"bufio"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func strGen(length int) string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = runes[rand.Intn(len(runes))]
	}

	randomString[len(randomString)-1] = '\n'

	return string(randomString)
}

func main() {
	if len(os.Args) < 3 {
		// todo:error
		return
	}
	strSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		// todo: error
		return
	}
	for {
		go connector(os.Args[1], (strSize + (rand.Intn(10))))
	}
}

func connector(host string, strLen int) {
	conn, err := net.Dial("tcp", os.Args[1])
	defer conn.Close()
	if err != nil {
		return
	}
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for i := 0; i < 20; i++ {

		str := strGen(strLen)
		_, err := conn.Write([]byte(str))
		if err != nil {
			break
		}
		readWriter.ReadBytes('\n')
	}
}
