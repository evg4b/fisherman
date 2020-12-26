---
id: expressions
title: Condition expressions
---

An expression on C link language. It allows you to define a condition for executing a script or validator.
<!-- TODO: Add additional description -->

## Operators

### Modifiers

#### Addition, concatenation `+`

If either left or right sides of the `+` operator are a `string`, then this operator will perform string concatenation and return that result. If neither are string, then both must be numeric, and this will return a numeric result.

Any other case is invalid.

#### Arithmetic `-` `*` `/` `**` `%`

`**` refers to "take to the power of". For instance, `3 ** 4` == 81.

* _Left side_: numeric
* _Right side_: numeric
* _Returns_: numeric

#### Bitwise shifts, masks `>>` `<<` `|` `&` `^`

All of these operators convert their `float64` left and right sides to `int64`, perform their operation, and then convert back.
Given how this library assumes numeric are represented (as `float64`), it is unlikely that this behavior will change, even though it may cause havoc with extremely large or small numbers.

* _Left side_: numeric
* _Right side_: numeric
* _Returns_: numeric

#### Negation `-`

Prefix only. This can never have a left-hand value.

* _Right side_: numeric
* _Returns_: numeric

#### Inversion `!`

Prefix only. This can never have a left-hand value.

* _Right side_: bool
* _Returns_: bool

#### Bitwise NOT `~`

Prefix only. This can never have a left-hand value.

* _Right side_: numeric
* _Returns_: numeric

### Logical Operators

For all logical operators, this library will short-circuit the operation if the left-hand side is sufficient to determine what to do. For instance, `true || expensiveOperation()` will not actually call `expensiveOperation()`, since it knows the left-hand side is `true`.

#### Logical AND/OR `&&` `||`

* _Left side_: bool
* _Right side_: bool
* _Returns_: bool

#### Ternary true `?`

Checks if the left side is `true`. If so, returns the right side. If the left side is `false`, returns `nil`.
In practice, this is commonly used with the other ternary operator.

* _Left side_: bool
* _Right side_: Any type.
* _Returns_: Right side or `nil`

#### Ternary false `:`

Checks if the left side is `nil`. If so, returns the right side. If the left side is non-nil, returns the left side.
In practice, this is commonly used with the other ternary operator.

* _Left side_: Any type.
* _Right side_: Any type.
* _Returns_: Right side or `nil`

#### Null coalescence `??`

Similar to the C# operator. If the left value is non-nil, it returns that. If not, then the right-value is returned.

* _Left side_: Any type.
* _Right side_: Any type.
* _Returns_: No specific type - whichever is passed to it.

### Comparators

#### Numeric/lexicographic comparators `>` `<` `>=` `<=`

If both sides are numeric, this returns the usual greater/lesser behavior that would be expected.
If both sides are string, this returns the lexicographic comparison of the strings. This uses Go's standard lexicographic compare.

* _Accepts_: Left and right side must either be both string, or both numeric.
* _Returns_: bool

#### Regex comparators `=~` `!~`

These use go's standard `regexp` flavor of regex. The left side is expected to be the candidate string, the right side is the pattern. `=~` returns whether or not the candidate string matches the regex pattern given on the right. `!~` is the inverted version of the same logic.

* _Left side_: string
* _Right side_: string
* _Returns_: bool

### Arrays

#### Separator `,`

The separator, always paired with parenthesis, creates arrays. It must always have both a left and right-hand value, so for instance `(, 0)` and `(0,)` are invalid uses of it.

Again, this should always be used with parenthesis; like `(1, 2, 3, 4)`.

#### Membership `IN`

The only operator with a text name, this operator checks the right-hand side array to see if it contains a value that is equal to the left-side value.
Equality is determined by the use of the `==` operator, and this library doesn't check types between the values. Any two values, when cast to `interface{}`, and can still be checked for equality with `==` will act as expected.

Note that you can use a parameter for the array, but it must be an `[]interface{}`.

* _Left side_: Any type.
* _Right side_: array
* _Returns_: bool
