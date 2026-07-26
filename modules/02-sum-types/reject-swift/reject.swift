// DOES NOT TYPE-CHECK: a switch over an enum must be exhaustive. The missing
// `.blue` case is a compile error, the check Go omits.

enum Color {
  case red
  case green
  case blue
}

func name(_ c: Color) -> String {
  switch c {
  case .red: return "red"
  case .green: return "green"
  }
}
