// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "gen",
    targets: [
        .target(name: "Gen"),
        .executableTarget(name: "demo", dependencies: ["Gen"]),
        .testTarget(name: "GenTests", dependencies: ["Gen"]),
    ]
)
