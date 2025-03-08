package core

import (
	"fmt"
	"log"
	"sync"
)

// TicketSystem is a struct that contains the state of the ticket system
type TicketSystem struct {
	TotalTickets  int
	Mutex         sync.Mutex
	Condition     sync.Cond
	BuyChannel    chan int
	RefundChannel chan int
	QuitChannel   chan interface{}
	UserTickets   sync.Map
	WaitingUsers  sync.Map
	History       sync.Map
	Logger        *log.Logger
}

// Start the ticket system
func (t *TicketSystem) Start() {
	for {
		select {
		case id := <-t.BuyChannel:
			// create a new goroutine for each buy request
			go t.HandleBuyRequest(id)
		case id := <-t.RefundChannel:
			// create a new goroutine for each refund request
			go t.HandleRefundRequest(id)
		case data := <-t.QuitChannel:
			switch v := data.(type) {
			case int:
				t.Logger.Printf("👤 user %d quit", v)
			case bool:
				t.Logger.Println("👤 ticket system quit")
				return
			}
		}
	}
}

// HandleBuyRequest handles a buy request
func (t *TicketSystem) HandleBuyRequest(id int) {
	t.Mutex.Lock()
	// check if user has ticket
	if count, ok := t.UserTickets.Load(id); ok && count.(int) > 0 {
		t.Logger.Printf("❌ user %d already has a ticket!", id)
		t.Mutex.Unlock()
		return
	}

	for t.TotalTickets == 0 {
		t.UserTickets.Store(id, 0)
		t.WaitingUsers.Store(id, true)
		t.Logger.Printf("🔄 user %d waiting for ticket...", id)
		t.Condition.Wait()
	}

	if t.TotalTickets > 0 {
		t.TotalTickets--
		t.UserTickets.Store(id, 1)
		t.WaitingUsers.Store(id, false)
		t.Logger.Printf("✅ user %d success purchase ticket! remaining tickets: %d", id, t.TotalTickets)
		t.AddHistory(id, fmt.Sprintf("✅ user %d success purchase ticket! remaining tickets: %d", id, t.TotalTickets))
	} else {
		t.Logger.Printf("🔄 user %d still waiting for ticket...", id)
		t.WaitingUsers.Store(id, true)
	}

	t.Mutex.Unlock()
}

// HandleRefundRequest handles a refund request
func (t *TicketSystem) HandleRefundRequest(id int) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if count, _ := t.UserTickets.Load(id); count.(int) > 0 {
		t.TotalTickets++
		t.UserTickets.Store(id, count.(int)-1)
		t.Logger.Printf("🔄 user %d success refund ticket! remaining tickets: %d", id, t.TotalTickets)
		t.AddHistory(id, fmt.Sprintf("🔄 user %d success refund ticket! remaining tickets: %d", id, t.TotalTickets))
		t.Condition.Signal()
	} else {
		t.Logger.Printf("❌ user %d doesn't have a ticket to refund!", id)
	}
}

// QueryHistory queries the history for a user
func (t *TicketSystem) QueryHistory(userID int) []string {
	if hist, ok := t.History.Load(userID); ok {
		return hist.([]string)
	}
	return []string{}
}

// AddHistory adds a history record for a user
func (t *TicketSystem) AddHistory(userID int, message string) {
	var history []string
	if existingHistory, ok := t.History.Load(userID); ok {
		history = existingHistory.([]string)
	} else {
		history = []string{}
	}

	history = append(history, message)
	t.History.Store(userID, history)
}

// QueryTickets queries the ticket count for a user
func (t *TicketSystem) QueryTickets(userID int) int {
	if count, ok := t.UserTickets.Load(userID); ok {
		return count.(int)
	}
	return 0
}

// GetUsersWithoutTickets gets users without tickets
func (t *TicketSystem) GetUsersWithoutTickets() []int {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	usersWithoutTickets := []int{}

	// add waiting users
	t.WaitingUsers.Range(func(key, value interface{}) bool {
		userID := key.(int)
		isWaiting := value.(bool)
		if isWaiting {
			usersWithoutTickets = append(usersWithoutTickets, userID)
		}
		return true
	})

	// add users without tickets
	t.UserTickets.Range(func(key, value interface{}) bool {
		userID := key.(int)
		ticketCount := value.(int)

		// if user has no ticket
		if ticketCount == 0 {
			// check if user is already in the list
			alreadyIncluded := false
			for _, id := range usersWithoutTickets {
				if id == userID {
					alreadyIncluded = true
					break
				}
			}

			// if user is not in the list, add it
			if !alreadyIncluded {
				usersWithoutTickets = append(usersWithoutTickets, userID)
			}
		}
		return true
	})

	return usersWithoutTickets
}

// GetUsersWithTickets gets users with tickets
func (t *TicketSystem) GetUsersWithTickets() []int {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	usersWithTickets := []int{}

	// get users with tickets from UserTickets
	t.UserTickets.Range(func(key, value interface{}) bool {
		userID := key.(int)
		ticketCount := value.(int)

		// if user has tickets
		if ticketCount > 0 {
			usersWithTickets = append(usersWithTickets, userID)
		}
		return true
	})

	return usersWithTickets
}
