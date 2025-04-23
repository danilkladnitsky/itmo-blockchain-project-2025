package authorize

type Model struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImdyQHF3ZXIud2VxIiwiZXhwIjoxNzM3MjE1MDE4LCJ1aWQiOiJlMDgxMzllNi1mYTQyLTQ4ZGMtYTRmYy0yOTNlNDY1N2NmNjAifQ.P81eWq85CbVAunISVzxHHDsXvIgNvwPJOCsidmdNUQA"`
}

type ResponseModel struct {
	Model Model `json:"auth"`
}

func (m Model) GetResponseData() interface{} {
	return Model{
		Token: m.Token,
	}
}
