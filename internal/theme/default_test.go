package theme

import "testing"

func TestAll_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		b          []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "nil",
			b:          nil,
			wantErr:    true,
			wantErrMsg: "unexpected end of JSON input",
		},
		{
			name:    "all",
			b:       []byte(`{"info": {"-": {"color": "white","bold": true,"icon": "","faint": true,"italics": true,"blink": true}},"permission": {"-": {"color": "BrightBlack"}},"size": {"-": {"color": "white"}},"user": {"owner": {"color": "yellow","bold": true}},"group": {"group": {"color": "yellow","bold": true}},"symlink": {"link-num": {"color": "red"}},"git": {"git-branch": {"color": "yellow"}},"name": {".azure": {"color": "white","icon": ""}},"special": {"char": {"color": "yellow","icon": ""}},"ext": {".profile": {"color": "BrightPreen","icon": ""}}}`),
			wantErr: false,
		},
		{
			name:       "failed key",
			b:          []byte(`{"info": {"-": {"color": "white","failed_key": true,"icon": "","faint": true,"italics": true,"blink": true}},"permission": {"-": {"color": "BrightBlack"}},"size": {"-": {"color": "white"}},"user": {"owner": {"color": "yellow","bold": true}},"group": {"group": {"color": "yellow","bold": true}},"symlink": {"link-num": {"color": "red"}},"git": {"git-branch": {"color": "yellow"}},"name": {".azure": {"color": "white","icon": ""}},"special": {"char": {"color": "yellow","icon": ""}},"ext": {".profile": {"color": "BrightPreen","icon": ""}}}`),
			wantErr:    true,
			wantErrMsg: "failed at key 'info': failed at key '-': unknown field: 'failed_key'",
		},
		{
			name:       "unknown field",
			b:          []byte(`{"unknown_field": {"-": {"color": "white","bold": true,"icon": "","faint": true,"italics": true,"blink": true}},"permission": {"-": {"color": "BrightBlack"}},"size": {"-": {"color": "white"}},"user": {"owner": {"color": "yellow","bold": true}},"group": {"group": {"color": "yellow","bold": true}},"symlink": {"link-num": {"color": "red"}},"git": {"git-branch": {"color": "yellow"}},"name": {".azure": {"color": "white","icon": ""}},"special": {"char": {"color": "yellow","icon": ""}},"ext": {".profile": {"color": "BrightPreen","icon": ""}}}`),
			wantErr:    true,
			wantErrMsg: "unknown field: 'unknown_field'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := All{}
			err := a.UnmarshalJSON(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("UnmarshalJSON() error = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

func TestAll_CheckLowerCase(t *testing.T) {
	DefaultAll.CheckLowerCase()
}

func TestAll_CheckLowerCase1(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should panic")
		} else {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	a := All{
		InfoTheme: Theme{"Info": {}},
	}
	a.CheckLowerCase()
}
