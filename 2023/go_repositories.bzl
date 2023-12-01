load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:Gvh4YaCaXNs6dKTlfgismwWZKyjVZXwOPfIyUaqU3No=",
        version = "v0.0.0-20231127185646-65229373498e",
    )
