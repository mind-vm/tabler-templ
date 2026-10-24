package theme

import "testing"

func TestAttrs(t *testing.T) {
	if a := (Settings{}).Attrs(); len(a) != 0 {
		t.Errorf("zero settings: %v", a)
	}
	if a := (Settings{Primary: Blue, Base: Neutral, Font: SansSerif, Radius: Radius1}).Attrs(); len(a) != 0 {
		t.Errorf("defaults are left out: %v", a)
	}
	a := Settings{Primary: Teal, Base: Zinc, Font: Serif, Radius: Radius0}.Attrs()
	for k, v := range map[string]string{"data-bs-theme-primary": "teal", "data-bs-theme-base": "zinc", "data-bs-theme-font": "serif", "data-bs-theme-radius": "0"} {
		if a[k] != v {
			t.Errorf("%s = %v, want %s", k, a[k], v)
		}
	}
}
