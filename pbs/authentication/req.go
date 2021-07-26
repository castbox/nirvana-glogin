package anti_authentication

func NewCheckRequest() *CheckRequest {
	return new(CheckRequest)
}

func NewCheckUnsafeRequest() *CheckRequest {
	return new(CheckRequest)
}

func NewQueryRequest() *QueryRequest {
	return new(QueryRequest)
}

func NewStateQueryRequest() *StateQueryRequest {
	return new(StateQueryRequest)
}
