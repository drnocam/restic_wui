package restic

// SnapshotSummary contains statistics for a snapshot
type SnapshotSummary struct {
	FilesNew            uint64 `json:"files_new"`
	FilesChanged        uint64 `json:"files_changed"`
	FilesUnmodified     uint64 `json:"files_unmodified"`
	DirsNew             uint64 `json:"dirs_new"`
	DirsChanged         uint64 `json:"dirs_changed"`
	DirsUnmodified      uint64 `json:"dirs_unmodified"`
	DataBlobs           int64  `json:"data_blobs"`
	TreeBlobs           int64  `json:"tree_blobs"`
	DataAdded           uint64 `json:"data_added"`
	TotalFilesProcessed uint64 `json:"total_files_processed"`
	TotalBytesProcessed uint64 `json:"total_bytes_processed"`
}

// Snapshot represents a restic snapshot
type Snapshot struct {
	ID       string           `json:"id"`
	Time     string           `json:"time"`
	Tree     string           `json:"tree"`
	Paths    []string         `json:"paths"`
	Hostname string           `json:"hostname"`
	Username string           `json:"username"`
	Tags     []string         `json:"tags"`
	ShortID  string           `json:"short_id"`
	Summary  *SnapshotSummary `json:"summary,omitempty"`
}

// RepositoryStats represents the output of 'restic stats'
type RepositoryStats struct {
	TotalSize              uint64  `json:"total_size"`
	TotalFileCount         uint64  `json:"total_file_count"`
	TotalBlobCount         uint64  `json:"total_blob_count"`
	SnapshotsCount         uint64  `json:"snapshots_count"`
	TotalUncompressedSize  uint64  `json:"total_uncompressed_size"`
	CompressionRatio       float64 `json:"compression_ratio"`
	CompressionProgress    float64 `json:"compression_progress"`
	CompressionSpaceSaving float64 `json:"compression_space_saving"`
	DiskSize               uint64  `json:"disk_size"` // Actual size on disk (from --mode raw-data)
}

// FindResult represents the output of 'restic find'
type FindResult struct {
	Hits     uint64      `json:"hits"`
	Snapshot string      `json:"snapshot"`
	Matches  []FindMatch `json:"matches"`
}

// FindMatch represents a single match in find results
type FindMatch struct {
	Path        string `json:"path"`
	Permissions string `json:"permissions"`
	Type        string `json:"type"`
	Mode        uint32 `json:"mode"`
	Mtime       string `json:"mtime"`
	Atime       string `json:"atime"`
	Ctime       string `json:"ctime"`
	UID         uint32 `json:"uid"`
	GID         uint32 `json:"gid"`
	User        string `json:"user"`
	Group       string `json:"group"`
	Inode       uint64 `json:"inode"`
	DeviceID    uint64 `json:"device_id"`
	Size        uint64 `json:"size"`
	Links       uint64 `json:"links"`
}

// LSNode represents a file/directory in snapshot listing (restic ls)
type LSNode struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Path        string `json:"path"`
	UID         uint32 `json:"uid"`
	GID         uint32 `json:"gid"`
	Size        uint64 `json:"size"`
	Mode        uint32 `json:"mode"`
	Permissions string `json:"permissions"`
	Mtime       string `json:"mtime"`
	Atime       string `json:"atime"`
	Ctime       string `json:"ctime"`
	Inode       uint64 `json:"inode"`
}

// RestoreMessage represents a message from restore command (JSON lines format)
type RestoreMessage struct {
	MessageType    string  `json:"message_type"`
	SecondsElapsed uint64  `json:"seconds_elapsed"`
	PercentDone    float64 `json:"percent_done"`
	TotalFiles     uint64  `json:"total_files"`
	FilesRestored  uint64  `json:"files_restored"`
	FilesSkipped   uint64  `json:"files_skipped"`
	FilesDeleted   uint64  `json:"files_deleted"`
	TotalBytes     uint64  `json:"total_bytes"`
	BytesRestored  uint64  `json:"bytes_restored"`
	BytesSkipped   uint64  `json:"bytes_skipped"`
}
