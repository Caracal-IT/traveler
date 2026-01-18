package main

// go get github.com/fatih/color

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}

		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// Keep track of which pizza we try to make
	var i = 0

	// Run forever or until we receive a quit notification
	// Try to make pizzas
	for {
		currentPizza := makePizza(i)

		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			// we tried to make a pizza (we send it to a data channel)
			case pizzaMaker.data <- *currentPizza:
				fmt.Println(currentPizza.message)
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)

				return
			}
		}
	}
}

func main() {
	// print out a message
	color.Cyan("The Pizzeria is open for business!!")
	color.Cyan("-----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run producer in the background
	go pizzeria(pizzaJob)

	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!\n", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer was really mad!\n")
			}
		} else {
			color.Cyan("We're done making pizzas.")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing the pizzeria: %s", err.Error())
			}
			break
		}
	}

	// print the ending message
	color.Cyan("We're done for the day!")
	color.Cyan("-----------------------------------")
	color.Cyan("Pizzas made: %d", pizzasMade)
	color.Cyan("Pizzas failed: %d", pizzasFailed)
	color.Cyan("Total orders: %d", total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed > 6:
		color.Red("It was not a good day...")
	case pizzasFailed > 4:
		color.Yellow("It was an ok day...")
	case pizzasFailed > 2:
		color.Yellow("It was a pretty good day...")
	default:
		color.Green("It was a great day!")
	}
}
