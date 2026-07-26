// DOES NOT COMPILE. Equatable synthesis requires every stored property to be
// Equatable; NotEq is not, so the compiler cannot synthesize == and rejects it:
//
//   error: type 'Wrapper' does not conform to protocol 'Equatable'

struct NotEq {}

struct Wrapper: Equatable {
    let inner: NotEq
}
