package release

import "testing"

func TestNextVersionMinor(t *testing.T) {
	version, err := NextVersion([]string{"v0.15.0", "v0.16.0", "ignored"}, KindMinor, "")

	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.17.0" {
		t.Fatalf("expected v0.17.0, got %s", version)
	}
}

func TestNextVersionPatch(t *testing.T) {
	version, err := NextVersion([]string{"v0.15.0", "v0.16.0", "v0.16.1-rc.0"}, KindPatch, "")

	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.16.1" {
		t.Fatalf("expected v0.16.1, got %s", version)
	}
}

func TestNextVersionRCStartsFromNextMinor(t *testing.T) {
	version, err := NextVersion([]string{"v0.15.0", "v0.16.0"}, KindRC, "")

	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.17.0-rc.0" {
		t.Fatalf("expected v0.17.0-rc.0, got %s", version)
	}
}

func TestNextVersionRCIncrementsExistingRC(t *testing.T) {
	version, err := NextVersion([]string{"v0.16.0", "v0.17.0-rc.0", "v0.17.0-rc.1"}, KindRC, "")

	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.17.0-rc.2" {
		t.Fatalf("expected v0.17.0-rc.2, got %s", version)
	}
}

func TestNextVersionCustomRequiresVPrefix(t *testing.T) {
	_, err := NextVersion([]string{"v0.16.0"}, KindCustom, "0.17.0")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestNextVersionCustomRejectsExistingTag(t *testing.T) {
	_, err := NextVersion([]string{"v0.16.0"}, KindCustom, "v0.16.0")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestNextVersionCustomAcceptsValidVersion(t *testing.T) {
	version, err := NextVersion([]string{"v0.16.0"}, KindCustom, "v0.18.0")

	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.18.0" {
		t.Fatalf("expected v0.18.0, got %s", version)
	}
}
