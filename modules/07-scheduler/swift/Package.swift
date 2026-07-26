// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "sched",
    platforms: [.macOS(.v13)],
    targets: [
        .target(name: "Sched"),
        .executableTarget(name: "demo", dependencies: ["Sched"]),
        .testTarget(name: "SchedTests", dependencies: ["Sched"]),
    ]
)
