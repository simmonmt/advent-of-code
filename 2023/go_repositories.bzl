load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=",
        version = "v0.6.0",
    )

    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:Gvh4YaCaXNs6dKTlfgismwWZKyjVZXwOPfIyUaqU3No=",
        version = "v0.0.0-20231127185646-65229373498e",
    )

    go_repository(
        name = "com_github_aclements_go_z3",
        importpath = "github.com/aclements/go-z3",
        commit = "4675d5f90ca5778e64e8686b1f11401b9d16521e",
    )
