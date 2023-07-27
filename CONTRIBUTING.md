# Contributing

Thank you for your interest in contributing to `xk6-coap`! Please make sure to
familiarize yourself with the following guidelines to ensure that your changes
are able to be reviewed and merged in an efficient manner.

## Proposing a Change

While potentially uneccessary for minor fixes, contributors are encouraged to
[create an issue](https://github.com/golioth/xk6-coap/issues/new) prior to
working on any new features or significant updates. This allows for maintainers
to provide initial feedback or guidance that may inform how the change is
implemented.

If a bug is being reported, information about the version of `xk6-coap` used, as
well as steps to reproduce should be included in the issue.

For both new functionality and bug fix issues, contributors should indicate
whether they intend to implement the change. Otherwise, it will be assumed that
the issue is a request for another contributor to implement.

## Creating a Fork

The first step in modifying the `xk6-coap` is to
[fork](https://github.com/golioth/xk6-coap/fork) the repository under your own
GitHub account. All changes should be made on a branch and pushed to your fork
of the repository, before being proposed as an update to the `golioth/xk6-coap`
repository.

## Commit Style

`xk6-coap` uses the [conventional commits](https://www.conventionalcommits.org/)
standard for all changes. Every commit in the `xk6-coap` history should follow
this pattern.

## Opening a Pull Request

When changes are ready for review, a [pull
request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests)
should be opened from the branch in the `xk6-coap` repository forked under your
GitHub account to the `main` branch in the `golioth/xk6-coap` repository. If you
are looking for initial feedback, but additional work is needed prior to merge,
you may choose to open your pull request as a
[draft](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/changing-the-stage-of-a-pull-request).

`xk6-coap` maintainers actively monitor activity on the repository and should
respond promptly with a review. If you do not receive a review within a week of
opening the pull request, it is appropriate to tag a maintainer in a comment
requesting they provide feedback.

The body of the pull request should include enough context for a maintainer to
understand why the change is being proposed. When appropriate, the pull request
should include `Fixes #{issue-number}` so that the linked issue will be
automatically closed upon merge.

When the pull request includes new or changed behavior, a relevant example
should be added, or an existing one updated, in the `examples/` directory. Any
manual testing steps performed, as well as their output, should be included in
the pull request body.

## Merging a Pull Request

Maintainers are responsible for reviewing and merging pull requests. Maintainers
may opt to delay reviewing a pull request until all automated tests are passing,
which is required for changes to be merged. One approving review from a
maintainer is required for a pull request to be merged, and all approved pull
requests should be merged promptly.
