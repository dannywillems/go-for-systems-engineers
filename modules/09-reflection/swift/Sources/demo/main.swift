// Shows the synthesized Equatable at work. Deterministic.
import Codegen

let a = Person(name: "Ada", age: 36)
print("describe:              \(a.describe())")
print("Equatable (synth):     \(a == Person(name: "Ada", age: 36))")
print("not equal after edit:  \(a == Person(name: "Ada", age: 37))")
