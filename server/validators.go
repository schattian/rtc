package server

func validateRm(body *reqBody) error {
	if body.Branch == "" {
		return errNoBranch
	}
	if body.Table == "" {
		return errNoTable
	}
	if body.Column == "" {
		return errNoColumn
	}
	return nil
}

func validateAdd(body *reqBody) error {
	if body.Branch == "" {
		return errNoBranch
	}
	if body.Table == "" {
		return errNoTable
	}
	return nil
}
