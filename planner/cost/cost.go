package cost

type Cost struct {
	CpuCost    float64
	MemoryCost float64
}

type CostModel interface {
	IsBetter(currentCost Cost, newCost Cost) bool
}
