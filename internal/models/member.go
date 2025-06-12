package models

import "time"

type MovieMember struct {
	MovieID   int64    `json:"movieId" gorm:"primaryKey;autoIncrement:false;index"`
	MemberID  int64    `json:"memberId" gorm:"primaryKey;autoIncrement:false;index"`
	Character string   `json:"character,omitempty" gorm:"default:NULL"`
	Roles     []string `json:"roles" gorm:"type:text;serializer:json"`
	Member    Member   `json:"member" gorm:"foreignKey:MemberID"`
}

type Member struct {
	Id            int64     `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"not null"`
	Photo         string    `json:"photo"`
	Roles         []string  `json:"roles" gorm:"type:text;serializer:json"`
	BirthDate     time.Time `json:"birthDate"`
	DeathDate     time.Time `json:"deathDate,omitempty" gorm:"default:NULL"`
	Biography     string    `json:"biography"`
	Nationality   string    `json:"nationality"`
	FeaturedFilms []Movie   `json:"featuredFilms" gorm:"many2many:movie_members;joinForeignKey:MemberID;joinReferences:MovieID;"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
