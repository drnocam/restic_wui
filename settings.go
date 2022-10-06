package main

var settings Settings;

type Repository struct {
	Path string
	Password string
	index int
}

type Settings struct {
	Repositories []Repository
}

func GetSettings() Settings {

	settings = Settings{
		Repositories: []Repository{
			Repository{
				Path: "/home/cam/public_html/backup",
				Password : "1",
				index : 1,
			},
		},
	}
	return settings;
}