package main

var settings Settings

/*
* snapshots repo infos.
 */
type RepositoryInfo struct {
	Time     string   `json: "time"`
	Paths    []string `json: "paths"`
	Tags     []string `json: "tags"`
	Hostname string   `json: "hostname"`
	Username string   `json: "username"`
	Short_id string   `json: "short_id"`
	Id       string   `json: "id"`
}

/*
* id to update uniqe
 */
type SavedRepository struct {
	Path     string `json: "path"`
	Password string `json: "password"`
	Name     string `json: "name"`
	Args     string `json: "args"`
	Size     int    `json: "size"`
	Id       int    `json: "id"`
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
				Name:     "www",
				Size:     0,
				Id:       1,
			},
		},
	}
	return settings
}
func (s Settings) FindMaxId() int {
	var local_max int = 0
	for _, v := range s.Repositories {
		if v.Id > local_max {
			local_max = v.Id
		}
	}
	return local_max
}
func (s Settings) FindIndexById(id int) int {
	for index, v := range s.Repositories {
		if v.Id == id {
			return index
		}
	}
	return -1
}
func (s *Settings) AddRepository(repo SavedRepository) {
	/*
		find max id and add
	*/
	repo.Id = s.FindMaxId() + 1
	s.Repositories = append(s.Repositories, repo)
}
func (s *Settings) UpdateRepository(id int, repo SavedRepository) {
	/*
		find max id and add
	*/
	repo.Id = s.FindMaxId() + 1
	s.Repositories = append(s.Repositories, repo)
}
