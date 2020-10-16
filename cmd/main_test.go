package main

import (
	"testing"
)

func BenchmarkRegular(b *testing.B) {
	want := 2000000
	for i := 0; i < b.N; i++ {
		result := regular()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}	
}

func BenchmarkConcurrently(b *testing.B) {
	want := 2000000
	for i := 0; i < b.N; i++ {
		result := concurrently()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}
}