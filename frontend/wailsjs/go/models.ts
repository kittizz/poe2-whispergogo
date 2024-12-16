export namespace main {
	
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
	export class VerifyChatResponse {
	    chat_id: string;
	    valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VerifyChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chat_id = source["chat_id"];
	        this.valid = source["valid"];
	    }
	}

}

