package lightlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	l := NewLogger(100)
	l.Log(LevelAll, "Test Log")
}
func BenchmarkLog(b *testing.B) {
	b.N = 100
	l := NewLogger(100)
	for i := 0; i < b.N; i++ {
		l.Log(LevelAll, "test log")
	}
}
