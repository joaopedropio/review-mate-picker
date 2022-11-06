package domain

type PickingType string

var (
	PickingTypeStateless = PickingType("stateless")
	PickingTypeStateful  = PickingType("stateful")
)
