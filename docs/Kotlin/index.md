# Kotlin

Kotlin is a statically typed programming language that runs on the Java Virtual Machine (JVM) and also compiles   
to JavaScript. It is designed to interoperate with Java and is often used for Android app development. Kotlin is   
known for its concise syntax, null safety features, and powerful type inference.

## Hallo World

```kotlin
fun main() {
    println("Hello, World!")
}    
```

## Variables

Kotlin variables can be declared using the `var` keyword for mutable variables and the `val` keyword for   
immutable variables. To declare a variable, use the `var` keyword for mutable   
variables and the `val` keyword for immutable variables. Types are inferred automatically. To declare a variable,   
of a specific type, use the type name after the variable name. Example:

```kotlin
var name = "John" // Mutable variable
val age = 30 // Immutable variable
val length : Int = 175
```

### Kotlin Types

Kotlin types are based on Java types, but with some additional features.   

The basic types in Kotlin include:
- `Byte` for bytes -128 to 127
- `Short` for shorts -32,768 to 32,767
- `Int` for integers -2,147,483,648 to 2,147,483,647
- `Long` for long integers -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
- `Float` for single-precision floating-point numbers 3.4028235E38 to 1.17549435E–38
- `Double` for floating-point numbers 1.7976931348623157E308 to 4.9E–324
- `Boolean` for true/false values True or False
- `String` for text
- `Char` for single characters
- `Array` for collections of elements
- `List` for ordered collections
- `Map` for key-value pairs
- `Set` for unordered collections of unique elements
- `Pair` for two-element tuples
- `Triple` for three-element tuples 


## Operators

Operators are used to perform operations on values. 

- + Addition
- - Subtraction
- * Multiplication
- / Division
- % Remainder
- == Equal to
- != Not equal to
- < Less than
- > Greater than
- <= Less than or equal to
- >= Greater than or equal to
- && Logical AND
- || Logical OR
- ! Logical NOT
- ++ Increment
- -- Decrement
- += Addition assignment
- -= Subtraction assignment
- *= Multiplication assignment
- /= Division assignment
- %= Remainder assignment
- === Identity
- !== Not identity 
- ?: Elvis operator
- !! Not null assertion

## Conditionals

conditionals are used to make decisions based on the value of a condition. 

- if-else
- when
- ?: Elvis operator
- !! Not null assertion
- && Logical AND
- || Logical OR
- ! Logical NOT

```kotlin
if (x > 5) {
    println("x is greater than 5")
} else {
    println("x is less than or equal to 5")
}
```

```kotlin
val result = if (x > 5) "x is greater than 5" else "x is less than or equal to 5"
```

```kotlin
val result = if (x > 5) {
    "x is greater than 5"
} else {
    "x is less than or equal to 5"
}
```

```kotlin
when (x) {
    1 -> print("x is 1")
    2 -> print("x is 2")
    else -> print("x is neither 1 nor 2")
}
```

```kotlin
val result = when (x) {
    1 -> "x is 1"
    2 -> "x is 2"
    else -> "x is neither 1 nor 2"
}
```

```kotlin
val result = when (x) {
    in 1..10 -> "x is between 1 and 10"
    in 11..20 -> "x is between 11 and 20"
    else -> "x is not between 1 and 20"
}
```

```kotlin
val result = when (x) {
    !in 1..10 -> "x is not between 1 and 10"
}
 ```

```kotlin
val result = when (x) {
    is Int -> "x is an integer"
    is String -> "x is a string"
    else -> "x is neither an integer nor a string"
}
```

## Loops

loops are used to repeat a block of code a specified number of times.
They are useful for iterating over collections, performing repetitive tasks, and controlling the flow of execution in a program.
There are several types of loops available in Kotlin, including for loops, while loops, and do-while loops.
For loops are used to iterate over a range of values or elements in a collection.
While loops are used to execute a block of code repeatedly as long as a condition is true.
Do-while loops are similar to while loops, but they guarantee that the block of code is executed at least once.
Loops can be nested to perform more complex operations and control the flow of execution in a program

```kotlin
for (i = 1; i <= 5; i++) {
    println(i)
}
```

