package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	SeatingCapacity int
	HairCutDuration time.Duration
	NumberOfBarbers int
	Open            bool
	CustomersChan   chan string
	BarbersDoneChan chan bool
}

func (shop *BarberShop) AddBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for customers", barber)

		for {
			// If there are no customers, sleep
			if len(shop.CustomersChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap", barber)
				isSleeping = true
			}

			customer, shopOpen := <-shop.CustomersChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up", customer, barber)
					isSleeping = false
				}

				shop.cutHair(barber, customer)
			} else {
				shop.sendBarberHome(barber)
				return
			}
		}

	}()
}

func (shop *BarberShop) CloseShopForDay() {
	color.Cyan("Closing shop for the day.")
	close(shop.CustomersChan)
	shop.Open = false

	for i := 0; i < shop.NumberOfBarbers; i++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	color.Green("=======================================================")
	color.Green("Shop is closed for the day, and everyone has gone home.")
}

func (shop *BarberShop) cutHair(barber, customer string) {
	color.Green("%s cuts %s's hair.", barber, customer)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is done cutting %s's hair.", barber, customer)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) AddCustomer(customer string) {
	color.Green("*** %s arrives at the shop.", customer)

	if shop.Open {
		select {
		case shop.CustomersChan <- customer:
			color.Yellow("%s takes a seat in the waiting room.", customer)
		default:
			color.Red("The shop is full, so %s leaves.", customer)
		}
	} else {
		color.Red("The shop is closed, so %s is not allowed in.", customer)
	}
}
