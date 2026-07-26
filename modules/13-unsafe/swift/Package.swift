// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "un",
    targets: [
        .target(name: "Un"),
        .executableTarget(name: "demo", dependencies: ["Un"]),
        .testTarget(name: "UnTests", dependencies: ["Un"]),
    ]
)
