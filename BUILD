load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/adamyi/easfs
gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = [
        "BookParser.go",
        "FooterParser.go",
        "MDParser.go",
        "ProjectParser.go",
        "YAMLParser.go",
        "error.go",
        "helper.go",
        "load.go",
        "main.go",
        "page.go",
        "redirector.go",
        "statusz.go",
    ],
    importpath = "github.com/adamyi/easfs",
    visibility = ["//visibility:private"],
    x_defs = {
        "github.com/adamyi/easfs.GitCommit": "{STABLE_GIT_COMMIT}",
        "github.com/adamyi/easfs.GitVersion": "{STABLE_GIT_VERSION}",
        "github.com/adamyi/easfs.Builder": "{STABLE_BUILDER}",
        "github.com/adamyi/easfs.BuildTimestamp": "{BUILD_TIMESTAMP}",
    },
    deps = [
        "@com_github_flosch_pongo2//:go_default_library",
        "@com_github_gholt_blackfridaytext//:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_gosimple_slug//:go_default_library",
        "@com_github_shirou_gopsutil//load:go_default_library",
        "@in_gopkg_russross_blackfriday_v2//:go_default_library",
        "@in_gopkg_yaml_v2//:go_default_library",
    ],
)

go_binary(
    name = "easfs",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
