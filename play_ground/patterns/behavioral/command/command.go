package main

import "fmt"

var overdraftLimit1 = -500

type BankAccount1 struct {
	balance int
}

func (b *BankAccount1) Deposit1(amount int) {
	b.balance += amount
	fmt.Println("Deposited", amount,
		"\b, balance is now", b.balance)
}

func (b *BankAccount1) Withdraw1(amount int) bool {
	if b.balance-amount >= overdraftLimit1 {
		b.balance -= amount
		fmt.Println("Withdrew", amount,
			"\b, balance is now", b.balance)
		return true
	}
	return false
}

type Command1 interface {
	Call1()
	Undo1()
}

type Action1 int

const (
	Deposit1 Action1 = iota
	Withdraw1
)

type BankAccountCommand1 struct {
	account   *BankAccount1
	action    Action1
	amount    int
	succeeded bool
}

func (b *BankAccountCommand1) Call1() {
	switch b.action {
	case Deposit1:
		b.account.Deposit1(b.amount)
		b.succeeded = true
	case Withdraw1:
		b.succeeded = b.account.Withdraw1(b.amount)
	}
}

func (b *BankAccountCommand1) Undo1() {
	if !b.succeeded {
		return
	}
	switch b.action {
	case Deposit1:
		b.account.Withdraw1(b.amount)
	case Withdraw1:
		b.account.Deposit1(b.amount)
	}
}

func NewBankAccountCommand1(account *BankAccount1, action Action1, amount int) *BankAccountCommand1 {
	return &BankAccountCommand1{account: account, action: action, amount: amount}
}

func main() {
	ba := BankAccount1{}
	cmd := NewBankAccountCommand1(&ba, Deposit1, 100)
	cmd.Call1()
	cmd2 := NewBankAccountCommand1(&ba, Withdraw1, 50)
	cmd2.Call1()
	fmt.Println(ba)
}
