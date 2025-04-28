// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"sort"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestMergeMatcher(t *testing.T) {
	cases := []struct {
		name            string
		existingMatcher *labels.Matcher
		newMatcher      *labels.Matcher
		expected        *labels.Matcher
		shouldStop      bool
	}{
		{
			name:            "same matchers",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "existing matcher matches empty value",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", ""),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", ""),
			shouldStop:      false,
		},
		{
			name:            "new matcher matches empty value",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", ""),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", ""),
			shouldStop:      false,
		},
		{
			name:            "existing matcher matches nothing",
			existingMatcher: labels.MustNewMatcher(labels.MatchNotRegexp, "label", ".*"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        labels.MustNewMatcher(labels.MatchNotRegexp, "label", ".*"),
			shouldStop:      false,
		},
		{
			name:            "new matcher matches nothing",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchNotRegexp, "label", ".*"),
			expected:        labels.MustNewMatcher(labels.MatchNotRegexp, "label", ".*"),
			shouldStop:      false,
		},
		{
			name:            "existing matcher matches all values",
			existingMatcher: labels.MustNewMatcher(labels.MatchRegexp, "label", ".+"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "new matcher matches all values",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchRegexp, "label", ".+"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "existing matcher matches everything",
			existingMatcher: labels.MustNewMatcher(labels.MatchRegexp, "label", ".*"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "new matcher matches everything",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchRegexp, "label", ".*"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "both not equal matchers",
			existingMatcher: labels.MustNewMatcher(labels.MatchNotEqual, "label", "value1"),
			newMatcher:      labels.MustNewMatcher(labels.MatchNotEqual, "label", "value2"),
			expected:        labels.MustNewMatcher(labels.MatchNotRegexp, "label", "value1|value2"),
			shouldStop:      false,
		},
		{
			name:            "equal matcher with regexp that matches",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchRegexp, "label", "val.*"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "equal matcher with regexp that doesn't match",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchRegexp, "label", "foo.*"),
			expected:        nil,
			shouldStop:      true,
		},
		{
			name:            "equal matcher with not equal that matches",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			newMatcher:      labels.MustNewMatcher(labels.MatchNotEqual, "label", "value2"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			shouldStop:      false,
		},
		{
			name:            "equal matcher with not equal that doesn't match",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchNotEqual, "label", "value"),
			expected:        nil,
			shouldStop:      true,
		},
		{
			name:            "not equal with equal matcher that matches",
			existingMatcher: labels.MustNewMatcher(labels.MatchNotEqual, "label", "value2"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			shouldStop:      false,
		},
		{
			name:            "not equal with equal matcher that doesn't match",
			existingMatcher: labels.MustNewMatcher(labels.MatchNotEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        nil,
			shouldStop:      true,
		},
		{
			name:            "equal matcher with not regexp that matches",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			newMatcher:      labels.MustNewMatcher(labels.MatchNotRegexp, "label", "value2|value3"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			shouldStop:      false,
		},
		{
			name:            "equal matcher with not regexp that doesn't match",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			newMatcher:      labels.MustNewMatcher(labels.MatchNotRegexp, "label", "val.*"),
			expected:        nil,
			shouldStop:      true,
		},
		{
			name:            "not regexp with equal matcher that matches",
			existingMatcher: labels.MustNewMatcher(labels.MatchNotRegexp, "label", "foo.*|bar.*"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			shouldStop:      false,
		},
		{
			name:            "not regexp with equal matcher that doesn't match",
			existingMatcher: labels.MustNewMatcher(labels.MatchNotRegexp, "label", "val.*"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value"),
			expected:        nil,
			shouldStop:      true,
		},
		{
			name:            "incompatible matchers",
			existingMatcher: labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
			newMatcher:      labels.MustNewMatcher(labels.MatchEqual, "label", "value2"),
			expected:        nil,
			shouldStop:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, stop := mergeMatcher(tc.existingMatcher, tc.newMatcher)
			testutil.Equals(t, tc.shouldStop, stop)
			if tc.expected == nil {
				testutil.Equals(t, tc.expected, result)
			} else {
				testutil.Assert(t, matcherEqual(tc.expected, result))
			}
		})
	}
}

func TestPropagateMatchers(t *testing.T) {
	cases := []struct {
		name     string
		binOp    *Binary
		expected *Binary
	}{
		{
			name: "non vector selector LHS",
			binOp: &Binary{
				Op:  parser.ADD,
				LHS: &NumberLiteral{Val: 1},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
						},
					},
				},
			},
			expected: &Binary{
				Op:  parser.ADD,
				LHS: &NumberLiteral{Val: 1},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
						},
					},
				},
			},
		},
		{
			name: "non vector selector RHS",
			binOp: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
						},
					},
				},
				RHS: &NumberLiteral{Val: 1},
			},
			expected: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
						},
					},
				},
				RHS: &NumberLiteral{Val: 1},
			},
		},
		{
			name: "non equal metric name matcher",
			binOp: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchRegexp, labels.MetricName, "metric.*"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
						},
					},
				},
			},
			expected: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchRegexp, labels.MetricName, "metric.*"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
						},
					},
				},
			},
		},
		{
			name: "empty metric names",
			binOp: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric1"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value2"),
						},
					},
				},
			},
			expected: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric1",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric1"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value2"),
						},
					},
				},
			},
		},
		{
			name: "same metric names - should skip",
			binOp: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value2"),
						},
					},
				},
			},
			expected: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value1"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric"),
							labels.MustNewMatcher(labels.MatchEqual, "label", "value2"),
						},
					},
				},
			},
		},
		{
			name: "vector matching on labels",
			binOp: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric1",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric1"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label2", "value2"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label3", "value3"),
						},
					},
				},
				VectorMatching: &parser.VectorMatching{
					On:             true,
					MatchingLabels: []string{"label1"},
				},
			},
			expected: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric1",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric1"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label2", "value2"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label3", "value3"),
						},
					},
				},
				VectorMatching: &parser.VectorMatching{
					On:             true,
					MatchingLabels: []string{"label1"},
				},
			},
		},
		{
			name: "vector matching ignoring labels",
			binOp: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric1",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric1"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label2", "value2"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label3", "value3"),
						},
					},
				},
				VectorMatching: &parser.VectorMatching{
					On:             false,
					MatchingLabels: []string{"label2", "label3"},
				},
			},
			expected: &Binary{
				Op: parser.ADD,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric1",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric1"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label2", "value2"),
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "metric2",
						LabelMatchers: []*labels.Matcher{
							labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "metric2"),
							labels.MustNewMatcher(labels.MatchEqual, "label1", "value1"),
							labels.MustNewMatcher(labels.MatchEqual, "label3", "value3"),
						},
					},
				},
				VectorMatching: &parser.VectorMatching{
					On:             false,
					MatchingLabels: []string{"label2", "label3"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			propagateMatchers(tc.binOp)

			// Compare LHS
			lhs1, ok1 := tc.binOp.LHS.(*VectorSelector)
			lhs2, ok2 := tc.expected.LHS.(*VectorSelector)
			testutil.Equals(t, ok1, ok2)
			if ok1 {
				testutil.Equals(t, lhs2.Name, lhs1.Name)
				testutil.Assert(t, matchersEqual(lhs2.LabelMatchers, lhs1.LabelMatchers))
			} else {
				testutil.Equals(t, tc.expected.LHS, tc.binOp.LHS)
			}

			// Compare RHS
			rhs1, ok1 := tc.binOp.RHS.(*VectorSelector)
			rhs2, ok2 := tc.expected.RHS.(*VectorSelector)
			testutil.Equals(t, ok1, ok2)
			if ok1 {
				testutil.Equals(t, rhs2.Name, rhs1.Name)
				testutil.Assert(t, matchersEqual(rhs2.LabelMatchers, rhs1.LabelMatchers))
			} else {
				testutil.Equals(t, tc.expected.RHS, tc.binOp.RHS)
			}

			// Compare VectorMatching
			if tc.expected.VectorMatching == nil {
				testutil.Equals(t, tc.expected.VectorMatching, tc.binOp.VectorMatching)
			} else {
				testutil.Equals(t, tc.expected.VectorMatching.On, tc.binOp.VectorMatching.On)
				testutil.Equals(t, tc.expected.VectorMatching.MatchingLabels, tc.binOp.VectorMatching.MatchingLabels)
			}
		})
	}
}

func matchersEqual(m1, m2 []*labels.Matcher) bool {
	if len(m1) != len(m2) {
		return false
	}
	sort.Slice(m1, func(i, j int) bool { return m1[i].Name < m1[j].Name })
	sort.Slice(m2, func(i, j int) bool { return m2[i].Name < m2[j].Name })
	for i := range m1 {
		if !matcherEqual(m1[i], m2[i]) {
			return false
		}
	}
	return true
}
