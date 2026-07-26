// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "testkit",
    targets: [
        .target(name: "Normalize"),
        .testTarget(name: "NormalizeTests", dependencies: ["Normalize"]),
    ]
)
