package slice

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randSlice(size int) []int {
	rand.Seed(time.Now().Unix())
	a := rand.Perm(size)
	b := rand.Perm(size)
	halfSize := size / 2
	return append(a[:halfSize], b[halfSize:]...)
}

func BenchmarkDeduplicate(b *testing.B) {
	testData := randSlice(1000)
	for n := 0; n < b.N; n++ {
		Deduplicate(testData)
	}
}

func TestDeduplicate(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"deduplicates", args{[]int{1, 2, 2, 3}}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := Deduplicate(tt.args.s)
				assert.ElementsMatch(t, got, tt.want)
			},
		)
	}
}
