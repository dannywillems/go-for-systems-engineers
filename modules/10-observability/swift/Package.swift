// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "observability",
    platforms: [.macOS(.v13)],
    targets: [
        .target(name: "Observability"),
        .executableTarget(name: "bench", dependencies: ["Observability"]),
        .testTarget(name: "ObservabilityTests", dependencies: ["Observability"]),
    ]
)
