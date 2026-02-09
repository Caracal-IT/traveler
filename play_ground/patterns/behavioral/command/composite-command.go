package main

import "fmt"

var overdraftLimit2 = -500

type BankAccount2 struct {
	balance int
}

func (b *BankAccount2) Deposit2(amount int) {
	b.balance += amount
	fmt.Println("Deposited", amount,
		"\b, balance is now", b.balance)
}

func (b *BankAccount2) Withdraw2(amount int) bool {
	if b.balance-amount >= overdraftLimit2 {
		b.balance -= amount
		fmt.Println("Withdrew", amount,
			"\b, balance is now", b.balance)
		return true
	}
	return false
}

type Command2 interface {
	Call2()
	Undo2()
	Succeeded2() bool
	SetSucceeded2(value bool)
}

type Action2 int

const (
	Deposit2 Action2 = iota
	Withdraw2
)

type BankAccountCommand2 struct {
	account   *BankAccount2
	action    Action2
	amount    int
	succeeded bool
}

func (b *BankAccountCommand2) SetSucceeded2(value bool) {
	b.succeeded = value
}

// additional member
func (b *BankAccountCommand2) Succeeded2() bool {
	return b.succeeded
}

func (b *BankAccountCommand2) Call2() {
	switch b.action {
	case Deposit2:
		b.account.Deposit2(b.amount)
		b.succeeded = true
	case Withdraw2:
		b.succeeded = b.account.Withdraw2(b.amount)
	}
}

func (b *BankAccountCommand2) Undo2() {
	if !b.succeeded {
		return
	}
	switch b.action {
	case Deposit2:
		b.account.Withdraw2(b.amount)
	case Withdraw2:
		b.account.Deposit2(b.amount)
	}
}

type CompositeBankAccountCommand2 struct {
	commands []Command2
}

func (c *CompositeBankAccountCommand2) Succeeded2() bool {
	for _, cmd := range c.commands {
		if !cmd.Succeeded2() {
			return false
		}
	}
	return true
}

func (c *CompositeBankAccountCommand2) SetSucceeded2(value bool) {
	for _, cmd := range c.commands {
		cmd.SetSucceeded2(value)
	}
}

func (c *CompositeBankAccountCommand2) Call2() {
	for _, cmd := range c.commands {
		cmd.Call2()
	}
}

func (c *CompositeBankAccountCommand2) Undo2() {
	// undo in reverse order
	for idx := range c.commands {
		c.commands[len(c.commands)-idx-1].Undo2()
	}
}

func NewBankAccountCommand2(account *BankAccount2, action Action2, amount int) *BankAccountCommand2 {
	return &BankAccountCommand2{account: account, action: action, amount: amount}
}

type MoneyTransferCommand2 struct {
	CompositeBankAccountCommand2
	from, to *BankAccount2
	amount   int
}

func NewMoneyTransferCommand2(from *BankAccount2, to *BankAccount2, amount int) *MoneyTransferCommand2 {
	c := &MoneyTransferCommand2{from: from, to: to, amount: amount}
	c.commands = append(c.commands,
		NewBankAccountCommand2(from, Withdraw2, amount))
	c.commands = append(c.commands,
		NewBankAccountCommand2(to, Deposit2, amount))
	return c
}

func (m *MoneyTransferCommand2) Call2() {
	ok := true
	for _, cmd := range m.commands {
		if ok {
			cmd.Call2()
			ok = cmd.Succeeded2()
		} else {
			cmd.SetSucceeded2(false)
		}
	}
}

func main() {
	ba := &BankAccount2{}
	cmdDeposit := NewBankAccountCommand2(ba, Deposit2, 100)
	cmdWithdraw := NewBankAccountCommand2(ba, Withdraw2, 1000)
	cmdDeposit.Call2()
	cmdWithdraw.Call2()
	fmt.Println(ba)
	cmdWithdraw.Undo2()
	cmdDeposit.Undo2()
	fmt.Println(ba)

	from := BankAccount2{100}
	to := BankAccount2{0}
	mtc := NewMoneyTransferCommand2(&from, &to, 100) // try 1000
	mtc.Call2()

	fmt.Println("from=", from, "to=", to)

	fmt.Println("Undoing...")
	mtc.Undo2()
	fmt.Println("from=", from, "to=", to)
}
