package cost

type Cost interface{}

type CostModel interface {
	IsBetter(currentCost Cost, newCost Cost) bool
}
