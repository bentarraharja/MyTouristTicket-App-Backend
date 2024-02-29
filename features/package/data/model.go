package data

import (
	packages "my-tourist-ticket/features/package"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	TourID      uint
	PackageName string
	Price       int
	JumlahTiket int
	Benefits    []Benefit
}

type Benefit struct {
	gorm.Model
	PackageID uint
	Benefit   string
}

func CoreToModel(input packages.Core) Package {
	return Package{
		TourID:      input.TourID,
		PackageName: input.PackageName,
		Price:       input.Price,
		JumlahTiket: input.JumlahTiket,
	}
}

func (b Benefit) ModelToCoreBenefits() packages.BenefitCore {
	return packages.BenefitCore{
		ID:        b.ID,
		PackageID: b.PackageID,
		Benefit:   b.Benefit,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (p Package) ModelToCore() packages.Core {
	core := packages.Core{
		ID:          p.ID,
		TourID:      p.TourID,
		PackageName: p.PackageName,
		Price:       p.Price,
		JumlahTiket: p.JumlahTiket,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	for _, benefit := range p.Benefits {
		core.Benefits = append(core.Benefits, benefit.ModelToCoreBenefits())
	}
	return core
}
