package dto

type UpdateCheckoutDTO struct {
	Status string `json:"status"`
}

type CreateCheckout struct {
	CheckoutURL string `json:"checkout_url"`
	Message     string `json:"message"`
}
