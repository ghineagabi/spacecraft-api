package main

type LoginFromHeader struct {
	Auth string `header:"Authorization" binding:"required"`
}

type UserCredentials struct {
	Password string `json:"password" binding:"required,pw"`
	Email    string `json:"email" binding:"required"`
}

type FilteredSpacecraft struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type DetailedSpacecraft struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Class    string     `json:"class"`
	Crew     int        `json:"crew"`
	Image    string     `json:"image"`
	Value    float64    `json:"value"`
	Status   string     `json:"status"`
	Armament []Armament `json:"armament"`
}

type CreateSpacecraft struct {
	Name   string  `json:"name" binding:"required"`
	Class  string  `json:"class" binding:"required"`
	Crew   int     `json:"crew" binding:"required"`
	Image  string  `json:"image" binding:"required"`
	Value  float64 `json:"value" binding:"required"`
	Status string  `json:"status" binding:"required"`
}

type UpdateSpacecraft struct {
	Id     int    `json:"id" binding:"required"`
	Name   string `json:"name"`
	Class  string `json:"class"`
	Crew   int    `json:"crew"`
	Image  string `json:"image"`
	Value  int    `json:"value"`
	Status string `json:"status"`
}

type SpacecraftId struct {
	Id int `json:"id" binding:"required"`
}

type Armament struct {
	Title string `json:"title"`
	Qty   int    `json:"qty"`
}
