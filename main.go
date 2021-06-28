package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
)

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	//favicon ignored
}

func main() {
	serverPort := flag.String("server", "", "server:port")
	concurrent := flag.Int("c",5,"concurrent go-routines")
	dur := flag.Int("d", 5, "duration in seconds")
	delay := flag.Int("delay", 100, "delay in ms between requests (within a go-routine)")

	flag.Parse()

	if *serverPort == "" {
		fmt.Printf("need server/port")
		os.Exit(0)
	}

	if *concurrent > 5000 {
		fmt.Printf("Over Limit")
		os.Exit(0)
	}

	if *dur > 3 * 60 * 60 {
		fmt.Printf("Over Duration Limit")
		os.Exit(0)
	}

	wg := sync.WaitGroup{}
	wg.Add(*concurrent)
	endTime := time.Now().Add( time.Duration(*dur) * time.Second)
	for i:=0; i< *concurrent; i++ {
		go func() {

			for time.Now().Before(endTime) {
				_, _, err := bot.PingAndList(*serverPort)
				if err != nil {
					log.Fatalf("Error while hitting %s : %s\n", *serverPort, err.Error())
				}
				<-time.After(time.Duration(*delay) * time.Millisecond)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func (s status) String() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "Server:", s.Version.Name)
	fmt.Fprintln(&sb, "Protocol:", s.Version.Protocol)
	fmt.Fprintln(&sb, "Description:", s.Description)
	fmt.Fprintf(&sb, "Players: %d/%d\n", s.Players.Online, s.Players.Max)
	for _, v := range s.Players.Sample {
		fmt.Fprintf(&sb, "- [%s] %v\n", v.Name, v.ID)
	}
	return sb.String()
}
