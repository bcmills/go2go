# Orthogonalizing Type Lists

The [type list][] interfaces in the Type Parameters [Type Parameters draft][]
lack significant properties common to other interface types — and, indeed, types
in general! — in Go as it exists today.

## Type lists are not coherent with interface types.

The draft specifies that “[i]nterface types with type lists may only be used as
constraints on type parameters. They may not be used as ordinary interface
types. The same is true of the predeclared interface type `comparable`. … This
restriction may be lifted in future language versions.” However, it would not be
possible to lift this restriction without giving type-list interfaces a
_different_ meaning as interface types than they have as type constraints.

An ordinary interface type `T` [implements][] itself: for all types `T`, all
values of `T` support the operations of `T` and are assignable to `T`. However,
a type-list interface _cannot_ “implement” itself for the purpose of satisfying
type constraints: that would provide unspecified (and perhaps inconsistent)
semantics for operators such as `+` and `<` when two variables of the type store
values with different [concrete types][], and the whole point of type-list
interfaces in the first place is to be able to safely allow those operators.

According to the Go specification, “[a] _[type][]_ determines a set of values
together with operations and methods specific to those values”. But a type list,
as defined in the draft, represents a set of _sets of_ values paired with
operations whose semantics vary per set, not a single set of values with uniform
operations. So in some sense, a type-list interface is not really even a type!

This mismatch in meaning not only makes interface types and type-lists more
complex to describe and reason about, but also limits the future extensibility
of the language. Without this semantic mismatch, the Type Parameters design
could be extended to allow type variables themselves as constraints, with the
meaning that “the argument _is_ or _implements_ the constraint type”, or perhaps
“the argument _is assignable to_ the constraint type”. However, if the meaning
of an interface varies depending on whether it is a type or a constraint the
[substitution lemma][] would no longer hold: a type-list interface passed as
parameter `Τ` and interpreted as a proper interface type would match any type
that “implements” `T` — including a type-list interface naming any subset of
`T`'s type list — but upon substituting `T`'s actual type for the constraint, it
would no longer allow those interface types.

In contrast, the built-in `comparable` constraint, if allowed as an ordinary
interface type, would have the same properties as other interface types: the
`==` and `!=` operations are defined uniformly and meaningfully for every pair
of `comparable` values, so the type `comparable` itself guarantees the
operations enabled by the `comparable` constraint.

## Type lists miss the opportunity to address general use-cases.

Type lists match on the underlying type of a type argument, but do not allow a
generic function to convert the argument type to its underlying type. That
prevents even simple type-based specializations, such as in the
`GeneralAbsDifference` example in the draft.

