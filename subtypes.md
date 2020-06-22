# A Theory of Subtyping for Go

In Featherweight Go, the subtyping judgement (`Δ ⊢ τ <: σ`) plays a central
role. However, subtyping is notably absent from both the [Go specification][]
(which instead relies on [assignability][]) and the [Type Parameters draft][]
(which uses a combination of [interface][] implementation and [type list][]
matching). Here, I examine a general notion of subtyping for the full Go
language, _derivable from_ assignability.

## What is a type?

According to the Go specification, “[a] <dfn>[type][]</dfn> determines a set of
values together with operations and methods specific to those values”.

## What is a subtype?

Formal definitions of a “subtype” vary. [Fuh and Mishra] defined subtyping
“based on type embedding or coercion: type _t₁_ is a subtype of type _t₂_,
written <code>_t₁_ ▹ _t₂_</code>, if we have some way of mapping every value
with type _t₁_ to a value of type _t₂_.” That would seem to correspond to Go's
notion of [conversion][] operation.

However, [Reynolds][] defined subtypes in terms of only _implicit_ conversions:
“When there is an implicit conversion from sort ω to sort ω′, we write ω ≤ ω′
and say that ω is a _subsort_ (or _subtype_) of ω′. Syntactically, this means
that a phrase of sort ω can occur in any context which permits a phrase of sort
ω′.” This definition maps more closely to Go's assignability, and to the
“implements” notion of subtyping used in [Featherweight Go][], and it is the
interpretation I will use here.

Notably, all commonly-used definitions of subtyping are reflexive: a type τ is
always a subtype of itself.

## What is a subtype _in Go_?

Go allows implicit conversion of a value `x` to type `T` only when `x` is
“assignable to” `T`. Following Reynolds' definition, I interpret a “context
which permits a phrase of sort ω′” in Go as “an operation for which ω′ is
assignable to the required operand type”. That leads to the following
definition:

In Go, `T1` is a <dfn>subtype</dfn> of `T2` if, for every `T3` to which a value
of type `T2` is assignable, every value `x` of type `T1` is _also_ assignable to
`T3`. (Note that the precondition is “`x` _of type_ `T1`”, not “`x` _assignable
to_ `T1`”. The distinction is subtle, but important.)

Let's examine the cases in the Go specification.

> *   `x`'s type is identical to `T`

This gives the two reflexive rules from the Featherweight Go paper:

```
----------
Δ ⊢ α <: α
```

```
------------
Δ ⊢ τS <: τS
```

For the full Go language, we can generalize `τS` to include the built-in
non-interface [composite types][] (arrays, structs, pointers, functions, slices,
maps, and channels), boolean types, numeric types, and string types, which we
collectively call the <dfn>[concrete types][]</dfn>.

This part of the assignability definition also allows a (redundant) rule for
interface types, which is omitted from Featherweight Go:

```
------------
Δ ⊢ τI <: τI
```

> *   `x`'s type `V` and `T` have identical [underlying types][] and at least
>     one of `V` or `T` is not a [defined][] type.

This case makes each [literal][] type a subtype of every _[defined][]_ type with
the same underlying type:

```
literal(τ)   Δ ⊢ underlying(σ) = τ
-----------------------------------
            Δ ⊢ τ <: σ
```

