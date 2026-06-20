package database

import (
	"fmt"
	"sync"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func NewPostgresDataBase(dbConfig *bootstrap.Database) *gorm.DB {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
	)

	dbOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect database"))
		}

		dbInstance = db

		err = dbInstance.AutoMigrate(
			&entities.User{},
			&entities.Role{},
			&entities.Permission{},
			&entities.Parking{},
			&entities.Address{},
			&entities.Province{},
			&entities.City{},
			&entities.ParkingAttribute{},
			&entities.AvailabilitySlot{},
			&entities.ParkingImage{},
			&entities.Vehicle{},
			&entities.Wallet{},
			&entities.MockPayment{},
			&entities.Reservation{},
			&entities.Transaction{},
		)

		if err != nil {
			panic(err)
		}
	})

	return dbInstance
}
