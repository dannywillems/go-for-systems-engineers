// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "errs",
    targets: [
        .target(name: "Errs"),
        .executableTarget(name: "demo", dependencies: ["Errs"]),
        .testTarget(name: "ErrsTests", dependencies: ["Errs"]),
    ]
)
