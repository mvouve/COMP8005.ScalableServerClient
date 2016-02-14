package main

import (
	"bufio"
	"container/list"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
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
	cInfo := make(chan clientInfo, 100)

	for i := 0; i < clients; i++ {
		waitGroup.Add(1)
		go client(addr, strSize, itterations, repeatition, cInfo)
	}
	audit(cInfo)

}

/*-----------------------------------------------------------------------------
-- FUNCTION:    audit
--
-- DATE:        February 13, 2016
--
-- REVISIONS:	  (DATE - DESCRIPTION)
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		audit(cInfo chan clientInfo)
-- 		 cInfo:		A channel to push client into into.
--
-- RETURNS: 		void
--
-- NOTES:			This function collects data from various go routines about connections
--						to the server. Before the program exits it saves the information to
--						an excel spreadsheet named by the current time, in the current
--            directory
------------------------------------------------------------------------------*/
func audit(cInfo chan clientInfo) {
	wait := make(chan bool)
	waitRoutine(wait)
	cList := new(list.List)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill)

	for {
		select {
		case c := <-cInfo:
			cList.PushBack(c)
		case <-wait:
			generateReport(time.Now().String(), cList)
			return
		case <-osSignals: // if for some reason wait never fires still gen a list.
			generateReport(time.Now().String(), cList)
			return
		}
	}
}

/*-----------------------------------------------------------------------------
-- FUNCTION:    waitRoutine
--
-- DATE:        February 13, 2016
--
-- REVISIONS:	  (DATE - DESCRIPTION)
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		waitRoutine(waitChan chan bool)
--	waitChan:		a channel to send a message about the program being finished.
--
-- RETURNS: 		void
--
-- NOTES:			This abstracts WaitGroup.Wait() into a channel for a select statment.
------------------------------------------------------------------------------*/
func waitRoutine(waitChan chan bool) {
	waitGroup.Wait()
	waitChan <- true
}

/*-----------------------------------------------------------------------------
-- FUNCTION:    waitRoutine
--
-- DATE:        February 13, 2016
--
-- REVISIONS:	  (DATE - DESCRIPTION)
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		waitRoutine(waitChan chan bool)
--	waitChan:		a channel to send a message about the program being finished.
--
-- RETURNS: 		void
--
-- NOTES:			This abstracts WaitGroup.Wait() into a channel for a select statment.
------------------------------------------------------------------------------*/
func parseInt(argName string, str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(argName, "must be an integer")
	}

	return num
}

/*-----------------------------------------------------------------------------
-- FUNCTION:    connect
--
-- DATE:        February 8, 2016
--
-- REVISIONS:	  February 13, 2016 - Added in auditing.
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		func connect(host string, strLen int, repeat int, itterations int, cInfo chan clientInfo)
--	    host:   The server to connect to.
--		strLen:	  The length of the string to send.
--		repeat:		The number of times to reconnect to the server.
--itterations:  The number of requests to make per connection.
--     cInfo:   The channel to pipe client info onto.
--
-- RETURNS: 		time.Duration: the average time the server took to echo.
--							int: The ammount of itterations actually completed, regardless of errors.
--
-- NOTES:			If this errors, it breaks out of echoing, therefore the input itterations
--						may not be acurate to the number of total itterations, which is returned.
------------------------------------------------------------------------------*/
func connect(conn net.Conn, strLen int, itterations int) (time.Duration, int) {
	defer waitGroup.Done()

	stopWatch := time.Now()
	str := strGen(strLen)
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	var i int
	for i = 0; i < itterations; i++ {

		_, err := conn.Write([]byte(str))
		if err != nil {
			log.Println(err)
			break
		}
		readWriter.ReadBytes('\n')
	}
	avgResponce := time.Duration(int64(time.Now().Sub(stopWatch)) / int64(itterations))

	return avgResponce, i
}

func client(host string, strLen int, repeat int, itterations int, cInfo chan clientInfo) {
	defer waitGroup.Done()
	for j := 0; j > repeat; j++ {

		conn, err := net.Dial("tcp", host)
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		avgTime, iCompleted := connect(conn, strLen, itterations)
		cInfo <- clientInfo{responseTime: avgTime, ammountOfData: strLen * iCompleted, requestsMade: iCompleted}
	}
}

/*-----------------------------------------------------------------------------
-- FUNCTION:    strGen
--
-- DATE:        February 8, 2016
--
-- REVISIONS:	  February 13, 2016 - Added in auditing.
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		func strGen(length int) string
--	  length:		the lenght of the string to generate.
--
-- RETURNS: 		string - a string of length length
--
-- NOTES:			  this generates a random alpha string of upper and lower case letters.
------------------------------------------------------------------------------*/
func strGen(length int) string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = runes[rand.Intn(len(runes))]
	}

	randomString[len(randomString)-1] = '\n'

	return string(randomString)
}
