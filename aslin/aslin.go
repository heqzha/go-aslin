package aslin

type(
	AsFunc func(*Context)

	AsLine []AsFunc

	AsEngine struct{
		Matrix []AsLine
	}
)

func (l AsLine) Last() AsFunc {
	length := len(l)
	if length > 0 {
		return l[length-1]
	}
	return nil
}
