package errors

func Handle(fn func() error, record string) error {
	err := fn()
	if err != nil {
		return Simplify(err, record)
	}
	return nil
}