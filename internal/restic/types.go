package restic

// Snapshot represents a restic snapshot
type Snapshot struct {
	ID       string   `json:"id"`
	Time     string   `json:"time"`
	Tree     string   `json:"tree"`
	Paths    []string `json:"paths"`
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	Tags     []string `json:"tags"`
	ShortID  string   `json:"short_id"`
}

// RepositoryStats represents the output of 'restic stats'
type RepositoryStats struct {
	TotalSize      int64 `json:"total_size"`
	TotalFileCount int64 `json:"total_file_count"`
	TotalBlobCount int64 `json:"total_blob_count"`
}
