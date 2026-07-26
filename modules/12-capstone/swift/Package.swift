// swift-tools-version:6.0
import PackageDescription

let package = Package(
    name: "capstone",
    platforms: [.macOS(.v13)],
    targets: [
        .target(name: "Cache"),
        .executableTarget(name: "bench", dependencies: ["Cache"]),
        .testTarget(name: "CacheTests", dependencies: ["Cache"]),
    ]
)
