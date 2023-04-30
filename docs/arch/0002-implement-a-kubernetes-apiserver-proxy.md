# 2. Implement a proxy for the kubernetes API server

Date: 2023-05-01

## Status

Accepted

## Context

Companies often needs to build user interfaces and dashboards to show the status of their Kubernetes clusters, and to allow their users to interact with them.
Such dashboards are built on top of a custom api layer that abstracts the kubernetes api in some way, and that talks to the kubernetes api server to perform the actual operations.
Given there are dozens -if not hundreds- of such dashboards, the amount of the work invested in writing such api layers is huge, and possibly redundant: the Kubernetes API Server -or a subset of it- is already exposed as a REST API, and it would be great if we could reuse it to some degree.

## Decision

We are going to build a configurable proxy for the Kubernetes API Server, that will allow the user to expose a subset of its endpoints as a REST, gRPC or GraphQL API: while the first two will likely be transparent proxies, the GraphQL one will be an entirely new implementation that will allow us to expose the Kubernetes API in a more frontend-friendly way, making the development of user interfaces and dashboards easier.

## Consequences

We hope that this project will allow us to reduce the amount of work needed to build new dashboards and user interfaces, discouraging the creation of new custom api layers and empowering the users to focus on the real business value of their projects.
