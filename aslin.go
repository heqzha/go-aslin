package aslin

type(
	AsFunc func(*Context)

	AsLine []AsFunc

	AsEngine struct{
		Matrix []AsLine
	}
)
