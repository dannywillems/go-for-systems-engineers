# 08 — Comparison: five encapsulation systems

**Environment.** Go 1.26.5, Rust 1.92.0, OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10.
Every reject's exact compiler error is in the README (injected from the
`reject-*` directories).

## The unit and the ladder

| Language | Unit of encapsulation | Visibility levels | Default |
| -------- | --------------------- | ----------------- | ------- |
| Go       | the package + `internal/` subtree | 2: exported (Capitalized) / package-private (lowercase) | package-private |
| Rust     | the module tree | 4: `pub` / `pub(crate)` / `pub(super)` / private | private |
| OCaml    | the module, ascribed by its `.mli` | binary per item, plus ABSTRACT vs concrete types | everything in the `.ml` is public unless the `.mli` hides it |
| Swift    | the file and the module | 5: `private` / `fileprivate` / `internal` / `public` / `open` | `internal` (the module) |
| Kotlin   | the module | 4: `private` / `protected` / `internal` / `public` | `public` |

Two design axes separate them:

- **Granularity.** Go is the coarsest: two levels, chosen by capitalizing a
  letter. Swift is the finest ladder (five levels, distinguishing the file from
  the module and adding `open` for subclassability). Rust and Kotlin sit between.
- **Default.** Go, Rust, OCaml-via-`.mli`, and Swift default to the more private
  option (you opt INTO exposure); Kotlin defaults to `public` (you opt into
  hiding). A public-by-default language leaks more by omission.

## What "hidden" can mean

The rejects all prove access is denied, but they deny different things:

- **Go / Swift / Kotlin** hide by making a FIELD inaccessible. The type name is
  visible; the caller simply cannot read or write the field, nor use a struct/
  class literal to build one (Go: no exported fields to set; Swift/Kotlin: the
  initializer is the gate). The representation is present but unreachable.
- **Rust** hides per item: `Account` is `pub` but `balance` is private, and you
  can dial each field/function to `pub(crate)` or `pub(super)` independently, so
  the boundary radius is chosen field by field.
- **OCaml** hides the REPRESENTATION ITSELF. `type t` in the `.mli` is abstract:
  the client cannot even learn that `t` is a record. This is strictly more than
  field privacy — there is no representation to be denied access to, because none
  is exported. It is the basis of ML's abstract data types and functors.

The reject errors line up with these mechanisms:

| Language | What the reject does | Compiler error (abridged) |
| -------- | -------------------- | ------------------------- |
| Go       | import an `internal/` package from outside its subtree | use of internal package ... not allowed |
| Rust     | read a private field from another module | E0616: field ... is private |
| OCaml    | reach for an abstract type's record field | Unbound record field Store.n |
| Swift    | read a `private` property from outside | 'balance' is inaccessible due to 'private' protection level |
| Kotlin   | read a `private` field from outside the class | cannot access 'balance': it is private |

## Enforcement, not convention

The point the whole module makes: none of these is documentation or a naming
convention that a linter flags. Each is a hard compiler error — the program does
not produce a binary. Python's leading-underscore "privacy" or a `// do not use`
comment has no equivalent here; the type checker is the boundary. That is what
lets a library author change a hidden representation freely, knowing no consumer
could have depended on it, because the compiler forbade it.

## Bottom line

For hiding a field, all five are equivalent in effect (compile-time denial),
differing only in ergonomics and default. For hiding a whole representation
behind an abstract type name, OCaml's `.mli` is in a class of its own; Rust
approximates it with a private-field newtype, and Go with an unexported struct or
an interface, but the type's concreteness still leaks in a way `type t` does not.
Go trades expressiveness for a rule so simple it needs no keywords — you can read
a Go identifier's visibility from its first letter, which is the same
minimalism-over-power trade seen across the rest of these modules.
