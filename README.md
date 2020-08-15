# Sample Go App

A simple HTTP service that doesn't do anything actually useful, but is ready to
go for Kubernetes and contains testing/monitoring.

Used for reference purposes for all the tooling/testing/boilerplate that goes
around a production-grade service, set up in a way that simplifies development.
The goal is to have as few manual steps or tribal knowledge gotchas as possible.

For context, this API allows a client to get information about some fictional
dogs as they roam around.

## Local Development Requirements

Specific tools (like Skaffold and Ginkgo) are handled by the [Makefile](./Makefile)
as much as possible, but you will need some general things to work on this.  We
want to minimize this as much as possible, and ideally this list should never
grow from here.

### Go

You'll need to [install Go](https://golang.org/doc/install).  This was built
using 1.14, but should work with newer versions.

### Docker

You'll need [Docker](https://www.docker.com).  Ideally get Docker for Desktop
so you can just enable Kubernetes and be done with it.

### Kubernetes

If you're using Docker for Desktop, you can just tick a box for Kubernetes.  If
you're on Linux or just using basic Docker, you'll have to set up your own stack
either locally or elsewhere.

### Make

You'll need GNU Make.  You probably already have it.  If you don't, google how
to install it for your OS.

## Running stuff

The [Makefile](./Makefile) has useful commands at the top with additional info.

To run things fully for the first time, you can use Skaffold.

```bash
# Runs our service in local k8s with all its dependencies and listens for any
# changes that we make locally.
make skaffold
```

## Tech choices

Uses [Go](https://golang.org) because it's great for microservices... and I love
Go.

Uses [Docker](https://www.docker.com) because Docker.  We'll use this both for
packaging up our own service and for using Dockerized tools like Swagger.

Uses [Kubernetes](https://kubernetes.io/) as the deployment target because it's
incredibly popular and great for services like this.

Uses [Ginkgo](https://github.com/onsi/ginkgo) for testing; generally I prefer
just using the base testing package, but Ginkgo is popular and I'm grudgingly
growing to like it anyway.  It also makes a great use case for vendoring a tool
and using it in the Makefile.

Uses [Skaffold](https://github.com/GoogleContainerTools/skaffold) for local
development to a k8s cluster.  This allows us to rapidly iterate on changes
but still test our service in its native k8s environment.

Uses good ol' [Makefiles](https://www.gnu.org/software/make/manual/make.html)
for all development dependencies and commands.  Make is ubiquitous and works
well with defining how to work with all our dependencies, tools, etc.

Uses [Swagger](https://swagger.io/) for API spec.  Open API Specification (OAS)
is a popular choice for defining an API spec, and Swagger has some great tools
for quickly generating docs and editing the spec as necessary.  A useful example
of a Dockerized tool, too!

Uses [MongoDB](https://www.mongodb.com/) to store some data.  Just an excuse
to play with some of its streaming features a little more, no real reason otherwise.
It also gives us an external dependency that our tooling will have to deal with.

Uses [Zap](https://github.com/uber-go/zap) for logging because it's fast and I
wanted to get some more hands on time with it to see if the API is as painful
as I thought it would be.

