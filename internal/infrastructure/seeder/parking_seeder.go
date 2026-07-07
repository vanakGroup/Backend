package seeder

import (
	"fmt"
	"log"
	"time"

	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/enums"
	"gorm.io/gorm"
)

type ParkingSeeder struct {
	db *gorm.DB
}

func NewParkingSeeder(db *gorm.DB) *ParkingSeeder {
	return &ParkingSeeder{db: db}
}

func (ps *ParkingSeeder) SeedParkings() {
	var count int64
	ps.db.Model(&entities.Parking{}).Count(&count)
	if count >= 100 {
		log.Println("parking seeder: already seeded, skipping")
		return
	}

	var user entities.User
	if err := ps.db.First(&user).Error; err != nil {
		log.Println("parking seeder: no user found, skipping")
		return
	}

	var cities []entities.City
	if err := ps.db.Limit(10).Find(&cities).Error; err != nil || len(cities) == 0 {
		log.Println("parking seeder: no cities found, skipping")
		return
	}

	var attributes []entities.ParkingAttribute
	ps.db.Find(&attributes)

	titles := []string{
		"پارکینگ مرکزی", "پارکینگ شمالی", "پارکینگ جنوبی", "پارکینگ شرقی", "پارکینگ غربی",
		"پارکینگ آزاد", "پارکینگ سرپوشیده", "پارکینگ روباز", "پارکینگ طبقاتی", "پارکینگ هوشمند",
	}

	// Tehran bounding box lat: 35.56–35.83, lng: 51.15–51.65
	lats := []float64{35.70, 35.72, 35.74, 35.68, 35.66, 35.76, 35.78, 35.62, 35.64, 35.80}
	lngs := []float64{51.38, 51.40, 51.42, 51.36, 51.34, 51.44, 51.46, 51.30, 51.32, 51.48}
	prices := []float64{5000, 8000, 10000, 12000, 15000, 20000, 25000, 30000, 7500, 18000}

	today := time.Now().Truncate(24 * time.Hour)

	for i := 0; i < 100; i++ {
		city := cities[i%len(cities)]
		title := fmt.Sprintf("%s %d", titles[i%len(titles)], i+1)
		lat := lats[i%len(lats)] + float64(i%5)*0.001
		lng := lngs[i%len(lngs)] + float64(i%5)*0.001
		price := prices[i%len(prices)]

		parking := entities.Parking{
			Title:       title,
			Latitude:    lat,
			Longitude:   lng,
			BasePrice:   price,
			Status:      enums.ParkingStatusActive,
			DocumentURL: "seeder/placeholder.pdf",
			UserID:      user.ID,
		}

		if err := ps.db.Create(&parking).Error; err != nil {
			log.Printf("parking seeder: failed to create parking %d: %v", i+1, err)
			continue
		}

		address := entities.Address{
			CityID:       city.ID,
			PostalCode:   fmt.Sprintf("1%08d", i+1),
			UnitNumber:   uint(i%20 + 1),
			ExactAddress: fmt.Sprintf("خیابان آزادی، کوچه %d، پلاک %d", i%50+1, i%10+1),
			ParkingID:    parking.ID,
		}
		if err := ps.db.Create(&address).Error; err != nil {
			log.Printf("parking seeder: failed to create address for parking %d: %v", i+1, err)
		}

		if len(attributes) > 0 {
			attrCount := 2 + i%3
			selected := make([]*entities.ParkingAttribute, 0, attrCount)
			for j := 0; j < attrCount && j < len(attributes); j++ {
				attr := attributes[(i+j)%len(attributes)]
				selected = append(selected, &attr)
			}
			ps.db.Model(&parking).Association("Attributes").Append(selected)
		}

		var slots []entities.AvailabilitySlot
		for day := 0; day < 7; day++ {
			date := today.AddDate(0, 0, day)
			for hour := uint(8); hour < 20; hour++ {
				slots = append(slots, entities.AvailabilitySlot{
					Date:      date,
					Hour:      hour,
					IsBooked:  false,
					ParkingID: parking.ID,
				})
			}
		}
		if err := ps.db.Create(&slots).Error; err != nil {
			log.Printf("parking seeder: failed to create slots for parking %d: %v", i+1, err)
		}
	}

	log.Println("parking seeder: 100 parkings seeded successfully")
}
