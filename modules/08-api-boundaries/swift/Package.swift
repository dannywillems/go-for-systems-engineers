// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "api",
    targets: [
        .target(name: "Account"),
        .testTarget(name: "AccountTests", dependencies: ["Account"]),
    ]
)
