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
- `Float` for single-precision floating-point numbers 3.4028235E38 to 1.17549435E-38
- `Double` for floating-point numbers 1.7976931348623157E308 to 4.9E-324
- `Boolean` for true/false values 
- `String` for text
- `Char` for single characters
- `Array` for collections of elements
- `List` for ordered collections
- `Map` for key-value pairs
- `Set` for unordered collections of unique elements
- `Pair` for two-element tuples
- `Triple` for three-element tuples 

