package errors

type Problem struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	Status        int    `json:"status"`
	Detail        string `json:"detail"`
	CorrelationID string `json:"correlationId"`
}
