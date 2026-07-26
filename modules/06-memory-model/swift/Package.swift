// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "conc",
    platforms: [.macOS(.v13)],
    targets: [
        .target(name: "Conc"),
        .executableTarget(name: "demo", dependencies: ["Conc"]),
        .testTarget(name: "ConcTests", dependencies: ["Conc"]),
    ]
)
