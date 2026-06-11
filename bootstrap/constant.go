package bootstrap

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Constants struct {
	Context        Context
	Role           Role
	RedisKey       RedisKey
	ErrorTag       ErrorTag
	ErrorField     ErrorField
	JWTKeysPath    JWTKeysPath
	Server         Server
	ProfileFilleds ProfileFilled
	VehicleFields  VehicleField
	S3BucketPath   BucketPath
	ValidationTag  ValidationTag
}

type Server struct {
	Port string
}

// type BucketPath struct {
// }

type ValidationTag struct {
	IranianPhone  string
	IraninaPostal string
}

type ErrorTag struct {
	AlreadyRegistered      string
	MinimumLength          string
	ContainsLowercase      string
	ContainsUppercase      string
	ContainsNumber         string
	ContainsSpecialChar    string
	Expired                string
	Invalid                string
	NotRegistered          string
	NotVerified            string
	NotActive              string
	InvalidAuthCredentials string
	ExpiredAuthToken       string
	InvalidAuthToken       string
	InvalidSignedToken     string
	Unauthorized           string
	AwaitingApproval       string
	Rejected               string
	NotExist               string
	AlreadyExist           string
	ForbiddenStatus        string
	Pending                string
	AlreadyBlocked         string
	AlreadyActive          string
	AlreadyResolved        string
	AlreadyArchived        string
	AlreadyCompleted       string
	NotAccepted            string
	StatusNotChange        string
	AlreadyCanceled        string
	AlreadyRejected        string
	AlreadyAccepted        string
	AlreadyDraft           string
	InvalidRecaptcha       string
	Required               string
	InvalidNumber          string
	AlreadyUsed            string
	ClaimReject            string
	EmptyToken             string
	EmptyHeader            string
	InvalidHeader          string
	NationalCodeLength     string
	Numeric                string
	ShabaLength            string
	WrongFormat            string
	TooLarge               string
	Empty                  string
}

type ErrorField struct {
	Hashed               string
	Parking              string
	Attribute            string
	Slot                 string
	Recaptcha            string
	Reserve              string
	User                 string
	Phone                string
	Email                string
	Password             string
	OTP                  string
	Corporation          string
	NationalID           string
	RegistrationNumber   string
	IBAN                 string
	InstallationRequest  string
	Bid                  string
	Address              string
	Name                 string
	Province             string
	City                 string
	Page                 string
	ContactType          string
	Room                 string
	NotificationType     string
	Notification         string
	NotificationSetting  string
	Panel                string
	MaintenanceRequest   string
	MaintenanceRecord    string
	Ticket               string
	Role                 string
	Permission           string
	TicketComment        string
	Report               string
	ContactInformation   string
	PaymentTerm          string
	Guarantee            string
	GuaranteeViolation   string
	News                 string
	Media                string
	Blog                 string
	Post                 string
	Like                 string
	CorporationReview    string
	CorporationStaff     string
	Jwt                  string
	Auth                 string
	LicensePlate         string
	RegistrationCard     string
	RegistrationCardSize string
	CarModel             string
	VehicleId            string
	Vehicle              string
}

type ProfileFilled struct {
	Firstname    string
	Lastname     string
	NationalCode string
	ShabaNumber  string
}

type VehicleField struct {
	ID                  string
	LicensePlate        string
	CarModel            string
	RegistrationCardUrl string
}

type RedisKey struct {
	OtpKey string
}

type JWTKeysPath struct {
	PublicKey  string
	PrivateKey string
}

type Context struct {
	Translator string
	UserID     string
}

type Role struct {
	DefaultType string
}

type BucketPath struct {
}

