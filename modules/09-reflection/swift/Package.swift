// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "reflectgen",
    targets: [
        .target(name: "Codegen"),
        .executableTarget(name: "demo", dependencies: ["Codegen"]),
        .testTarget(name: "CodegenTests", dependencies: ["Codegen"]),
    ]
)
