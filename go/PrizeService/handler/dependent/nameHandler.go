package dependent

type NameHandler struct{}

func (e *NameHandler) IsPrizeNameEmpty(name string) bool{
	return name == ""
}