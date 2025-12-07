export namespace config {
	
	export class Config {
	    ResticCommand: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ResticCommand = source["ResticCommand"];
	    }
	}

}

export namespace restic {
	
	export class FindMatch {
	    path: string;
	    permissions: string;
	    type: string;
	    mode: number;
	    mtime: string;
	    atime: string;
	    ctime: string;
	    uid: number;
	    gid: number;
	    user: string;
	    group: string;
	    inode: number;
	    device_id: number;
	    size: number;
	    links: number;
	
	    static createFrom(source: any = {}) {
	        return new FindMatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.permissions = source["permissions"];
	        this.type = source["type"];
	        this.mode = source["mode"];
	        this.mtime = source["mtime"];
	        this.atime = source["atime"];
	        this.ctime = source["ctime"];
	        this.uid = source["uid"];
	        this.gid = source["gid"];
	        this.user = source["user"];
	        this.group = source["group"];
	        this.inode = source["inode"];
	        this.device_id = source["device_id"];
	        this.size = source["size"];
	        this.links = source["links"];
	    }
	}
	export class FindResult {
	    hits: number;
	    snapshot: string;
	    matches: FindMatch[];
	
	    static createFrom(source: any = {}) {
	        return new FindResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hits = source["hits"];
	        this.snapshot = source["snapshot"];
	        this.matches = this.convertValues(source["matches"], FindMatch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LSNode {
	    name: string;
	    type: string;
	    path: string;
	    uid: number;
	    gid: number;
	    size: number;
	    mode: number;
	    permissions: string;
	    mtime: string;
	    atime: string;
	    ctime: string;
	    inode: number;
	
	    static createFrom(source: any = {}) {
	        return new LSNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.path = source["path"];
	        this.uid = source["uid"];
	        this.gid = source["gid"];
	        this.size = source["size"];
	        this.mode = source["mode"];
	        this.permissions = source["permissions"];
	        this.mtime = source["mtime"];
	        this.atime = source["atime"];
	        this.ctime = source["ctime"];
	        this.inode = source["inode"];
	    }
	}
	export class RepositoryStats {
	    total_size: number;
	    total_file_count: number;
	    total_blob_count: number;
	    snapshots_count: number;
	    total_uncompressed_size: number;
	    compression_ratio: number;
	    compression_progress: number;
	    compression_space_saving: number;
	    disk_size: number;
	
	    static createFrom(source: any = {}) {
	        return new RepositoryStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_size = source["total_size"];
	        this.total_file_count = source["total_file_count"];
	        this.total_blob_count = source["total_blob_count"];
	        this.snapshots_count = source["snapshots_count"];
	        this.total_uncompressed_size = source["total_uncompressed_size"];
	        this.compression_ratio = source["compression_ratio"];
	        this.compression_progress = source["compression_progress"];
	        this.compression_space_saving = source["compression_space_saving"];
	        this.disk_size = source["disk_size"];
	    }
	}
	export class RestoreMessage {
	    message_type: string;
	    seconds_elapsed: number;
	    percent_done: number;
	    total_files: number;
	    files_restored: number;
	    files_skipped: number;
	    files_deleted: number;
	    total_bytes: number;
	    bytes_restored: number;
	    bytes_skipped: number;
	
	    static createFrom(source: any = {}) {
	        return new RestoreMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message_type = source["message_type"];
	        this.seconds_elapsed = source["seconds_elapsed"];
	        this.percent_done = source["percent_done"];
	        this.total_files = source["total_files"];
	        this.files_restored = source["files_restored"];
	        this.files_skipped = source["files_skipped"];
	        this.files_deleted = source["files_deleted"];
	        this.total_bytes = source["total_bytes"];
	        this.bytes_restored = source["bytes_restored"];
	        this.bytes_skipped = source["bytes_skipped"];
	    }
	}
	export class SnapshotSummary {
	    files_new: number;
	    files_changed: number;
	    files_unmodified: number;
	    dirs_new: number;
	    dirs_changed: number;
	    dirs_unmodified: number;
	    data_blobs: number;
	    tree_blobs: number;
	    data_added: number;
	    total_files_processed: number;
	    total_bytes_processed: number;
	
	    static createFrom(source: any = {}) {
	        return new SnapshotSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files_new = source["files_new"];
	        this.files_changed = source["files_changed"];
	        this.files_unmodified = source["files_unmodified"];
	        this.dirs_new = source["dirs_new"];
	        this.dirs_changed = source["dirs_changed"];
	        this.dirs_unmodified = source["dirs_unmodified"];
	        this.data_blobs = source["data_blobs"];
	        this.tree_blobs = source["tree_blobs"];
	        this.data_added = source["data_added"];
	        this.total_files_processed = source["total_files_processed"];
	        this.total_bytes_processed = source["total_bytes_processed"];
	    }
	}
	export class Snapshot {
	    id: string;
	    time: string;
	    tree: string;
	    paths: string[];
	    hostname: string;
	    username: string;
	    tags: string[];
	    short_id: string;
	    summary?: SnapshotSummary;
	
	    static createFrom(source: any = {}) {
	        return new Snapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.time = source["time"];
	        this.tree = source["tree"];
	        this.paths = source["paths"];
	        this.hostname = source["hostname"];
	        this.username = source["username"];
	        this.tags = source["tags"];
	        this.short_id = source["short_id"];
	        this.summary = this.convertValues(source["summary"], SnapshotSummary);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

