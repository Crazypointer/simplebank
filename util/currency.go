package util

const (
	// USD is the currency code of US dollar
	USD = "USD"
	// EUR is the currency code of Euro
	EUR = "EUR"
	// CAD is the currency code of Canadian dollar
	CAD = "CAD"

	// DefaultCurrency is the default currency code
	DefaultCurrency = USD
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