func NewConstant() *Constants {
	return &Constants{
		ValidationTag: ValidationTag{
			IranianPhone:  "ir_phone",
			IraninaPostal: "ir_postal",
		},
		Server: Server{
			Port: "8080",
		},
		Context: Context{
			Translator: "translator",
			UserID:     "userID",
		},
		ProfileFilleds: ProfileFilled{
			Firstname:    "first_name",
			Lastname:     "last_name",
			NationalCode: "national_code",
			ShabaNumber:  "shaba_number",
		},
		VehicleFields: VehicleField{
			ID:                  "vehicle_id",
			LicensePlate:        "license_plate",
			CarModel:            "car_model",
			RegistrationCardUrl: "registration_card_url",
		},
		JWTKeysPath: JWTKeysPath{
			PublicKey:  "./internal/infrastructure/jwt/public.pem",
			PrivateKey: "./internal/infrastructure/jwt/private.pem",
		},
		Role: Role{
			DefaultType: "customer",
		},
		RedisKey: RedisKey{
			OtpKey: "otp",
		},
		ErrorField: ErrorField{
			Hashed:               "hashed",
			Parking:              "parking",
			Reserve:              "reserve",
			Slot:                 "slot",
			Recaptcha:            "recaptcha",
			User:                 "user",
			Phone:                "phone",
			Email:                "email",
			Password:             "password",
			OTP:                  "otp",
			Corporation:          "corporation",
			NationalID:           "nationalID",
			RegistrationNumber:   "registrationNumber",
			IBAN:                 "iban",
			InstallationRequest:  "installationRequest",
			Bid:                  "bid",
			Address:              "address",
			Name:                 "name",
			Province:             "province",
			City:                 "city",
			Page:                 "page",
			ContactType:          "contactType",
			Room:                 "room",
			NotificationType:     "notificationType",
			Notification:         "notification",
			Panel:                "panel",
			MaintenanceRequest:   "maintenanceRequest",
			MaintenanceRecord:    "maintenanceRecord",
			Ticket:               "ticket",
			Role:                 "role",
			Permission:           "permission",
			TicketComment:        "ticketComment",
			Report:               "report",
			ContactInformation:   "contactInformation",
			NotificationSetting:  "notificationSetting",
			PaymentTerm:          "paymentTerm",
			Guarantee:            "guarantee",
			GuaranteeViolation:   "guaranteeViolation",
			News:                 "news",
			Media:                "media",
			Blog:                 "blog",
			Post:                 "post",
			Like:                 "like",
			CorporationReview:    "corporationReview",
			CorporationStaff:     "corporationStaff",
			Jwt:                  "jwt",
			Auth:                 "auth",
			LicensePlate:         "licensePlate",
			RegistrationCard:     "registrationCard",
			RegistrationCardSize: "registrationCardSize",
			CarModel:             "carModel",
			VehicleId:            "vehiclId",
			Vehicle:              "vehicle",
			Attribute:            "attribute",
		},
		ErrorTag: ErrorTag{
			AlreadyRegistered:      "alreadyRegistered",
			MinimumLength:          "minimumLength",
			ContainsLowercase:      "containsLowercase",
			ContainsUppercase:      "containsUppercase",
			ContainsNumber:         "containsNumber",
			ContainsSpecialChar:    "containsSpecialChar",
			Expired:                "Expired",
			Invalid:                "invalid",
			NotRegistered:          "notRegistered",
			NotVerified:            "notVerified",
			NotActive:              "notActive",
			InvalidAuthCredentials: "invalidAuthCredentials",
			ExpiredAuthToken:       "expiredAuthToken",
			InvalidAuthToken:       "invalidAuthToken",
			InvalidSignedToken:     "invalidSignedToken",
			Unauthorized:           "unauthorized",
			AwaitingApproval:       "awaitingApproval",
			Rejected:               "rejected",
			NotExist:               "notExist",
			AlreadyExist:           "alreadyExist",
			ForbiddenStatus:        "forbiddenStatus",
			Pending:                "pending",
			AlreadyBlocked:         "alreadyBlocked",
			AlreadyActive:          "alreadyActive",
			AlreadyResolved:        "alreadyResolved",
			AlreadyArchived:        "alreadyArchived",
			AlreadyCompleted:       "alreadyCompleted",
			NotAccepted:            "notAccepted",
			StatusNotChange:        "statusNotChange",
			AlreadyCanceled:        "alreadyCanceled",
			AlreadyRejected:        "alreadyRejected",
			AlreadyAccepted:        "alreadyAccepted",
			AlreadyDraft:           "alreadyDraft",
			InvalidRecaptcha:       "invalidRecaptcha",
			Required:               "required",
			InvalidNumber:          "invalidNumber",
			AlreadyUsed:            "alreadyUsed",
			ClaimReject:            "claimReject",
			EmptyToken:             "emptyToken",
			EmptyHeader:            "emptyHeader",
			InvalidHeader:          "invalidHeader",
			NationalCodeLength:     "nationalCodeLength",
			Numeric:                "numeric",
			ShabaLength:            "shabaLength",
			WrongFormat:            "wrongFormat",
			TooLarge:               "tooLarge",
			Empty:                  "empty",
		},
	}
}

func (path *BucketPath) GenerateVehicleImagePath(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	uniqueName := uuid.New().String()
	return fmt.Sprintf("vehicles/%s%s", uniqueName, ext)
}

func (path *BucketPath) GetParkingDocPath(parkingID uint, fileName string) string {
	return fmt.Sprintf("parkings/%d/document_%d_%s", parkingID, time.Now().Unix(), fileName)
}

func (path *BucketPath) GetParkingImagePath(parkingID uint, index int, fileName string) string {
	return fmt.Sprintf("parkings/%d/images/img_%d_%d_%s", parkingID, index, time.Now().Unix(), fileName)
}
