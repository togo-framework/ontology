package ontology

import "testing"

func TestReferencedEntityRanksHigher(t *testing.T) {
	g := New()
	// c is referenced (derived-from) by a and b → higher authority
	g.Add(Entity{ID: "a", Term: "revenue", Definition: "a", Refs: []string{"c"}})
	g.Add(Entity{ID: "b", Term: "revenue", Definition: "b", Refs: []string{"c"}})
	g.Add(Entity{ID: "c", Term: "revenue", Definition: "canonical"})
	res := g.Resolve("revenue", nil)
	if len(res) != 3 {
		t.Fatalf("want 3, got %d", len(res))
	}
	if res[0].Entity.ID != "c" {
		t.Fatalf("expected canonical 'c' ranked first, got %q (%v)", res[0].Entity.ID, res)
	}
}

func TestPermissionFilter(t *testing.T) {
	g := New()
	g.Add(Entity{ID: "pub", Term: "x", Source: "public"})
	g.Add(Entity{ID: "sec", Term: "x", Source: "secret"})
	res := g.Resolve("x", func(e *Entity) bool { return e.Source != "secret" })
	if len(res) != 1 || res[0].Entity.ID != "pub" {
		t.Fatalf("permission filter failed: %v", res)
	}
}

func TestCertifiedAndUsesBoost(t *testing.T) {
	g := New()
	g.Add(Entity{ID: "plain", Term: "m"})
	g.Add(Entity{ID: "blessed", Term: "m", Certified: true, Uses: 100})
	res := g.Resolve("m", nil)
	if res[0].Entity.ID != "blessed" {
		t.Fatalf("certified+used should rank first: %v", res)
	}
}
