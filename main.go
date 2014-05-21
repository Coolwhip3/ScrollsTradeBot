package main

import (
	"log"
	"os"
	"time"
)

const (
    MyRoom = "tradebot"
    TradeOne = "trading-1"
    TradeTwo = "trading-2"
    GeneralOne = "general-1"
    GeneralTwo = "general-2"

func main() {
	Log.Println("Starting up")
	
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	deny(err)
	log.SetOutput(f)
	// log.SetOutput(ioutil.Discard)

	// Get email and password from the login.txt file (2 lines)
	login, err := ioutil.ReadFile("login.txt")
	deny(err)

	lines := strings.Split(string(login), "\r\n")
	if len(lines) != 2 { // try unix line endings
		lines = strings.Split(string(login), "\n")
	}
	if len(lines) != 2 {
		panic("could not read email/password from login.txt")
	}

	email, password := lines[0], lines[1]
	go startWebServer()

	startBot(email, password, "")
	for {
		startBot(email, password, "I live again!")
	}
}
func StartBot(email, password) {
	defer func() {
		log.Println("Shut bot down.")
	}()

	s, chAlive := Connect(email, password)
	log.Println(s, chAlive)

	s.JoinRoom(MyRoom)
	s.JoinRoom(TradeOne)
	s.Joinroom(TradeTwo)
	s.JoinRoom(GeneralOne)
	s.JoinRoom(GeneralTwo)
	s.Say()
	s.SendRequest()
	s.StartMessageHandling()

	for {
		timeout := time.After(time.Minute * 1)
		InnerLoop:
		for {
			select {
				case <-chAlive:
					break InnerLoop
				case <-s.chQuit:
					log.Println("Bot Quit")
					s.chQuit <- true
					return
				case <-timeout:
					log.Println("Time out")
					return
			}
		}
	}
}

// Some error handling, could be improved
func deny(err error) {
	if err != nil {
		panic(err)
	}
}
