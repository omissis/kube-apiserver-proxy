# 6. Dependency alert tool

Date: 2023-05-01

## Status

Accepted

## Context

We need a tool that helps us update the project dependencies automatically, keeping it as secure as possible.

## Decision

We are going to use [Renovate bot](https://github.com/renovatebot/renovate) as the dependency maintainer: it is a GitHub app that simplifies the update process, opening new Pull Requests when new dependencies' updates are available.

## Consequences

This will allow us to keep our codebase dependecies updated at a lower cost, and it possibly help avoiding security issues.
