git_repository(
    name = "io_bazel_rules_python",
    remote = "https://github.com/bazelbuild/rules_python.git",
    commit = "8b5d0683a7d878b28fffe464779c8a53659fc645",
)

git_repository(
    name = "subpar",
    remote = "https://github.com/google/subpar",
    tag = "1.3.0",
)

load("@io_bazel_rules_python//python:pip.bzl", 'pip_import')
pip_import(
    name = 'requirements',
    requirements = 'requirements.txt'
)

load("@requirements//:requirements.bzl", 'pip_install')
pip_install()
