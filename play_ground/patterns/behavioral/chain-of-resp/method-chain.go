package main

import "fmt"

type Creature2 struct {
	Name            string
	Attack, Defense int
}

func (c *Creature2) String() string {
	return fmt.Sprintf("%s (%d/%d)",
		c.Name, c.Attack, c.Defense)
}

func NewCreature2(name string, attack int, defense int) *Creature2 {
	return &Creature2{Name: name, Attack: attack, Defense: defense}
}

type Modifier2 interface {
	Add(m Modifier2)
	Handle()
}

type CreatureModifier2 struct {
	creature *Creature2
	next     Modifier2 // singly linked list
}

func (c *CreatureModifier2) Add(m Modifier2) {
	if c.next != nil {
		c.next.Add(m)
	} else {
		c.next = m
	}
}

func (c *CreatureModifier2) Handle() {
	if c.next != nil {
		c.next.Handle()
	}
}

func NewCreatureModifier2(creature *Creature2) *CreatureModifier2 {
	return &CreatureModifier2{creature: creature}
}

type DoubleAttackModifier2 struct {
	CreatureModifier2
}

func NewDoubleAttackModifier2(c *Creature2) *DoubleAttackModifier2 {
	return &DoubleAttackModifier2{CreatureModifier2{
		creature: c}}
}

type IncreasedDefenseModifier2 struct {
	CreatureModifier2
}

func NewIncreasedDefenseModifier2(
	c *Creature2) *IncreasedDefenseModifier2 {
	return &IncreasedDefenseModifier2{CreatureModifier2{
		creature: c}}
}

func (i *IncreasedDefenseModifier2) Handle() {
	if i.creature.Attack <= 2 {
		fmt.Println("Increasing",
			i.creature.Name, "\b's defense")
		i.creature.Defense++
	}
	i.CreatureModifier2.Handle()
}

func (d *DoubleAttackModifier2) Handle() {
	fmt.Println("Doubling", d.creature.Name,
		"attack...")
	d.creature.Attack *= 2
	d.CreatureModifier2.Handle()
}

type NoBonusesModifier2 struct {
	CreatureModifier2
}

func NewNoBonusesModifier2(
	c *Creature2) *NoBonusesModifier2 {
	return &NoBonusesModifier2{CreatureModifier2{
		creature: c}}
}

func (n *NoBonusesModifier2) Handle() {
	// nothing here!
}

func main() {
	goblin := NewCreature2("Goblin", 1, 1)
	fmt.Println(goblin.String())

	root := NewCreatureModifier2(goblin)

	//root.Add(NewNoBonusesModifier2(goblin))

	root.Add(NewDoubleAttackModifier2(goblin))
	root.Add(NewIncreasedDefenseModifier2(goblin))
	root.Add(NewDoubleAttackModifier2(goblin))

	// eventually process the entire chain
	root.Handle()
	fmt.Println(goblin.String())
}
