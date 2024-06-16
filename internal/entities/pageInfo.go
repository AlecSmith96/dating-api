package entities

type PageInfo struct {
	MinAge           int      `json:"minAge"`
	MaxAge           int      `json:"maxAge"`
	PreferredGenders []string `json:"preferredGenders"`
}
