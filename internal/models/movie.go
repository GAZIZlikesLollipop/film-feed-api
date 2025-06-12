package models

import "time"

type Genre struct {
	Id   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null"`
}

type Movie struct {
	Id           int64         `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name" gorm:"not null"`
	PosterURL    string        `json:"posterURL"`
	MovieURL     string        `json:"movieURL"`
	TrailerURL   string        `json:"trailerURL"`
	Duration     int           `json:"duration" gorm:"default:0"`
	Age          int           `json:"age" gorm:"default:0"`
	Genres       []Genre       `json:"genres" gorm:"many2many:movie_genres;"`
	Rating       float32       `json:"rating" gorm:"default:0"`
	Reviews      int64         `json:"reviews" gorm:"default:0"`
	Description  string        `json:"description"`
	Country      string        `json:"country"`
	Year         int           `json:"year" gorm:"default:0"`
	Budget       int64         `json:"budget" gorm:"default:0"`
	BoxOffice    int64         `json:"boxOffice" gorm:"default:0"`
	MovieMembers []MovieMember `json:"movieMembers" gorm:"foreignKey:MovieID"`
	Members      []Member      `json:"members" gorm:"many2many:movie_members;"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
