package addressdto

type AddressRequest struct {
    CityID       uint   
    PostalCode   string 
    UnitNumber   uint   
    ExactAddress string 
}