package main

import (
	"log"
	"os"
	"time"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

)

const (
    MyRoom = "tradebot"
    TradeOne = "trading-1"
    TradeTwo = "trading-2"
    GeneralOne = "general-1"
    GeneralTwo = "general-2" )
   
var WTBrequests = make(map[Player]map[Card]int)
var Bot Player
var currentState *State


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
			s.joinRoomsAndSayHi()
	                s.startTradeThread(queue)
	                s.startMessageHandlingThread(queue, chKillThread)
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
        
        if HelloMessage != "" {
            s.Say(MyRoom, HelloMessage)
            }
        
        upSince := time.Now()


	for {
		timeout := time.After(time.Minute * 1)
	InnerLoop:
		for {
			select {
			case <-chAlive:
				break InnerLoop
			case <-s.chQuit:
				log.Println("!!!main QUIT!!!")
				s.chQuit <- true
				return
			case <-timeout:
				log.Println("!!!TIMEOUT!!!")
				return
			}
		}
	}
}

func (s *State) startMessageHandlingThread(queue chan<- Player, chKillThread chan bool) {
	go func() {
		messages := s.Listen()
		defer s.Shut(messages)
		for {
			select {
			case <-chKillThread:
				return
			case m := <-messages:
				s.HandleMessages(m, queue)
			}
		}
	}()
}

func (s *State) startTradeThread(queue <-chan Player) {

	go func() {
		for {
			player := <-queue
			s.Trade(player)
			if len(queue) == 0 {
				s.Say(MyRoom, "Finished trading.")
			}
		}
	}()
}

// Some error handling, could be improved
func deny(err error) {
	if err != nil {
		panic(err)
	}
}
