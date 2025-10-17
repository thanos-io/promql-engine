// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package extwarnings

func Coalesce(a, b error) error {
	if a != nil {
		return a
	}
	return b
}
