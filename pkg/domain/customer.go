package domain

// Customer defines a customer
type Customer struct {
	Name string `json:"name"`
	City string `json:"city"`
}

// ValidateNewCustomer function to add validate a new customer
func ValidateNewCustomer(newCustomer Customer) error {
	return nil
}
