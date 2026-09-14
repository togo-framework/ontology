<!-- togo-header -->
<div align="center">
  <picture><source media="(prefers-color-scheme: dark)" srcset=".github/assets/togo-mark-dark.svg" /><img src=".github/assets/togo-mark.svg" alt="ToGO" height="64" /></picture>
  <h1>togo-framework/ontology</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1F8A99" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/ontology"><img src="https://pkg.go.dev/badge/github.com/togo-framework/ontology.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/ontology
```

<!-- /togo-header -->

> **Experimental.** A permission-aware, authority-ranked **context graph** that
> grounds agents — the open, Go-native take on Genie-style *onto-rank*.

Feed it entities (metric/term definitions from resources, queries, dashboards)
with references between them; `Resolve(term)` returns the definitions to trust,
ranked by **authority** (PageRank over references × certified/usage/freshness)
and filtered by what the caller may see.

```go
g, _ := ontology.FromKernel(k)
g.Add(ontology.Entity{ID: "c", Term: "revenue", Definition: "canonical", Certified: true})
best := g.Resolve("revenue", func(e *ontology.Entity) bool { return canSee(e) })
```
<!-- togo-sponsors -->
---

<div align="center">
  <h3>💎 Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><img src=".github/assets/id8media.svg" height="44" alt="ID8 Media" /></a>
    &nbsp;&nbsp;&nbsp;&nbsp;
    <a href="https://one-studio.co"><img src=".github/assets/one-studio.jpeg" height="44" alt="One Studio" /></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
