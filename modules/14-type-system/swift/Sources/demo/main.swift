import Ts

let _door: Door<Closed> = Door<Closed>().open().close()
print("Door<Closed>.open().close() typechecks (Swift typestate)")
