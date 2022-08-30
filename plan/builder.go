package plan

import "github.com/prometheus/alertmanager/pkg/labels"

type Builder interface {
	Aggregate(expr Operator, by bool, labels []string) Builder
	Select(matchers []labels.Matcher) Builder
}

type builder struct {
	plan plan
}

func NewBuilder() Builder {
	return builder{}
}

func (b builder) Aggregate(expr Operator, by bool, labels []string) Builder {
	return builder{
		plan: plan{
			rootOperator: expr,
		},
	}
}

func (b builder) Select(matchers []labels.Matcher) Builder {
	//TODO implement me
	panic("implement me")
}
