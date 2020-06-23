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
meaning that “the argument must be a [subtype][] of the constraint type”.
However, if the meaning of an interface varies depending on whether it is a type
or a constraint the [substitution lemma][] would no longer hold: a type-list
interface passed as parameter `Τ` and interpreted as a proper interface type
would match _any_ [subtype][] of `T`, including a type-list interface naming any
subset of `T`'s type list, but upon substituting `T`'s actual type for the
constraint, it would no longer allow those interface subtypes.

In contrast, The built-in `comparable` constraint, if allowed as an ordinary
interface type, would have the same properties as other interface types: it is a
subtype of itself, and the `==` and `!=` operations are defined uniformly and
meaningfully for every pair of `comparable` values.

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

Type lists could be made more orthogonal in one of several ways.

We could acknowledge the difference between constraints and interface types and
move on:

1.  Define a separate `constraint` declaration for type constraints that are not
    interface types, and move the type-list feature from interfaces to
    constraints.

Or we could preserve as much of the current type-lists as we can as interface
types, and split out a constraint for the remaining behaviors:

1.  Define “type list” interfaces as underlying-type matchers, implemented by
    any type that has the same underlying type as an entry in the list.

2.  Add a built-in `concrete(T)` constraint (which is not an interface type),
    which constrains the type argument to be any [concrete type][] that is a
    [subtype][] of `T`.

Or we could break down type-lists into smaller orthogonal parts: disjunction
interfaces, underlying-type interfaces, and concrete-type constraints.

1.  Define “type list” interfaces as disjunctions of types, implemented by any
    type that _is_ or _implements_ at least one of the types in the list.

2.  Add a built-in `underlying(T)` interface, implemented by any type whose
    underlying type is or implements `T`.

3.  Add a built-in `concrete(T)` constraint (which is not an interface type),
    which constrains the type argument to be any [concrete type][] that is a
    [subtype][] of `T`.

Under the orthogonal option, the interface types would be defined as follows.

### Type-list interfaces

An interface type containing a type list is implemented by any type that _is_ or
_implements_ at least one of the types in the list. The method set of a
type-list interface is the intersection of the method sets of each type in the
list; no other methods may be defined.

The zero value of a type-list interface is the `nil` interface value, even if
none of the types in the list is itself an interface type.

A type switch or type assertion on a variable of a type-list interface may use
only the types in the list and interface types _implemented by_ at least one
type in the list. To allow lists to be expanded over time, a `default` case is
permitted even if the switch is exhaustive.

A type-list interface is assignable to any interface implemented by all of the
types in the list.

If all of the types in the list are convertible to some type `T`, then the
type-list interface is also convertible to `T`. If `T` is a concrete type, the
`nil` interface value converts to the zero-value of the type.

A type-list interface embedded in another interface restricts the other
interface to only the types in the list

### Underlying-type interfaces

The interface type `underlying(T)` is defined only for types `T` that are
predeclared boolean, numeric, or string types, type literals, or type-list
interfaces comprising the same. `underlying(T)` is itself a type-list interface,
but its list of types is unbounded.

A type `R` implements `underlying(T)` if either:

*   the underlying type of `R` is identical to `T` or occurs in the type-list of
    `T`,

*   or `R` is a pointer type that is not a defined type, and the type list of
    `T` includes a pointer type `P` with a pointer base type identical to the
    underlying type of the pointer base type of `R`.

(The second case above is intended to handle the special case of [conversions][]
for non-defined pointer types with compatible base types.)

A value of type `underlying(T)` can be converted to `T`. If `T` is an interface
type, the result of the conversion is an interface value whose dynamic type is
in the type list of `T`. The `nil` interface value converts to the zero-value of
`T`.

## Examples

With the orthogonal approach, the examples in the draft design become a bit more
verbose, but more explicit — and, notably, the GeneralAbsDifference function
becomes possible to implement.

```go
package constraints

type Ordered interface {
    underlying(interface {
        type int, int8, …, string
    })
}
```

```go
func Smallest(type T concrete(constraints.Ordered))(s []T) T {
    …
}
```

```go
// StringableSignedInteger is a type constraint that matches any
// type that is both 1) defined as a signed integer type;
// 2) has a String method.
type StringableSignedInteger interface {
    underlying(interface{ type int, int8, int16, int32, int64 })
    String() string
}
```

```go
// SliceConstraint is a type constraint that matches a slice of
// the type parameter.
type SliceConstraint(type T) interface {
    underlying([]T)
}
```

```go
type integer interface {
    underlying(interface{
        type int, int8, int16, int32, int64,
            uint, uint8, uint16, uint32, uint64, uintptr
    })
}

// Convert converts x, which may be any integer type
// (including an interface type whose allowable values are all integers),
// to the integer type To, which must be any specific (concrete) integer type.
func Convert(type To concrete(integer), From integer)(x From) To {
    to := To(x)
    if From(to) != x {
        panic("conversion out of range")
    }
    return to
}
```

```go
type builtinInteger interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr
}
type integer interface {
    underlying(builtinInteger)
}

func Add10(type T concrete(integer))(s []T) {
    for i, v := range s {
        s[i] = v + 10 // OK: 10 can convert to any concrete integer type
    }
}

// This function is INVALID.
func Add1024(type T concrete(integer))(s []T) {
    for i, v := range s {
        s[i] = v + 1024 // INVALID: 1024 not permitted by int8/uint8
    }
}

// This function is INVALID.
func Add10Interface(type T integer)(s []T) {
    for i, v := range s {
        s[i] = v + 10 // INVALID: operation + not permitted for type integer that may be an interface type
    }
}

```

```go

type BuiltinNumeric interface {
    type int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr,
        float32, float64,
        complex64, complex128
}

type Numeric interface {
    underlying(BuiltinNumeric)
}

type BuiltinNumericPointer interface {
    type *int, *int8, *int16, *int32, *int64
        *uint, *uint8, *uint16, *uint32, *uint64, *uintptr
        *float32, *float64
        *complex64, *complex128
}

type NumericPointer interface {
    underlying(BuiltinNumericPointer)
}

func GeneralAbsDifference(type T concrete(Numeric))(a, b T) T {
    var result T
    rp := BuiltinNumericPointer(NumericPointer(&result))  // Convert result to a pointer to its underlying type.

    // Now we can switch exhaustively on the possible underlying types,
    // and safely type-assert rp, a, and b to obtain suitable operands.
    switch a := BuiltinNumeric(a).(type) {
    case int:
        *(rp.(*int)) = OrderedAbsDifference(a, BuiltinNumeric(b).(int))
    case int8:
        *(rp.(*int8)) = OrderedAbsDifference(a, BuiltinNumeric(b).(int8))
    …

    case complex64:
        *(rp.(*complex64)) = ComplexAbsDifference(a, BuiltinNumeric(b).(complex64))
    …

    default:
        panic(fmt.Sprintf("%T is not a recognized numeric type", v))
    }

    return r
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
[subtype]: ./subtypes.md
