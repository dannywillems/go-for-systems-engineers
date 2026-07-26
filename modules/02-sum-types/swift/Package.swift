// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "expr",
    targets: [
        .target(name: "Calc"),
        .executableTarget(name: "demo", dependencies: ["Calc"]),
        .testTarget(name: "CalcTests", dependencies: ["Calc"]),
    ]
)
