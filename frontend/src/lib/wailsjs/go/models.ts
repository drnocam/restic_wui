export namespace restic {
	
	export class RepositoryStats {
	    total_size: number;
	    total_file_count: number;
	    total_blob_count: number;
	
	    static createFrom(source: any = {}) {
	        return new RepositoryStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_size = source["total_size"];
	        this.total_file_count = source["total_file_count"];
	        this.total_blob_count = source["total_blob_count"];
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
	    }
	}

}

