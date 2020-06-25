# Coherent interface constraints for assignability and convertibility?

This document is a sketch of two interface types that could be added to the
[draft Type Parameters design](http://golang.org/design/go2draft-type-parameters)
to enable conversion and assignment in generic code.

--------------------------------------------------------------------------------

The interface type `convertible.To(T)` is implemented by any dynamic type that
can be converted to `T`, including `T` itself, any type assignable to `T`, and
any type that has the same _[underlying type][]_ as `T`. A variable of type
`convertible.To(T)` may be [converted][] to `T` itself. If the variable is the
`nil` interface value, the result of the conversion is the zero-value of `T`.

If `T` is itself an interface type, `convertible.To(T)` has the same method set
as `T`.

--------------------------------------------------------------------------------

The interface type `assignable.To(T)` is implemented by any dynamic type that is
assignable to `T`. A variable of the interface type `assignable.To(T)` is
assignable to `T` and, if the underlying type of `T` is a not a
[defined type][], also to that underlying type. If the variable is `nil`, the
value assigned is the zero-value of `T`.

If `T` is itself an interface type, `assignable.To(T)` has the same method set
as `T`.

(If `T` is _not_ an interface type, a variable of type `assignable.To(T)` must
not be assignable to any other defined type whose underlying type is `T` — even
if `T` itself is not a defined type — because `assignable.To(T)` may store
values of _other_ defined types.)

For example, given:

```go
type MyChan <-chan int
type OtherChan <-chan int

var (
    a chan int
    b <-chan int
    c assignable.To(MyChan)
    d MyChan
    e OtherChan
    f assignable.To(<-chan int)
    g chan<- int
)

```

-   `a` is assignable to `b`, `c`, `d`, `e`,`f`, and `g`.
-   `b` is assignable to `c`, `d`, `e`, and `f`.
-   `c` is assignable to `b`, `d`, and `f` (but not `e`, because `c` could
    contain a variable of type `MyChan` — which is not assignable to
    `OtherChan`).
-   `d` is assignable to `b`, `c`, and `f`.
-   `e` is assignable to `b` and `f`.
-   `f` is assignable to `b` (but not `d`, `d`, or `e`, because `f` could
    contain a variable of either type `OtherChan` or `MyChan`).
-   `g` is assignable to nothing.

[defined type]: https://golang.org/ref/spec#Type_definitions
[underlying type]: https://golang.org/ref/spec#Types
[implements]: https://golang.org/ref/spec#Interface_types
[converted]: https://golang.org/ref/spec#Conversions
[literal]: https://golang.org/ref/spec#Types
