import Foundation
import Shapes

let shapes: [any Shape] = [Circle(r: 1.0), Square(s: 2.0)]
print(String(format: "a Circle as any Shape has area %.4f", shapes[0].area()))
print(String(format: "sum via existential dispatch = %.4f", sumDynamic(shapes)))