Type lists enable the use of operators that are not defined for interface types.
However, they fail to cover other use cases that require non-interface types
(such as
[wrappers around `atomic.Value`](https://go2goplay.golang.org/p/1bfbQ1MDy6i), or
types that embed a field of type `*T`). (Admittedly, these use-cases are rare.)

Type-list interfaces also have a clear similarity to the oft-requested
“[sum types][]” that could be used to enable safer type-switches, but because
type-lists match on underlying types instead of concrete types, they do not
directly address the use cases for sum types.

## How could we address the use-cases of type lists more orthogonally?

Type lists could be made more orthogonal in one of several ways, but we must
start by acknowledging that one of the type constraints we want to express —
namely, the constraint that a parameter must be a _non-interface_ type — is not
itself a coherent Go type.

We could make a minimal change to the design to distinguish between types and
constraints:

1.  Define that a type constraint may be either an interface type or a type-list
    constraint. Define type-list constraints as constraints (not types) using
    tokens other than `type` and `interface`.

    For example, use the token `constraint`instead:

    ```go
    constraint T {
        type X, Y, Z
        someInterface
    }
    ```

Or, we could split out the “must be a non-interface type” constraint, and
preserve as much of the current type-lists as we can as interface types:

1.  Define “type list interface types”, implemented by any type for which all
    values are known to have _a dynamic type with the same underlying type as_
    one of the types in the list. Either reject or filter out any interface
    types in the list, since an interface type cannot be the dynamic underlying
    type of any value at run time.

2.  Define that a type constraint may be either a type, or a type restricted to
    its concrete (non-interface) implementations.

Or we could break down type-lists into smaller orthogonal parts: disjunction
interfaces, underlying-type interfaces, and concrete-type constraints.

1.  Define “sum interface types”, implemented by any type that _is or
    implements_ any of the listed types.

2.  Define “defined-sum interface types”, implemented by any type whose
    underlying type is any of the listed types.

3.  Define a mapping between defined-sum interface types and the corresponding
    ordinary sum types.

4.  Define that a type constraint may be either a type, or a type restricted to
    its concrete (non-interface) implementations.

Under the orthogonal option, the interface types would be defined as follows.

### Sum interface types

A sum-interface type (or “sum type”) is declared using the syntax `type = T1, …,
Tn`, and is implemented by any type that _is_ or _implements_ at least one of
the types in the list. (When a sum type `T` is implemented by a concrete type
`R`, we say that `R` is “in“ the sum `T`.)

The method set of a sum type is the intersection of the method sets of each type
in the sum. No other methods may be defined and no other interfaces may be
embedded, because the interfaces and methods implemented by the types in the sum
are already known.

The zero value of a sum type is the `nil` interface value, even if none of the
types in the list is itself an interface type.

A type switch or type assertion on a variable of a sum type may use only the
types in the sum and interface types _implemented by_ at least one type in the
sum. To allow lists to be expanded over time, a `default` case is permitted even
if the switch is exhaustive.

A sum type is assignable to any interface implemented by all of the types in the
sum.

If all of the types in the sum are convertible to a type `T`, then the sum type
is also convertible to `T`. If `T` is a concrete type, the `nil` interface value
converts to the zero-value of the type.

A sum type embedded in another interface `S` restricts `S` to the types in the
sum that also implement the remainder of `S`. (If multiple sum interface types
are embedded in an interface, the types in the sums are intersected.)

### Defined-sum interface types

A defined-sum interface type (or “defined sum type”), declared using the syntax
`type T1, … Tn`, is a sum type that includes in the sum every type whose
underlying type is one of the listed types (`T1` through `Tn`). (An ordinary sum
type includes a fixed set of types, but the set of types included in a
defined-sum type is unbounded.)

The underlying types listed in a defined sum type declaration must be
predeclared boolean, numeric, or string types, type literals, or sum types
comprising the same.

A defined-sum interface type _may_ include additional methods. A defined-sum
interface may embed other (ordinary and defined) sum interfaces.

A type `R` implements a defined sum type `T` if the underlying type of `R` is in
the underlying-type list of `T`, the method set of `R` includes all of the
methods declared in `T`, and `R` implements all of the interfaces embedded in
`T`.

### The `Underlying` type alias and `underlying` function

The built-in type `Underlying(type T)` is an alias for the type encompassing the
underlying types of the _concrete values of_ any type `T`:

*   If `T` is a sum interface type (including a defined-sum interface type),
    `Underlying(T)` is the sum interface type containing `Underlying(Tᵢ)` for
    each `Tᵢ` in `T`.
*   If `T` is any other interface type, `Underlying(T)` is the empty interface
    type. (That is, the set of possible underlying types for the values of type
    `T` is unrestricted.)
*   Otherwise, `T` is a concrete type, and `Underlying(T)` is its underlying
    type.

The built-in function `underlying` converts a value of any type to a value of
its underlying type.

```go
func underlying(type T)(x T) Underlying(T)
```

### Constraining type parameters to concrete types

The last piece of the puzzle is a constraint that restricts a type parameter to
only [concrete types][]. If all of the concrete types in a sum interface type
support a given operation, then a function with a type parameter constrained to
those concrete types can safely use that operation.

We could imagine a lot of options for the syntax of such a constraint. For the
purpose of this document, I suggest the keyword `in` immediately before the
parameter's type constraint, meaning “the type must be one of the concrete
dynamic types for values stored _in_ the constraint type”. (However, note that
this constraint has a well-defined meaning even for types that are not sum
types.)

```go
func
```

## Examples

With the orthogonal approach, the examples in the draft design become a bit more
verbose, but more explicit — and, notably, the GeneralAbsDifference function
becomes possible to implement.

```go
package constraints

type Ordered interface {
    type int, int8, …, string
}
```

```go
func Smallest(type T in constraints.Ordered)(s []T) T {
    …
}
```

```go
// StringableSignedInteger is an interface type implemented by any
// type that both 1) is defined as a signed integer type, and
// 2) has a String method.
type StringableSignedInteger interface {
    type int, int8, int16, int32, int64
    String() string
}
```

```go
// SliceConstraint is an interface type implemented by any slice type with
// an element type identical to the type parameter.
type SliceConstraint(type T) interface {
    type []T
}
```

```go
// integer is a type implemented by any defined integer type.
type integer interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr
}

// Convert converts a value of type From to the type To.
// To must be a specific concrete integer type.
// From may be any type that implements the integer interface,
// including any sum type whose members are all integer types.
func Convert(type To in integer, From integer)(x From) To {
    to := To(x)
    if From(to) != x {
        panic("conversion out of range")
    }
    return to
}
```

```go
type integer interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr
}

func Add10(type T in integer)(s []T) {
    for i, v := range s {
        s[i] = v + 10 // OK: 10 can convert to any concrete integer type
    }
}

// This function is INVALID.
func Add1024(type T in integer)(s []T) {
    for i, v := range s {
        s[i] = v + 1024 // INVALID: 1024 not permitted by types int8 and uint8 in integer
    }
}

// This function is INVALID.
func Add10Interface(type T integer)(s []T) {
    for i, v := range s {
        s[i] = v + 10 // INVALID: operation + not permitted for type T that may be an interface type (use "T in integer" to restrict to concrete types)
    }
}

```

```go
type Numeric interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr,
        float32, float64,
        complex64, complex128
}

func GeneralAbsDifference(type T in Numeric)(a, b T) T {
    // T is a concrete type, so *T can be converted to *Underlying(T).
    // The set of possible types for T is infinite (any defined integer type),
    // but the set of possible types for Underlying(T) is small
    // (only the built-in integer types).
    var result T
    var rp interface{} = (*Underlying(T))(&result)

    // Convert a and b to their underlying types so that we
    // can type-assert them to a concrete type from a finite list.
    var au, bu interface{} = underlying(a), underlying(b)

    // Now we can write an exhaustive type switch over the list of possible types of
    // rp, au, and bu.
    switch rp := rp.(type) {
    case *int:
        *rp = OrderedAbsDifference(au.(int), bu.(int))
    case *int8:
        *rp = OrderedAbsDifference(au.(int8), bu.(int8))
    …

    case *complex64:
        *rp = ComplexAbsDifference(au.(complex64), bu.(complex64))
    …

    default:
        // If, say, int128 is added to a future version of the language,
        // we will need to add a case for it.
        panic(fmt.Sprintf("%T is not a recognized numeric type", au))
    }

    return result
}
```

[type]: https://golang.org/ref/spec#Types
[implements]: https://golang.org/ref/spec#Interface_types
[concrete type]: https://golang.org/ref/spec#Variables
[concrete types]: https://golang.org/ref/spec#Variables
[conversions]: https://golang.org/ref/spec#Conversions
[sum types]: https://golang.org/issue/19412
[Type Parameters draft]: https://golang.org/design/go2draft-type-parameters
[type list]: https://golang.org/design/go2draft-type-parameters#type-lists-in-constraints
[pointer methods]: https://golang.org/design/go2draft-type-parameters#pointer-methods
[substitution lemma]: http://twelf.org/wiki/Substitution_lemma
[Featherweight Go]: https://arxiv.org/abs/2005.11710 "Featherweight Go, 2020"
