package repository

import (
	"testing"
)

// Teste de benchmark
func BenchmarkNewCampaignWithOnceInside(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewCampaign()
	}
}
