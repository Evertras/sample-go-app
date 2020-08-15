# Sample Go App

A super simple HTTP service that doesn't actually do anything useful on its own.
Used for reference purposes for all the tooling/testing/boilerplate that goes
around an actually useful service.

Uses [Ginkgo](https://github.com/onsi/ginkgo) for testing.

Uses [Skaffold](https://github.com/GoogleContainerTools/skaffold) for local
development to a k8s cluster.

Uses good ol' [Makefiles](https://www.gnu.org/software/make/manual/make.html)
for all development dependencies.

## Requirements

Specific tools (like Skaffold and Ginkgo) are handled by the [Makefile](./Makefile)
as much as possible, but you will need some general things to work on this.

### Go 1.14

Might even work on lower versions, just haven't tested.

### Docker

If you want to run Docker things, you'll need Docker.

