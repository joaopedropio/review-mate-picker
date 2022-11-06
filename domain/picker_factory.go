package domain

type PickerFactory interface {
	Build() Picker
}

type pickerFactory struct {
	pickingType PickingType
	persons     Persons
}

func NewPickerFactory(pickingType PickingType) PickerFactory {
	return &pickerFactory{
		pickingType: pickingType,
	}
}

func (f *pickerFactory) Build() Picker {
	switch f.pickingType {
	case PickingTypeStateful:
		return NewStatefulPicker()
	case PickingTypeStateless:
		return NewStatelessPicker()
	default:
		return NewStatelessPicker()
	}
}
