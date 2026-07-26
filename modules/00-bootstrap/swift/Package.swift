// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "bootstrap",
    targets: [
        .target(name: "Bootstrap"),
        .executableTarget(name: "demo", dependencies: ["Bootstrap"]),
        .testTarget(name: "BootstrapTests", dependencies: ["Bootstrap"]),
    ]
)
