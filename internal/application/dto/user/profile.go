package userdto

type UserFilled struct {
	UserID       uint
	FirstName    *string
	LastName     *string
	NationalCode *string
	ShabaNumber  *string
}

type ProfileDTO struct {
	Phone        string `json:"phone"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NationalCode string `json:"national_code"`
	ShabaNumber  string `json:"shaba_number"`
}
