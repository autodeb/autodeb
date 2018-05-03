package errorchecks

//IsInputError returns true if the error is an input error. The meaning of this
//is that the function has failed because of the input you provided and not
//some other external reason. For example, if you are building a REST API, this
//would be reason to return a 4XX status.
func IsInputError(err error) bool {

	type isInputErrorer interface {
		IsInputError() bool
	}

	switch err := cause(err).(type) {
	case isInputErrorer:
		return err.IsInputError()
	default:
		return false
	}

}
