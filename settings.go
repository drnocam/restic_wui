package main

var settings Settings

/*
* snapshots repo infos.
 */
type RepositoryInfo struct {
	time     string   `json: "time"`
	paths    []string `json: "paths"`
	tags     []string `json: "tags"`
	hostname string   `json: "hostname"`
	username string   `json: "username"`
	short_id string   `json: "short_id"`
	id       string   `json: "id"`
}

/*
* id to update uniqe
 */
type SavedRepository struct {
	Path     string `json: "path"`
	Password string `json: "password"`
	name     string `json: "name"`
	args     string `json: "args"`
	size     int    `json: "size"`
	id       int    `json: "id"`
}

/*
* default runcommand restic
 */
type Settings struct {
	RunCommand   string            `json: "runcmd`
	Repositories []SavedRepository `json: "repos`
}

func GetSettings() Settings {

	settings = Settings{
		RunCommand: "restic",
		Repositories: []SavedRepository{
			{
				Path:     "/home/cam/public_html/backup",
				Password: "1",
				name:     "www",
				size:     0,
				id:       1,
			},
		},
	}
	return settings
}
