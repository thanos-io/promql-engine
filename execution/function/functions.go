// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"math"
	"time"

	"github.com/prometheus/prometheus/model/histogram"
)

var invalidSample = sample{T: -1, F: 0}

type sample struct {
	T int64
	F float64
	H *histogram.FloatHistogram
}

type functionArgs struct {
	Samples      []sample
	StepTime     int64
	ScalarPoints []float64
}

type functionCall func(f functionArgs) sample

var instantVectorFuncs = map[string]functionCall{
	"abs":   simpleFunc(math.Abs),
	"ceil":  simpleFunc(math.Ceil),
	"exp":   simpleFunc(math.Exp),
	"floor": simpleFunc(math.Floor),
	"sqrt":  simpleFunc(math.Sqrt),
	"ln":    simpleFunc(math.Log),
	"log2":  simpleFunc(math.Log2),
	"log10": simpleFunc(math.Log10),
	"sin":   simpleFunc(math.Sin),
	"cos":   simpleFunc(math.Cos),
	"tan":   simpleFunc(math.Tan),
	"asin":  simpleFunc(math.Asin),
	"acos":  simpleFunc(math.Acos),
	"atan":  simpleFunc(math.Atan),
	"sinh":  simpleFunc(math.Sinh),
	"cosh":  simpleFunc(math.Cosh),
	"tanh":  simpleFunc(math.Tanh),
	"asinh": simpleFunc(math.Asinh),
	"acosh": simpleFunc(math.Acosh),
	"atanh": simpleFunc(math.Atanh),
	"rad": simpleFunc(func(v float64) float64 {
		return v * math.Pi / 180
	}),
	"deg": simpleFunc(func(v float64) float64 {
		return v * 180 / math.Pi
	}),
	"sgn": simpleFunc(func(v float64) float64 {
		var sign float64
		if v > 0 {
			sign = 1
		} else if v < 0 {
			sign = -1
		}
		if math.IsNaN(v) {
			sign = math.NaN()
		}
		return sign
	}),
	"round": func(f functionArgs) sample {
		if len(f.Samples) != 1 || len(f.ScalarPoints) > 1 {
			return invalidSample
		}

		toNearest := 1.0
		if len(f.ScalarPoints) > 0 {
			toNearest = f.ScalarPoints[0]
		}
		toNearestInverse := 1.0 / toNearest
		return sample{
			T: f.StepTime,
			F: math.Floor(f.Samples[0].F*toNearestInverse+0.5) / toNearestInverse,
		}
	},
	"pi": func(f functionArgs) sample {
		return sample{
			T: f.StepTime,
			F: math.Pi,
		}
	},
	"label_join": func(f functionArgs) sample {
		// This is specifically handled by functionOperator Series()
		return sample{}
	},
	"label_replace": func(f functionArgs) sample {
		// This is specifically handled by functionOperator Series()
		return sample{}
	},
	"time": func(f functionArgs) sample {
		return sample{
			T: f.StepTime,
			F: float64(f.StepTime) / 1000,
		}
	},
	"vector": func(f functionArgs) sample {
		if len(f.Samples) == 0 {
			return invalidSample
		}
		return sample{
			T: f.StepTime,
			F: f.Samples[0].F,
		}
	},
	"clamp": func(f functionArgs) sample {
		if len(f.Samples) == 0 || len(f.ScalarPoints) < 2 {
			return invalidSample
		}

		v := f.Samples[0].F
		min := f.ScalarPoints[0]
		max := f.ScalarPoints[1]

		if max < min {
			return invalidSample
		}

		return sample{
			T: f.StepTime,
			F: math.Max(min, math.Min(max, v)),
		}
	},
	"clamp_min": func(f functionArgs) sample {
		if len(f.Samples) == 0 || len(f.ScalarPoints) == 0 {
			return invalidSample
		}

		v := f.Samples[0].F
		min := f.ScalarPoints[0]

		return sample{
			T: f.StepTime,
			F: math.Max(min, v),
		}
	},
	"clamp_max": func(f functionArgs) sample {
		if len(f.Samples) == 0 || len(f.ScalarPoints) == 0 {
			return invalidSample
		}

		v := f.Samples[0].F
		max := f.ScalarPoints[0]

		return sample{
			T: f.StepTime,
			F: math.Min(max, v),
		}
	},
	"histogram_sum": func(f functionArgs) sample {
		if len(f.Samples) == 0 || f.Samples[0].H == nil {
			return invalidSample
		}
		return sample{
			T: f.StepTime,
			F: f.Samples[0].H.Sum,
		}
	},
	"histogram_count": func(f functionArgs) sample {
		if len(f.Samples) == 0 || f.Samples[0].H == nil {
			return invalidSample
		}
		return sample{
			T: f.StepTime,
			F: f.Samples[0].H.Count,
		}
	},
	"histogram_fraction": func(f functionArgs) sample {
		if len(f.Samples) == 0 || f.Samples[0].H == nil {
			return invalidSample
		}
		return sample{
			T: f.StepTime,
			F: histogramFraction(f.ScalarPoints[0], f.ScalarPoints[1], f.Samples[0].H),
		}
	},
	"histogram_quantile": func(f functionArgs) sample {
		// This is handled specially by operator.
		return sample{}
	},
	"days_in_month": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(32 - time.Date(t.Year(), t.Month(), 32, 0, 0, 0, 0, time.UTC).Day())
		})
	},
	"day_of_month": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.Day())
		})
	},
	"day_of_week": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.Weekday())
		})
	},
	"day_of_year": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.YearDay())
		})
	},
	"hour": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.Hour())
		})
	},
	"minute": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.Minute())
		})
	},
	"month": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.Month())
		})
	},
	"year": func(f functionArgs) sample {
		return dateWrapper(f, func(t time.Time) float64 {
			return float64(t.Year())
		})
	},
}

func simpleFunc(f func(float64) float64) functionCall {
	return func(fa functionArgs) sample {
		if len(fa.Samples) == 0 {
			return invalidSample
		}
		return sample{
			T: fa.StepTime,
			F: f(fa.Samples[0].F),
		}
	}
}

// Common code for date related functions.
func dateWrapper(fa functionArgs, f func(time.Time) float64) sample {
	if len(fa.Samples) == 0 {
		return sample{
			F: f(time.Unix(fa.StepTime/1000, 0).UTC()),
		}
	}
	t := time.Unix(int64(fa.Samples[0].F), 0).UTC()
	return sample{
		F: f(t),
	}
}
