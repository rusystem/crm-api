package domain

// Country представляет информацию о стране.
// @Description Country represents a country entity.
// @name Country
type Country struct {
	CountryId string `json:"country_id" example:"RU"` // Уникальный идентификатор страны.
	Name      string `json:"name" example:"Россия"`   // Название страны.
}

// Region представляет информацию о регионе в составе страны.
// @Description Region represents a region entity.
// @name Region
type Region struct {
	RegionId string `json:"region_id" example:"65"`           // Административный код первого уровня
	Name     string `json:"name" example:"Самарская Область"` // Название региона.
}

// City представляет информацию о городе.
// @Description City represents a city entity.
// @name City
type City struct {
	CityId int64  `json:"city_id" example:"499099"` // Уникальный идентификатор города.
	Name   string `json:"name" example:"Самара"`    // Название города.
}
