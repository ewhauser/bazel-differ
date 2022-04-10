load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file", "http_jar")

#--------
# Go
#--------
## rules_go
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains(version = "1.17.6")

gazelle_dependencies()

#------------
# Bazel tools
#------------
# buildifier
http_file(
    name = "buildifier",
    executable = True,
    sha256 = "3ed7358c7c6a1ca216dc566e9054fd0b97a1482cb0b7e61092be887d42615c5d",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/5.0.1/buildifier-linux-amd64"],
)

http_file(
    name = "buildifier_osx",
    executable = True,
    sha256 = "2cb0a54683633ef6de4e0491072e22e66ac9c6389051432b76200deeeeaf93fb",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/5.0.1/buildifier-darwin-amd64"],
)

http_file(
    name = "buildozer",
    executable = True,
    sha256 = "78204dac0ac6a94db499c57c5334b9c0c409d91de9779032c73ad42f2362e901",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/5.0.1/buildozer-linux-amd64"],
)

http_file(
    name = "buildozer_osx",
    executable = True,
    sha256 = "17a093596f141ead6ff70ac217a063d7aebc86174faa8ab43620392c17b8ee61",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/5.0.1/buildozer-darwin-amd64"],
)

http_file(
    name = "protoc_bin",
    executable = True,
    sha256 = "75d8a9d7a2c42566e46411750d589c51276242d8b6247a5724bac0f9283e05a8",
    urls = ["https://github.com/google/protobuf/releases/download/v3.20.0/protoc-3.20.0-linux-x86_64.zip"],
)

http_file(
    name = "protoc_bin_osx",
    executable = True,
    sha256 = "8b35a679c99b36caef5899e596281fe0b943ed248f7d5f70b3e705684bf67cb4",
    urls = ["https://github.com/google/protobuf/releases/download/v3.20.0/protoc-3.20.0-osx-x86_64.zip"],
)

#---------------
# Bazel diff
#---------------
http_jar(
    name = "bazel_diff",
    sha256 = "0dc9166097c181796cfcbde4d16ab304e70cc853c536257c3dfa9af60ddd31e6",
    urls = [
        "https://github.com/Tinder/bazel-diff/releases/download/3.3.0/bazel-diff_deploy.jar",
    ],
)
