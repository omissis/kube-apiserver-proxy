# 4. Setup main project libraries

Date: 2023-05-01

## Status

Accepted

## Context

As we are starting to develop the tool, we need to decide on a few architectural patterns to follow and some libraries to use.

## Decision

We are going to use [Cobra] as the main CLI library: it is a well-known library that is used by many other projects, it is easy to use and has a lot of features that we will need. We recognize there are lighter options available, but we think that the benefits of using a well-known library outweigh the cons of using a heavier one.

We will introduce a lightweight dependency injection container to make it easier to test the code and to relieve the code handling the business logic from being cluttered by initializations.

We will explore the use of [Slog] to handle logging, to try to stick to the standard library as much as possible and avoid vendor lock-in.

## Consequences

We will speed up development by using well-suppoerted, consolidated libraries, but we'll also leave some space for experimentation where we think it is needed, such as the logging library.

[cobra]: https://github.com/spf13/cobra
[slog]: https://pkg.go.dev/golang.org/x/exp/slog
