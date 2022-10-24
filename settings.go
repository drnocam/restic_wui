package main

import (
	"encoding/json"
	"fmt"
	"restic_wui/fileop"
)

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

/*
* returns default settings if file not read
 */
func GetSettings() Settings {

	if v, ok := fileop.ReadSettings(); ok == nil {
		json.Unmarshal([]byte(v), &settings)
		if settings.RunCommand == "" {
			settings.RunCommand = "restic"
		}
		fmt.Println("getted settings:", settings)

	} else {
		fmt.Println(ok.Error())
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

/*
find max id and add
*/
func (s *Settings) AddRepoSettings(repo SavedRepository) {

	repo.Id = s.FindMaxId() + 1
	s.Repositories = append(s.Repositories, repo)
	if s, ok := json.Marshal(settings); ok == nil {
		fileop.WriteSettings(s)
	}
}

/*
find index and change value
*/
func (s *Settings) UpdateRepoSettings(id int, repo SavedRepository) {

	ind := s.FindIndexById(id)
	s.Repositories[ind] = repo

	if s, ok := json.Marshal(settings); ok == nil {
		fileop.WriteSettings(s)
	}

}

/*
find index and delete
*/
func (s *Settings) DeleteRepoFromSettings(id int) {

	ind := s.FindIndexById(id)
	if ind != -1 {
		s.Repositories = append(s.Repositories[:ind], s.Repositories[ind+1:]...)
		if s, ok := json.Marshal(settings); ok == nil {
			fileop.WriteSettings(s)
		}

	}

}

func (s *Settings) RepoNickNames() string {

	names := make(map[int]string)

	for _, v := range settings.Repositories {
		names[v.Id] = v.Name
	}
	nms, _ := json.Marshal(names)

	return string(nms)

}
