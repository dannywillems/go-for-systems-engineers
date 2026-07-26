// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "mem",
    targets: [
        .target(name: "Mem"),
        .executableTarget(name: "demo", dependencies: ["Mem"]),
        .testTarget(name: "MemTests", dependencies: ["Mem"]),
    ]
)
