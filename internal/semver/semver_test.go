package semver

import (
	"errors"
	"testing"
)

func TestLatest(t *testing.T) {
	t.Parallel()

	tags := []string{
		"v1.2.0", "v1.10.0", "v1.9.9", "", "not-a-tag", "v2.0", // plain
		"api/v1.9.0", "api/v1.10.1", "api/v1.10.0", // api
		"admin/v3.0.0", "v.1.0", "api/v1.11.0-rc1", // noise
	}

	tests := []struct {
		name      string
		prefix    string
		want      string
		wantFound bool
	}{
		{name: "plain tags only, numeric order", prefix: "", want: "v1.10.0", wantFound: true},
		{name: "prefixed tags only", prefix: "api/", want: "v1.10.1", wantFound: true},
		{name: "another prefix", prefix: "admin/", want: "v3.0.0", wantFound: true},
		{name: "unknown prefix starts at zero", prefix: "web/", want: "v0.0.0", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := Latest(tags, tt.prefix)
			if got.String() != tt.want || found != tt.wantFound {
				t.Fatalf("Latest(%q) = %s, %v; want %s, %v", tt.prefix, got, found, tt.want, tt.wantFound)
			}
		})
	}
}

func TestNext(t *testing.T) {
	t.Parallel()

	base := Version{Major: 1, Minor: 4, Patch: 2}

	tests := []struct {
		name string
		bump Bump
		want string
	}{
		{name: "patch", bump: Bump{Patch: true}, want: "v1.4.3"},
		{name: "minor resets patch", bump: Bump{Minor: true}, want: "v1.5.0"},
		{name: "major resets minor and patch", bump: Bump{Major: true}, want: "v2.0.0"},
		{name: "minor and patch (as some repos pass)", bump: Bump{Minor: true, Patch: true}, want: "v1.5.1"},
		{name: "major and minor", bump: Bump{Major: true, Minor: true}, want: "v2.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Next(base, tt.bump)
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != tt.want {
				t.Fatalf("Next(%s, %+v) = %s; want %s", base, tt.bump, got, tt.want)
			}
		})
	}
}

func TestNextWithoutBump(t *testing.T) {
	t.Parallel()

	if _, err := Next(Version{}, Bump{}); !errors.Is(err, ErrNoBump) {
		t.Fatalf("want ErrNoBump, got %v", err)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	for _, s := range []string{"v1", "1.2.3", "v1.2.3-rc1", "api/v1.2.3", "v.1.0"} {
		if _, ok := Parse(s); ok {
			t.Errorf("Parse(%q) should fail", s)
		}
	}

	if v, ok := Parse("v12.0.7"); !ok || v != (Version{12, 0, 7}) {
		t.Errorf("Parse(v12.0.7) = %v, %v", v, ok)
	}
}
