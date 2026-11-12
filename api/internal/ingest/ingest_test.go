package ingest

import "testing"

func TestTitleFrom(t *testing.T) {
	cases := []struct{ name, prompt, want string }{
		{"a-cute-classic-school", "a cute classic school backpack in soft suede, with leather straps", "Cute Classic School Backpack in Soft Suede"},
		{"curved-katana-sword-with", "curved katana sword with glowing neon edge", "Curved Katana Sword with Glowing Neon Edge"},
		{"Golden Crown", "a golden crown", "Golden Crown"}, // renamed manually in ugc_studio
		{"", "", "job-1"},
		{"mechanical-cyber-wings-with", "mechanical cyber wings with neon circuits", "Mechanical Cyber Wings with Neon Circuits"},
	}
	for _, c := range cases {
		if got := titleFrom(c.name, c.prompt, "job-1"); got != c.want {
			t.Errorf("titleFrom(%q, %q) = %q, want %q", c.name, c.prompt, got, c.want)
		}
	}
}

func TestSlugify(t *testing.T) {
	if got := slugify("Cute Classic — School Backpack!"); got != "cute-classic-school-backpack" {
		t.Errorf("slugify = %q", got)
	}
}
