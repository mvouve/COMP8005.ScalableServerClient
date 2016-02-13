package main

import (
	"bufio"
	"container/list"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/docopt/docopt-go"
)

var waitGroup sync.WaitGroup

type clientInfo struct {
	ammountOfData int
	requestsMade  int
	responseTime  time.Duration
}

/*-----------------------------------------------------------------------------
-- FUNCTION:    main
--
-- DATE:        February 6, 2016
--
-- REVISIONS:	  February  8, 2016 - Modified for Select
--              February 12, 2015 - changed command line inputs to use docopt
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		main()
--
-- RETURNS: 		void
--
-- NOTES:			this function sets the port into non blocking mode.
------------------------------------------------------------------------------*/
func main() {
	usage := `
  Usage:
  	client_scaleable_server -i ITTERATIONS -c CLIENTS -d AMMOUNTOFDATA -r REPITITION <address>
Options:
		-i ITTERATIONS    Number of itterations to run
		-c CLIENTS        Number of clients to connect
		-d AMMOUNTOFDATA  Ammount of data to send
		-r REPITITION     Number of times to reconnect and send data`

	arguments, _ := docopt.Parse(usage, nil, false, "1.0.0", false)

	strSize := parseInt("-d DATA", arguments["-d"].(string))
	clients := parseInt("-c CLIENTS", arguments["-c"].(string))
	itterations := parseInt("-i ITTERATIONS", arguments["-i"].(string))
	repeatition := parseInt("-r REPITITION", arguments["-r"].(string))
	addr := arguments["<address>"].(string)

	for i := 0; i < clients; i++ {
		waitGroup.Add(1)
		go connect(addr, strSize, itterations, repeatition)
	}
	waitGroup.Wait()
}

func audit(cInfo chan clientInfo) {
	wait := make(chan int)
	cList := new(list.List)
	for {
		select {
		case c := <-cInfo:
			cList.PushBack(c)
		case <-wait:
			generateReport(time.Now().String(), cList)
			return
		}
	}
}

func waitRoutine(waitChan chan bool) {
	waitGroup.Wait()
	waitChan <- true
}

func parseInt(argName string, str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(argName, "must be an integer")
	}

	return num
}

func connect(host string, strLen int, repeat int, itterations int) {
	defer waitGroup.Done()
	for j := 0; j > repeat; j++ {
		conn, err := net.Dial("tcp", host)
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}

		str := strGen(strLen)
		readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		for i := 0; i < itterations; i++ {

			_, err := conn.Write([]byte(str))
			if err != nil {
				log.Println(err)
				break
			}
			readWriter.ReadBytes('\n')
		}
	}
}

func strGen(length int) string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = runes[rand.Intn(len(runes))]
	}

	randomString[len(randomString)-1] = '\n'

	return string(randomString)
}