```kotlin
for (i in 1..5) {
    println(i)
}
```

```kotlin
val numbers = listOf(1, 2, 3, 4, 5)`    
for (number in numbers) {
    println(number)
}
```

```kotlin
var i = 1
while (i <= 5) {
    println(i)
    i++
}
``` 

```kotlin
var i = 1
do {
    println(i)
    i++
} while (i <= 5)
```

## Break and Continue

The `break` and `continue` statements are used to exit or skip parts of a loop.

```kotlin
for (i in 1..10) {
    if (i == 5) {
        break
    }
    println(i)
}
```

```kotlin
for (i in 1..10) {
    if (i == 5) {
        continue
    }
    println(i)
}
````

## Arrays

Arrays are used to store a collection belonging to elements of the same type. 

```kotlin
val numbers = arrayOf(1, 2, 3, 4, 5)
```

```kotlin
val numbers = intArrayOf(1, 2, 3, 4, 5)

for (num in numbers) {
    println(num)
}

numbers.forEach { it -> println(it) }
```

## Ranges

Ranges are used to define a sequence of numbers or elements in a collection.
They are used to iterate over a range of values or elements in a collection.
The range operator `..` is used to define a range.

Close ranges are used to define a range of numbers or elements in a collection.  The range operator `..` is used to    
define a range.

Half-open ranges are used to define a range of numbers or elements in a collection. The range operator `until`   
is used to define a half-open range.

```kotlin
// Closed Range
val numbers = 1..5

// Half Open Range
val numbers = 1 until 5
```

```kotlin
for (i in 1..5) {
    println(i)
}
````

## Functions

Functions are used to define blocks of code that perform a specific task.
Functions can be defined using the `fun` keyword. Functions can have parameters and return a value.

```kotlin
fun sum(a: Int, b: Int): Int {
    return a + b
}
```

```kotlin
fun sum(a: Int, b: Int) = a + b
```

### Function Overloading

Functions can be overloaded by providing different parameter lists with the same name.

```kotlin
fun sum(a: Int, b: Int): Int {
    return a + b
}

fun sum(a: Double, b: Double): Double {
    return a + b
}
```

## OOP in Kotlin

Kotlin is a multi-paradigm language. It supports object-oriented programming (OOP) concepts such as classes, objects,   
inheritance, polymorphism, and encapsulation. Kotlin also supports functional programming concepts such as lambdas,   
higher-order functions, and immutability. Kotlin is designed to be concise, expressive, and safe, making it a   
popular choice for Android app development and other applications.

### Class

Classes are used to define data types and objects. Classes are used to create objects that represent real-world   
entities. Classes are used to define data types and objects. Classes are used to create objects that represent   
real-world entities. Classes can have properties and methods.

```kotlin
class Person(val name: String, val age: Int) {
    fun introduce() {
        println("Hi, my name is $name and I am $age years old.")
    }
}
```

#### Primary Constructor

primary constructor is used to initialize the properties of the class. 1 primary constructor is allowed in a class.

```kotlin
class Person(val name: String, val age: Int) {
    fun introduce() {
        println("Hi, my name is $name and I am $age years old.")
    }
}
```

#### Secondary Constructor

secondary constructor is used to initialize the properties of the class. 1 or more secondary constructors are 
allowed in a class. Secondary constructors are defined using the `constructor` keyword.

```kotlin
class Person(val name: String, val age: Int) {
    constructor(name: String) : this(name, 0)
    
    fun introduce() {
        println("Hi, my name is $name and I am $age years old.")
    }
}
```

#### Inheritance

Inheritance is used to create subclasses that inherit properties and methods from a superclass. Inheritance is     
used to create subclasses that inherit properties and methods from a superclass. Inheritance is used to create   
subclasses that inherit properties and methods from a superclass. Inheritance is used to create subclasses that   
inherit properties and methods from a superclass.

```kotlin
open class Animal(val name: String) {
    fun eat() {
        println("$name is eating.")
    }
} 

class Dog(name: String) : Animal(name) {
    fun bark() {
        println("$name is barking.")
    }
}
```




