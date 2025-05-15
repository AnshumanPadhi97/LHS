export namespace models {
	
	export class StackUIDto {
	    Id: number;
	    Name: string;
	    Created_at: string;
	    Status: string;
	    Service_count: number;
	
	    static createFrom(source: any = {}) {
	        return new StackUIDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Created_at = source["Created_at"];
	        this.Status = source["Status"];
	        this.Service_count = source["Service_count"];
	    }
	}

}

