package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const seatingCapacity = 15
const arrivalRate = 150
const cutDuration = 1000 * time.Millisecond
const timeOpen = 10 * time.Second

func main() {
	// Print a welcome message
	color.Yellow("Welcome to the Sleeping Barber Shop")
	color.Yellow("We have %d seats available.", seatingCapacity)
	color.Yellow("-----------------------------------")

	// Create channels if we need any
	customersChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// Create the barber shop
	shop := BarberShop{
		SeatingCapacity: seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		CustomersChan:   customersChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	// Add barbers
	shop.AddBarber("Frank")
	shop.AddBarber("Susan")
	shop.AddBarber("Pat")
	shop.AddBarber("Ann")
	shop.AddBarber("George")
	shop.AddBarber("Phil")

	// Start the barber shop as a go routine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.CloseShopForDay()
		closed <- true
	}()

	// Add customers
	i := 1

	go func() {
		for {
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Duration(randomMilliseconds) * time.Millisecond):
				shop.AddCustomer(fmt.Sprintf("Customer %d", i))
				i++
			}
		}
	}()

	// Wait for all customers to finish
	<-closed

	// Print a goodbye message
	color.Green("The shop is now closed. Goodbye!")
}
