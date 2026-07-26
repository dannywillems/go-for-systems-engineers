(* The implementation. The record type is concrete HERE, but because account.mli
   exposes `type t` abstractly, none of this representation escapes the module.
   Values are immutable, so deposit/withdraw return a new t. *)

type t = { balance : int }

exception Overdraft

let open_ initial = if initial < 0 then raise Overdraft else { balance = initial }
let deposit a amount = { balance = a.balance + amount }

let withdraw a amount =
  if amount > a.balance then raise Overdraft else { balance = a.balance - amount }

let balance a = a.balance
