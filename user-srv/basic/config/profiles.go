package config

type Profiles interface {
	GetInclude() string
}

type defaultProfiles struct {
	Include string `json:"include"`
}

func (p deaultProfiles) GetInclude() string {
	return p.Include
}
