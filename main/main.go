package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	core "ticket-booking/core"
)

func main() {
	totalUsers := 5

	// create log file
	logFile, err := os.Create("ticket-booking.log")
	if err != nil {
		fmt.Println("❌ cannot create log file!")
		return
	}
	defer logFile.Close()
	logger := log.New(logFile, "[TicketSystem] ", log.LstdFlags)

	// initialize ticket system, set 3 tickets
	system := core.TicketSystem{
		TotalTickets:  3,
		BuyChannel:    make(chan int, 10),
		RefundChannel: make(chan int, 10),
		QuitChannel:   make(chan interface{}, 1),
		UserTickets:   sync.Map{},
		WaitingUsers:  sync.Map{},
		Logger:        logger,
	}
	system.Condition = *sync.NewCond(&system.Mutex)
	wg := sync.WaitGroup{}

	go system.Start()

	// set timeout
	buyCtx, buyCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer buyCancel()

	// simulate 5 users purchasing tickets
	for i := 0; i < totalUsers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// add sleep for user 2 to test timeout
			if id == 2 {
				time.Sleep(3 * time.Second)
				system.Logger.Printf("👤 user %d slept for 3 seconds", id)
				system.UserTickets.Store(id, 0)
			}

			select {
			case <-buyCtx.Done():
				system.Logger.Printf("👤 user %d timed out", id)
				system.QuitChannel <- id
			default:
				system.BuyChannel <- id
				system.Logger.Printf("👤 user %d is purchasing ticket...", id)
			}
		}(i)
	}
	wg.Wait()

	system.Logger.Printf("Users without tickets: %v", system.GetUsersWithoutTickets())
	system.Logger.Printf("Users with tickets: %v", system.GetUsersWithTickets())

	// refund tickets from users with tickets
	wg = sync.WaitGroup{}
	// simulate 2 users refunding tickets (only users 0 and 1)
	for _, id := range system.GetUsersWithTickets()[:2] {
		wg.Add(1)
		go func(userId int) {
			defer wg.Done()

			select {
			case system.RefundChannel <- userId:
				system.Logger.Printf("👤 user %d is refunding ticket...", userId)
			default:
				system.Logger.Printf("👤 user %d is not refunding ticket...", userId)
			}
		}(id)
	}

	wg.Wait()
	time.Sleep(time.Second * 10)

	system.Logger.Println("👤 all operations completed")
	for i := 0; i < totalUsers; i++ {
		tickets := system.QueryTickets(i)
		history := system.QueryHistory(i)
		system.Logger.Printf("👤 user %d has %d tickets", i, tickets)
		system.Logger.Printf("👤 user %d has following history:", i)
		for _, record := range history {
			system.Logger.Printf("  %s", record)
		}
	}

	system.QuitChannel <- true
}
