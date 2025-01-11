package domain

// Country представляет информацию о стране.
// @Description Country represents a country entity.
// @name Country
type Country struct {
	Id   string `json:"id" example:"RU"`       // Уникальный идентификатор страны.
	Name string `json:"name" example:"Россия"` // Название страны.
}

// Region представляет информацию о регионе в составе страны.
// @Description Region represents a region entity.
// @name Region
type Region struct {
	Id   string `json:"id" example:"65"`                  // Административный код первого уровня
	Name string `json:"name" example:"Самарская Область"` // Название региона.
}

// City представляет информацию о городе.
// @Description City represents a city entity.
// @name City
type City struct {
	Id   int64  `json:"id" example:"499099"`   // Уникальный идентификатор города.
	Name string `json:"name" example:"Самара"` // Название города.
}
