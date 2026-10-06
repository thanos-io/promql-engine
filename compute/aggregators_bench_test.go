package compute

import (
	"math"
	"testing"
)

func TestNewMaxAcc(t *testing.T) {
	t.Run("empty accumulator has NoValue", func(t *testing.T) {
		acc := NewMaxAcc()
		if got := acc.ValueType(); got != NoValue {
			t.Errorf("ValueType() = %v, want NoValue", got)
		}
		val, _ := acc.Value()
		if val != 0 {
			t.Errorf("Value() = %v, want 0", val)
		}
	})

	t.Run("single Add sets max", func(t *testing.T) {
		acc := NewMaxAcc()
		_ = acc.Add(42, nil)
		val, _ := acc.Value()
		if val != 42 {
			t.Errorf("Value() = %v, want 42", val)
		}
		if acc.ValueType() != SingleTypeValue {
			t.Errorf("ValueType() = %v, want SingleTypeValue", acc.ValueType())
		}
	})

	t.Run("multiple Add keeps max", func(t *testing.T) {
		acc := NewMaxAcc()
		for _, v := range []float64{3, 7, 1, 9, 2} {
			_ = acc.Add(v, nil)
		}
		val, _ := acc.Value()
		if val != 9 {
			t.Errorf("Value() = %v, want 9", val)
		}
	})

	t.Run("AddVector takes max of slice", func(t *testing.T) {
		acc := NewMaxAcc()
		vec := []float64{10, 5, 20, 15}
		_ = acc.AddVector(vec, nil)
		val, _ := acc.Value()
		if val != 20 {
			t.Errorf("Value() = %v, want 20", val)
		}
	})

	t.Run("Add and AddVector combined", func(t *testing.T) {
		acc := NewMaxAcc()
		_ = acc.Add(100, nil)
		_ = acc.AddVector([]float64{50, 80, 120}, nil)
		val, _ := acc.Value()
		if val != 120 {
			t.Errorf("Value() = %v, want 120", val)
		}
	})

	t.Run("AddVectorSIMD takes max of slice", func(t *testing.T) {
		acc := NewMaxAcc()
		vec := []float64{10, 5, 20, 15}
		_ = acc.AddVectorSIMD(vec, nil)
		val, _ := acc.Value()
		if val != 20 {
			t.Errorf("Value() = %v, want 20", val)
		}
	})

	t.Run("Add and AddVectorSIMD combined", func(t *testing.T) {
		acc := NewMaxAcc()
		_ = acc.Add(100, nil)
		_ = acc.AddVectorSIMD([]float64{50, 80, 120}, nil)
		val, _ := acc.Value()
		if val != 120 {
			t.Errorf("Value() = %v, want 120", val)
		}
	})

	t.Run("NaN is replaced by real number", func(t *testing.T) {
		acc := NewMaxAcc()
		_ = acc.Add(math.NaN(), nil)
		_ = acc.Add(5, nil)
		val, _ := acc.Value()
		if val != 5 {
			t.Errorf("Value() = %v, want 5", val)
		}
	})

	t.Run("Reset clears state", func(t *testing.T) {
		acc := NewMaxAcc()
		_ = acc.Add(99, nil)
		acc.Reset(0)
		if acc.ValueType() != NoValue {
			t.Errorf("after Reset, ValueType() = %v, want NoValue", acc.ValueType())
		}
		val, _ := acc.Value()
		if val != 0 {
			t.Errorf("after Reset, Value() = %v, want 0", val)
		}
	})
}

func BenchmarkMaxAcc_AddVector(b *testing.B) {
	acc := NewMaxAcc()
	vec := make([]float64, 100)
	for i := range vec {
		vec[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = acc.AddVector(vec, nil)
		_, _ = acc.Value()
	}
}

func BenchmarkMaxAcc_AddVectorSIMD(b *testing.B) {
	acc := NewMaxAcc()
	vec := make([]float64, 100)
	for i := range vec {
		vec[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = acc.AddVectorSIMD(vec, nil)
		_, _ = acc.Value()
	}
}
