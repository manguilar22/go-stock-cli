package stock

// Struct For JSON Responses

// Historical Data
type History struct {
	Chart struct {
		Result []struct {
			Meta       map[string]interface{} `json:"meta"`
			Timestamp  []int64                `json:"timestamp"`
			Indicators struct {
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
				Quote []struct {
					Open   []float64 `json:"open"`
					Close  []float64 `json:"close"`
					Low    []float64 `json:"low"`
					High   []float64 `json:"high"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
	Error struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"error"`
}
