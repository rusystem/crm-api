package domain

type Country struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Region struct {
	AdminCode1 string `json:"adminCode1"`
	Name       string `json:"name"`
}

type City struct {
	Name string `json:"name"`
}
