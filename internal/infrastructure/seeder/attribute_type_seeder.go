package seeder

import (
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/ports"
)

var AttributeNames = []string{
	"سقف دار",
	"سرپوشیده",
	"روباز",
	"چند طبقه",
	"دوربین مداربسته",
	"نگهبان",
	"حصار",
	"درب اتوماتیک",
	"نورپردازی مناسب",
	"کارواش",
	"تعمیر لاستیک",
	"شارژ خودرو برقی",
	"باتری زدن",
	"روغن موتور",
	"۲۴ ساعته",
	"پارکبان",
	"مناسب معلولین",
	"جایگاه موتورسیکلت",
	"مناسب ون",
	"رزرو آنلاین",
	"اجاره ماهانه",
	"اجاره ساعتی",
	"وی آی پی",
	"پارکینگ طولانی مدت",
}

type AttributeTypeSeeder struct {
	unitOfWork ports.UnitOfWork
}

func NewAttributeTypeSeeder(unitOfWork ports.UnitOfWork) *AttributeTypeSeeder {
	return &AttributeTypeSeeder{
		unitOfWork: unitOfWork,
	}
}

func (ats *AttributeTypeSeeder) SeedAttributes() {
	err := ats.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		parkingRepository := rf.ParkingRepository()

		for _, attributeName := range AttributeNames {
			attribute, err := parkingRepository.GetAttributeByName(attributeName)

			if err != nil {
				return err
			}

			if attribute != nil {
				continue
			}

			attribute = &entities.ParkingAttribute{
				Name: attributeName,
			}

			err = parkingRepository.CreateAttribute(attribute)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
