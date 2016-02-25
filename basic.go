package jsval

func (dv defaultValue) HasDefault() bool {
	return dv.initialized
}

func (dv defaultValue) DefaultValue() interface{} {
	return dv.value
}

func (nc nilConstraint) Validate(_ interface{}) error {
	return nil
}

func (nc nilConstraint) HasDefault() bool {
	return false
}

func (nc nilConstraint) DefaultValue() interface{} {
	return nil
}
