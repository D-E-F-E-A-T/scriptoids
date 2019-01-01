load("@subpar//:subpar.bzl", "par_binary")
load("@requirements//:requirements.bzl", "requirement")

par_binary(
    name='dist',
    srcs=['main.py', 'scriptoid.py', 'output.py'],
    main='main.py',
    deps=[
        requirement('Click'),
        requirement('colorama'),
        requirement('toml'),
        requirement('semver')
    ],
    visibility=['//visibility:public']
)