# Ticket Booking System

A concurrent ticket booking system implemented in Go that simulates users purchasing and refunding tickets with thread-safe operations.

## Overview

This project demonstrates a ticket booking system with the following features:
- Limited ticket availability
- Concurrent user purchase requests
- Ticket refund capability
- Waiting queue for users when tickets are unavailable
- Request timeout handling
- Operation history tracking

## System Components

### Core Package

The core package contains the `TicketSystem` struct and its methods:

- **TicketSystem**: Main system that manages ticket inventory and user operations
- **Start()**: Runs the main event loop to handle purchase/refund requests
- **HandleBuyRequest()**: Processes ticket purchase requests
- **HandleRefundRequest()**: Processes ticket refund requests
- **QueryHistory()**: Retrieves operation history for a specific user
- **QueryTickets()**: Checks how many tickets a user has
- **GetUsersWithoutTickets()**: Lists users who don't have tickets
- **GetUsersWithTickets()**: Lists users who have tickets

### Main Package

The main package simulates multiple users interacting with the ticket system:
- Creates 5 simulated users
- Sets up a ticket system with 3 available tickets
- Implements timeout for purchase requests
- Demonstrates ticket refunding
- Logs all operations to a file

## Usage

To run the ticket booking system:

```bash
go run main/main.go
```

The system will generate a log file `ticket-booking.log` with detailed operation records.

## Simulation Flow

1. System initializes with 3 available tickets
2. 5 users attempt to purchase tickets concurrently
3. User 2 has a deliberate delay to demonstrate timeout behavior
4. After purchases, the system displays which users have tickets
5. Two users with tickets attempt to refund them
6. The system prints the final state and operation history for each user

## Concurrency Features

- Uses goroutines for concurrent request handling
- Implements mutex locks for thread safety
- Uses condition variables for waiting queue management
- Employs channels for communication between components
- Utilizes sync.Map for thread-safe data storage
- Implements context with timeout for request cancellation

## Log Output

The system logs all operations to `ticket-booking.log`, including:
- Ticket purchases
- Ticket refunds
- Waiting queue status
- Timeout events
- User operation history
