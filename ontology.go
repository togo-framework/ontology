// Package ontology is an EXPERIMENTAL togo capability: a permission-aware,
// authority-ranked context graph that grounds agents (the open, Go-native take
// on Genie-style "onto rank"). A togo app feeds it entities — metric/term
// definitions extracted from resources, queries, dashboards — with references
// between them; Resolve(term) returns the definitions to trust, ranked by
// authority (PageRank over references × source signals) and filtered by what the
// caller may see.
//
// Status: experimental — the API may change.
package ontology

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/togo-framework/togo"
)

// Entity is one definition/term in the graph.
type Entity struct {
	ID         string
	Term       string
	Definition string
	Source     string    // where it came from (table, dashboard, doc…)
	Author     string
	Certified  bool      // blessed by an owner
	Uses       int       // how often teams rely on it
	Updated    time.Time // freshness
	Refs       []string  // IDs this entity derives from / cites
}

// Ranked is an entity with its computed authority score.
type Ranked struct {
	Entity Entity
	Score  float64
}

// Graph holds the entities + their reference edges.
type Graph struct {
	nodes map[string]*Entity
}

// New returns an empty graph.
func New() *Graph { return &Graph{nodes: map[string]*Entity{}} }

// Add inserts or replaces an entity.
func (g *Graph) Add(e Entity) {
	cp := e
	g.nodes[e.ID] = &cp
}

const Name = "ontology"

func init() {
	togo.RegisterProviderFunc(Name, togo.PriorityService, func(k *togo.Kernel) error {
		k.Set(Name, New())
		if k.Log != nil {
			k.Log.Info("plugin active", "plugin", Name)
		}
		return nil
	})
}

// FromKernel returns the ontology graph from the kernel container.
func FromKernel(k *togo.Kernel) (*Graph, bool) {
	v, ok := k.Get(Name)
	if !ok {
		return nil, false
	}
	g, ok := v.(*Graph)
	return g, ok
}

// pageRank runs damped PageRank over the reference edges (authority flows from an
// entity to the ones it derives from / cites).
func (g *Graph) pageRank(iters int, damping float64) map[string]float64 {
	n := len(g.nodes)
	rank := make(map[string]float64, n)
	if n == 0 {
		return rank
	}
	for id := range g.nodes {
		rank[id] = 1.0 / float64(n)
	}
	for i := 0; i < iters; i++ {
		next := make(map[string]float64, n)
		base := (1 - damping) / float64(n)
		for id := range g.nodes {
			next[id] = base
		}
		for id, e := range g.nodes {
			var outs []string
			for _, r := range e.Refs {
				if _, ok := g.nodes[r]; ok {
					outs = append(outs, r)
				}
			}
			if len(outs) == 0 { // dangling: spread evenly
				share := damping * rank[id] / float64(n)
				for k := range g.nodes {
					next[k] += share
				}
				continue
			}
			share := damping * rank[id] / float64(len(outs))
			for _, r := range outs {
				next[r] += share
			}
		}
		rank = next
	}
	return rank
}

// score combines PageRank authority with source signals: certified boost, usage,
// and freshness decay.
func score(pr float64, e *Entity, now time.Time) float64 {
	s := pr
	if e.Certified {
		s *= 1.5
	}
	s *= 1 + math.Log1p(float64(e.Uses))
	if !e.Updated.IsZero() {
		ageDays := now.Sub(e.Updated).Hours() / 24
		if ageDays > 0 {
			s *= 1 / (1 + ageDays/30) // ~halve per month old
		}
	}
	return s
}

// Resolve returns the entities matching term (case-insensitive), ranked by
// authority, keeping only those visible(e) returns true for (permission filter).
// visible may be nil (everything visible).
func (g *Graph) Resolve(term string, visible func(*Entity) bool) []Ranked {
	pr := g.pageRank(30, 0.85)
	now := time.Now()
	var out []Ranked
	tl := strings.ToLower(term)
	for id, e := range g.nodes {
		if !strings.Contains(strings.ToLower(e.Term), tl) {
			continue
		}
		if visible != nil && !visible(e) {
			continue
		}
		out = append(out, Ranked{Entity: *e, Score: score(pr[id], e, now)})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
