package ec2fzf

import "testing"

func TestShowVersion(t *testing.T) {
	BuildTime = "BuildTime"
	Commit = "Commit"
	Release = "Release"
	Arch = "Arch"

	if BuildTime != "BuildTime" {
		t.Errorf("BuildTime is incorrect, got: %s, want: %s.", BuildTime, "BuildTime")
	}

	if Commit != "Commit" {
		t.Errorf("Commit is incorrect, got: %s, want: %s.", Commit, "Commit")
	}

	if Release != "Release" {
		t.Errorf("Release is incorrect, got: %s, want: %s.", Release, "Release")
	}

	if Arch != "Arch" {
		t.Errorf("Arch is incorrect, got: %s, want: %s.", Arch, "Arch")
	}
}
