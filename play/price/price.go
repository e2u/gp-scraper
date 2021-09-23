package price

type Price struct {
	Currency string  `json:"currency,omitempty"`
	Value    float64 `json:"value"`
}

type Type int

const (
	All Type = iota
	Free
	Paid
)