Note that the converse does not hold: even though a defined type `σ` is
_assignable to_ the corresponding literal type `τ`, `σ` is not assignable to
_every_ type that permits a value of type `τ`. Consider
[this program](https://play.golang.org/p/lUVS33lRKGM):

```go
    type (
        T0 struct{}
        T1 = struct{}
        T2 struct{}
    )
    var (
        x0 T0
        x1 T1
        x2 T2
    )

    x1 = x0 // ok: T0 is assignable to struct{}
    x2 = x1 // ok: struct{} is assignable to T2

    x2 = x0 // error: cannot use x0 (type T0) as type T2 in assignment
```

Even though `T0` is assignable to `T1`, and `T1` is assignable to `T2`, a value
of type `T1` is permitted in a context where a value of type `T0` cannot occur —
specifically, an assigment to `T2`. Thus, `T0` cannot be a subtype of `T1`.

> *   `T` is an interface type and `x` [implements][] `T`.

This leads to the interface-subtyping rule `<:I` from Featherweight Go:

```
methodsΔ(τ) ⊇ methodsΔ(τI)
--------------------------
       Δ ⊢ τ <: τI
```

> *   `x` is a bidirectional channel value, `T` is a channel type, `x`'s type
>     `V` and `T` have identical element types, and at least one of `V` or `T`
>     is not a defined type.

This case is analogous to the case for “identical underlying types”. It makes
the literal type `chan T` a subtype of both `<-chan T` and `chan<- T`, and by
extension a subtype of any defined type with `<-chan T` or `chan<- T` as its
underlying type.

```
----------------------
Δ ⊢ chan τ <: <-chan τ
```

```
----------------------
Δ ⊢ chan τ <: chan<- τ
```

As with the case for identical underlying types, we cannot make a _defined_ type
with underlying type `chan T` a subtype of `<-chan T`, because it is
[possible](https://play.golang.org/p/R_fgYsEJw7S) to use a `<-chan T` in a
context that excludes other defined types.

> *   `x` is the predeclared identifier `nil` and `T` is a pointer, function,
>     slice, map, channel, or interface type.

The predeclared identifier `nil` does not itself have a type, so this case does
not contribute any subtypes. (If it did, then the `nil` type would be a subtype
of all pointer, function, slice, map, channel, and interface types.)

> *   `x` is an untyped [constant][] [representable][] by a value of type `T`.

As with `nil`, this case does not contribute any subtypes.

## Why “of type” instead of “assignable to”?

I could have chosen a different definition of subtyping, based on _assignability
to_ type `T1` instead of a value _of_ type `T1`. That would give a rule
something like: “`T1` is a subtype of `T2` if, for every value `x` _assignable
to_ `T1`, `x` is also assignable to `T2`.“ Why not use that definition instead?

Unfortunately, that alternate definition would mean that even interface types do
not have subtypes: if [defined][] type `T` has a [literal][] underlying type `U`
and implements inteface `I`, then a value of type `U` is assignable to `T`, but
not to `I`, so `U` could not be a subtype of `I`. A notion of subtyping that
does not even capture the simple interface-implementation matching of
Featherweight Go does not seem useful.

<!-- Go citations -->

[Go specification]: https://golang.org/ref/spec
[type]: https://golang.org/ref/spec#Types
[underlying types]: https://golang.org/ref/spec#Types
[composite types]: https://golang.org/ref/spec#Types
[literal]: https://golang.org/ref/spec#Types
[concrete types]: https://golang.org/ref/spec#Variables
[type switch]: https://golang.org/ref/spec#Type_switches
[defined]: https://golang.org/ref/spec#Type_definitions
[assignability]: https://golang.org/ref/spec#Assignability
[conversion]: https://golang.org/ref/spec#Conversions
[interface]: https://golang.org/ref/spec#Interface_types
[implements]: https://golang.org/ref/spec#Interface_types
[constants]: https://golang.org/ref/spec#Constants
[constant]: https://golang.org/ref/spec#Constants
[representable]: https://golang.org/ref/spec#Representability
[Type Parameters draft]: https://golang.org/design/go2draft-type-parameters
[type list]: https://golang.org/design/go2draft-type-parameters#type-lists-in-constraints

<!-- Academic citations -->

[Fuh and Mishra]: https://link.springer.com/content/pdf/10.1007/3-540-19027-9_7.pdf "Type Inference with Subtypes, 1988"
[Reynolds]: https://link.springer.com/content/pdf/10.1007%2F3-540-10250-7_24.pdf "Using Category Theory to Design Implicit Conversions and Generic Operators, 1980"
[Featherweight Go]: https://arxiv.org/abs/2005.11710
