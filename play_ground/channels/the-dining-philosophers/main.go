package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	Name      string
	RightFork int
	LeftFork  int
}

var philosophers = []Philosopher{
	{Name: "Plato", LeftFork: 4, RightFork: 0},
	{Name: "Socrates", LeftFork: 0, RightFork: 1},
	{Name: "Aristotle", LeftFork: 1, RightFork: 2},
	{Name: "Pascal", LeftFork: 2, RightFork: 3},
	{Name: "Spinoza", LeftFork: 3, RightFork: 4},
}

var hunger = 3 // how many times to eat before full.
var eatTime = 1 * time.Second
var thinkTime = 1400 * time.Millisecond
var sleepTime = 1 * time.Second

func main() {
	fmt.Println("The Dining Philosophers")
	fmt.Println("======================")
	fmt.Println("The table is empty.")

	// Start the meal
	dine()

	fmt.Println("The table is empty.")
}

func dine() {
	var wg sync.WaitGroup
	wg.Add(len(philosophers))

	var seated sync.WaitGroup
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks
	forks := make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// Start the meal

	for _, philosopher := range philosophers {
		go diningProblem(philosopher, &wg, &seated, forks)
	}

	wg.Wait()
}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, seated *sync.WaitGroup, forks map[int]*sync.Mutex) {
	defer wg.Done()

	time.Sleep(sleepTime)

	// Seat the philosopher
	fmt.Printf("%s is seated at the table.\n", philosopher.Name)
	seated.Done()
	seated.Wait()

	// eat three times
	for i := hunger; i > 0; i-- {
		// Get a lock on both forks

		if philosopher.LeftFork > philosopher.RightFork {
			forks[philosopher.RightFork].Lock()
			fmt.Printf("\t%s picked up right fork %d.\n", philosopher.Name, philosopher.RightFork)
			forks[philosopher.LeftFork].Lock()
			fmt.Printf("\t%s picked up left fork %d.\n", philosopher.Name, philosopher.LeftFork)
		} else {
			forks[philosopher.LeftFork].Lock()
			fmt.Printf("\t%s picked up left fork %d.\n", philosopher.Name, philosopher.LeftFork)
			forks[philosopher.RightFork].Lock()
			fmt.Printf("\t%s picked up right fork %d.\n", philosopher.Name, philosopher.RightFork)
		}

		fmt.Printf("%s is eating.\n", philosopher.Name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s put down left fork %d.\n", philosopher.Name, philosopher.LeftFork)
		forks[philosopher.LeftFork].Unlock()
		fmt.Printf("\t%s put down right fork %d.\n", philosopher.Name, philosopher.RightFork)
		forks[philosopher.RightFork].Unlock()

		fmt.Printf("%s is thinking.\n", philosopher.Name)
		time.Sleep(thinkTime)
	}

	fmt.Printf("%s is satisfied.\n", philosopher.Name)
	fmt.Printf("%s left the table.\n", philosopher.Name)
}
