// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "ts",
    targets: [
        .target(name: "Ts"),
        .executableTarget(name: "demo", dependencies: ["Ts"]),
        .testTarget(name: "TsTests", dependencies: ["Ts"]),
    ]
)
