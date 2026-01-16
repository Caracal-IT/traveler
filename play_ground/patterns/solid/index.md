# Solid Principles in go

To understand the SOLID principles in Go, it's essential to first grasp the acronym itself. SOLID stands for:

- **S**ingle Responsibility Principle (SRP)
- **O**pen/Closed Principle (OCP)
- **L**iskov Substitution Principle (LSP)
- **I**nterface Segregation Principle (ISP)
- **D**ependency Inversion Principle (DIP)

Each principle plays a crucial role in designing maintainable and scalable software systems. Let's delve into 
each one in detail.

## Single Responsibility Principle

- A type should have only one reason to change.
- **Separation of concerns** is a key principle in SOLID design. - Different types / packages handling different, 
 independent tasks/problems.

## Open/Closed Principle

- Types Should be open for extension but closed for modification.

## Liskov Substitution Principle

- You should be able to substitute an embedding type in place of its embedded part.

## Interface Segregation Principle

- Don't put too much into an interface; split it into smaller ones.
- **YAGNI** – You Aren't Gonna Need It.

## Dependency Inversion Principle

- High-level modules should not depend on low-level modules. Both should depend on abstractions.
- Abstractions should not depend on details. Details should depend on abstractions.
