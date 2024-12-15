export namespace main {
	
	export enum Ntfy {
	    NTFY_BASE_URL = "https://ntfy.sh",
	    NTFY_PREFIX_TOPICS = "wpgogo",
	}
	export class Keyword {
	    Keyword: string;
	    Enable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Keyword(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Keyword = source["Keyword"];
	        this.Enable = source["Enable"];
	    }
	}

}

