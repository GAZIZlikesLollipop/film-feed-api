package models

import "time"

type MovieMember struct {
	Id        int64    `json:"id" gorm:"primaryKey"`
	MovieID   int64    `json:"movieId" gorm:"index"`
	MemberID  int64    `json:"memberId" gorm:"index"`
	Character string   `json:"character,omitempty" gorm:"default:NULL"`
	Roles     []string `json:"roles" gorm:"type:text;serializer:json"`
}

type Member struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Photo       string    `json:"photo"`
	Roles       []string  `json:"roles" gorm:"type:text;serializer:json"`
	BirthDate   time.Time `json:"birthDate"`
	DeathDate   time.Time `json:"deathDate,omitempty" gorm:"default:NULL"`
	Biography   string    `json:"biography"`
	Nationality string    `json:"nationality"`
	// FeaturedFilms []Movie   `json:"featuredFilms" gorm:"many2many:movie_members;joinForeignKey:MemberID;joinReferences:MovieID;"`
	FeaturedFilms []Movie `json:"featuredFilms" gorm:"many2many:movie_members;"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
