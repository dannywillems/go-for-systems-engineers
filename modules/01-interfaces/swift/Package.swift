// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "shapes",
    targets: [
        .target(name: "Shapes"),
        .executableTarget(name: "demo", dependencies: ["Shapes"]),
        .testTarget(name: "ShapesTests", dependencies: ["Shapes"]),
    ]
)
