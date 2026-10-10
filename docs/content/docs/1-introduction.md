---
title: "1. Introduction"
weight: 1
---

# Introduction & Architecture

While Go Fiber is incredibly fast and easy to use, large enterprise applications often struggle with state management,
dependency injection, and scattered route definitions across dozens of files.

**Rays Fiber** solves this by treating your web layer as standard IoC components.

## Core Benefits

1. **Injected Controllers:** Define your HTTP handlers as methods on a struct, allowing you to seamlessly inject repositories,
  services, and configurations using the Rays `@:"Inject"` tag.
2. **Centralized Lifecycle:** The HTTP server lifecycle (startup and graceful shutdown) is automatically managed by the Rays
  application context.
3. **Declarative Routing:** Define endpoint HTTP methods and paths directly via struct tags, keeping your routing logic exactly
  where the handler code lives.
