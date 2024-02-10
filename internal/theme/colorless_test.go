package theme

import "testing"

func TestSetClassic(t *testing.T) {
	SetClassic()
	checker := func(m Theme) {
		for k := range m {
			if m[k].Color != "" {
				t.Errorf("SetClassic() failed, got color %v, want %v", m[k].Color, "")
			}
			if m[k].Underline != false {
				t.Errorf("SetClassic() failed, got underline %v, want %v", m[k].Underline, false)
			}
			if m[k].Bold != false {
				t.Errorf("SetClassic() failed, got bold %v, want %v", m[k].Bold, false)
			}
			if m[k].Faint != false {
				t.Errorf("SetClassic() failed, got faint %v, want %v", m[k].Faint, false)
			}
			if m[k].Italics != false {
				t.Errorf("SetClassic() failed, got italics %v, want %v", m[k].Italics, false)
			}
			if m[k].Blink != false {
				t.Errorf("SetClassic() failed, got blink %v, want %v", m[k].Blink, false)
			}
		}
	}
	DefaultAll.Apply(checker)
}
