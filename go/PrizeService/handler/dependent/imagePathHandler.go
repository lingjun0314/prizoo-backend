package dependent

type ImagePathHandler struct {}

func (e *ImagePathHandler) IsPathEmpty(path string) bool {
	return path == ""
}