package logicalplan

import (
	"maps"
	"slices"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/thanos-io/promql-engine/query"
)

// setProjectionLabels is an optimizer that sets the projection labels for all vector selectors.
// If a projection is already set as a matcher, it will be materialized in the selector
// and the matcher will be removed.
// This is useful for sending projections as part of a remote query.
type setProjectionLabels struct{}

const seriesIDLabel = "__series__id"

func (s setProjectionLabels) Optimize(expr Node, _ *query.Options) (Node, annotations.Annotations) {
	var hasProjections bool
	TraverseBottomUp(nil, &expr, func(_ *Node, current *Node) bool {
		switch e := (*current).(type) {
		case *VectorSelector:
			// Check if a projection is already set in the node.
			hasProjections = e.Projection.Include || len(e.Projection.Labels) > 0
		}
		return hasProjections
	})
	if hasProjections {
		return expr, annotations.Annotations{}
	}

	var projection Projection
	return s.optimize(expr, projection)
}

func (s setProjectionLabels) optimize(expr Node, projection Projection) (Node, annotations.Annotations) {
	var stop bool
	Traverse(&expr, func(current *Node) {
		if stop {
			return
		}
		switch e := (*current).(type) {
		case *Aggregation:
			switch e.Op {
			case parser.TOPK, parser.BOTTOMK:
				projection = Projection{}
			default:
				projection = Projection{
					Labels:  slices.Clone(e.Grouping),
					Include: !e.Without,
				}
			}
			return
		case *FunctionCall:
			switch e.Func.Name {
			case "absent_over_time", "absent", "scalar":
				projection = Projection{Include: true}
			case "label_replace":
				switch projection.Include {
				case true:
					if slices.Contains(projection.Labels, UnsafeUnwrapString(e.Args[1])) {
						projection.Labels = append(projection.Labels, UnsafeUnwrapString(e.Args[3]))
					}
				case false:
					if !slices.Contains(projection.Labels, UnsafeUnwrapString(e.Args[1])) {
						projection.Labels = slices.DeleteFunc(projection.Labels, func(s string) bool {
							return s == UnsafeUnwrapString(e.Args[3])
						})
					}
				}
			}
		case *Binary:
			var highCard, lowCard = e.LHS, e.RHS
			if e.VectorMatching == nil || (!e.VectorMatching.On && len(e.VectorMatching.MatchingLabels) == 0) {
				if IsConstantExpr(lowCard) {
					s.optimize(highCard, projection)
				} else {
					s.optimize(highCard, Projection{})
				}

				if IsConstantExpr(highCard) {
					s.optimize(lowCard, projection)
				} else {
					s.optimize(lowCard, Projection{})
				}
				stop = true
				return
			}
			if e.VectorMatching.Card == parser.CardOneToMany {
				highCard, lowCard = lowCard, highCard
			}

			hcProjection := extendProjection(projection, e.VectorMatching.MatchingLabels)
			s.optimize(highCard, hcProjection)
			lcProjection := extendProjection(Projection{
				Include: e.VectorMatching.On,
				Labels:  append([]string{seriesIDLabel}, e.VectorMatching.MatchingLabels...),
			}, e.VectorMatching.Include)
			s.optimize(lowCard, lcProjection)
			stop = true
		case *VectorSelector:
			slices.Sort(projection.Labels)
			projection.Labels = slices.Compact(projection.Labels)
			e.Projection = projection
			projection = Projection{}
		}
	})

	return expr, annotations.Annotations{}
}

func extendProjection(projection Projection, lbls []string) Projection {
	var extendedLabels []string
	if projection.Include {
		extendedLabels = union(projection.Labels, lbls)
	} else {
		extendedLabels = intersect(projection.Labels, lbls)
	}
	return Projection{
		Include: projection.Include,
		Labels:  extendedLabels,
	}
}

// union returns the union of two string slices.
func union(l1 []string, l2 []string) []string {
	m := make(map[string]struct{})
	for _, s := range l1 {
		m[s] = struct{}{}
	}
	for _, s := range l2 {
		m[s] = struct{}{}
	}
	return slices.Collect(maps.Keys(m))
}

// intersect returns the intersection of two string slices.
func intersect(l1 []string, l2 []string) []string {
	m := make(map[string]struct{})
	var result []string
	for _, s := range l1 {
		m[s] = struct{}{}
	}
	for _, s := range l2 {
		if _, ok := m[s]; ok {
			result = append(result, s)
		}
	}
	return result
}
