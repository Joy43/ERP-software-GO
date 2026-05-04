package errors

//-----------  Simplify converts any error to an AppError ------------------//
func Simplify(err error, record string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return Internal(err.Error())
}