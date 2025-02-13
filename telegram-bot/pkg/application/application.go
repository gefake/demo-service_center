package application

type ApplicationForCall struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Date        int    `json:"date"`
}
