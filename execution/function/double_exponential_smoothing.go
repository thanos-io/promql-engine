package function

// Based on https://github.com/prometheus/prometheus/blob/8baad1a73e471bd3cf3175a1608199e27484f179/promql/functions.go#L438
// doubleExponentialSmoothing calculates the smoothed out value for the given series.
// It is similar to a weighted moving average, where historical data has exponentially less influence on the current data.
// It also accounts for trends in data. The smoothing factor (0 < sf < 1), aka "alpha", affects how historical data will affect the current data.
// A lower smoothing factor increases the influence of historical data.
// The trend factor (0 < tf < 1), aka "beta", affects how trends in historical data will affect the current data.
// A higher trend factor increases the influence of trends.
// Algorithm taken from https://en.wikipedia.org/wiki/Exponential_smoothing
func doubleExponentialSmoothing(vals []float64, sf, tf float64) (float64, bool) {
	// Check that the input parameters are valid
	if sf <= 0 || sf >= 1 || tf <= 0 || tf >= 1 {
		return 0, false
	}

	l := len(vals)
	// Can't do the smoothing operation with less than two points
	if l < 2 {
		return 0, false
	}

	var s0, s1, b float64
	// Set initial values
	s1 = vals[0]
	b = vals[1] - vals[0]

	// Run the smoothing operation
	var x, y float64
	for i := 1; i < l; i++ {
		// Scale the raw value against the smoothing factor
		x = sf * vals[i]
		// Scale the last smoothed value with the trend at this point
		b = calcTrendValue(i-1, tf, s0, s1, b)
		y = (1 - sf) * (s1 + b)
		s0, s1 = s1, x+y
	}

	return s1, true
}

// calcTrendValue calculates the trend value at the given index i.
// This is somewhat analogous to the slope of the trend at the given index.
// The argument "tf" is the trend factor.
// The argument "s0" is the previous smoothed value.
// The argument "s1" is the current smoothed value.
// The argument "b" is the previous trend value.
func calcTrendValue(i int, tf, s0, s1, b float64) float64 {
	if i == 0 {
		return b
	}
	x := tf * (s1 - s0)
	y := (1 - tf) * b
	return x + y
}
