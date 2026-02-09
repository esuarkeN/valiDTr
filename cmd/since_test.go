package cmd

import (
	"testing"
	"time"
)

func TestParseSince_EpochAndRFC3339(t *testing.T) {
	got, err := parseSince("epoch")
	if err != nil {
		t.Fatalf("parseSince(epoch): %v", err)
	}
	want := time.Unix(0, 0).UTC()
	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want.Format(time.RFC3339), got.Format(time.RFC3339))
	}

	got, err = parseSince("2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("parseSince(rfc3339): %v", err)
	}
	want = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want.Format(time.RFC3339), got.Format(time.RFC3339))
	}
}

func TestParseSince_Invalid(t *testing.T) {
	if _, err := parseSince("nope"); err == nil {
		t.Fatalf("expected error for invalid timestamp")
	}
}
