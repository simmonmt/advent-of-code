load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:4iLhBPcpqFmylhnkbY3W0ONLUYYkDAW9xMFLfxgsvCw=",
        version = "v0.0.0-20221208152030-732eee02a75a",
    )
