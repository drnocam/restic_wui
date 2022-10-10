package main

var settings Settings

type Repository struct {
	Path     string `json: "path"`
	Password string `json: "password"`
	index    int    `json: "index"`
}

/*
* default runcommand restic
 */
type Settings struct {
	RunCommand   string       `json: "runcmd`
	Repositories []Repository `json: "repos`
}

func GetSettings() Settings {

	settings = Settings{
		RunCommand: "restic",
		Repositories: []Repository{
			{
				Path:     "/home/cam/public_html/backup",
				Password: "1",
				index:    1,
			},
		},
	}
	return settings
}
