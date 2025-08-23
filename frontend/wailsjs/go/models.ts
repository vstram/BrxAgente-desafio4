export namespace config {
	
	export class AgentConfig {
	    enabled: boolean;
	    model: string;
	    temperature: number;
	    max_tokens: number;
	    worker_pool_size: number;
	    cache_enabled: boolean;
	    cache_size: number;
	    tools_enabled: string[];
	
	    static createFrom(source: any = {}) {
	        return new AgentConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.model = source["model"];
	        this.temperature = source["temperature"];
	        this.max_tokens = source["max_tokens"];
	        this.worker_pool_size = source["worker_pool_size"];
	        this.cache_enabled = source["cache_enabled"];
	        this.cache_size = source["cache_size"];
	        this.tools_enabled = source["tools_enabled"];
	    }
	}
	export class OllamaConfig {
	    base_url?: string;
	    model?: string;
	
	    static createFrom(source: any = {}) {
	        return new OllamaConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base_url = source["base_url"];
	        this.model = source["model"];
	    }
	}
	export class Config {
	    openai_key?: string;
	    ollama_config?: OllamaConfig;
	    agent_config?: AgentConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.openai_key = source["openai_key"];
	        this.ollama_config = this.convertValues(source["ollama_config"], OllamaConfig);
	        this.agent_config = this.convertValues(source["agent_config"], AgentConfig);
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

export namespace main {
	
	export class AgentMetrics {
	    totalWorkflowsExecuted: number;
	    successfulWorkflows: number;
	    collaboratorsProcessed: number;
	    reportsGenerated: number;
	    anomaliesDetected: number;
	    uptime: number;
	
	    static createFrom(source: any = {}) {
	        return new AgentMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalWorkflowsExecuted = source["totalWorkflowsExecuted"];
	        this.successfulWorkflows = source["successfulWorkflows"];
	        this.collaboratorsProcessed = source["collaboratorsProcessed"];
	        this.reportsGenerated = source["reportsGenerated"];
	        this.anomaliesDetected = source["anomaliesDetected"];
	        this.uptime = source["uptime"];
	    }
	}
	export class LogEntry {
	    id: string;
	    // Go type: time
	    timestamp: any;
	    level: string;
	    message: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.level = source["level"];
	        this.message = source["message"];
	        this.source = source["source"];
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
	export class WorkflowStep {
	    id: string;
	    name: string;
	    status: string;
	    // Go type: time
	    startTime?: any;
	    // Go type: time
	    endTime?: any;
	    duration: number;
	    errorMsg: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.status = source["status"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.duration = source["duration"];
	        this.errorMsg = source["errorMsg"];
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
	export class WorkflowInfo {
	    id: string;
	    name: string;
	    status: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime?: any;
	    steps: WorkflowStep[];
	    progress: number;
	    errorMsg: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.status = source["status"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.steps = this.convertValues(source["steps"], WorkflowStep);
	        this.progress = source["progress"];
	        this.errorMsg = source["errorMsg"];
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
	export class AgentStatus {
	    status: string;
	    // Go type: time
	    lastUpdated: any;
	    currentWorkflow?: WorkflowInfo;
	    availableWorkflows: string[];
	    metrics: AgentMetrics;
	    recentLogs: LogEntry[];
	
	    static createFrom(source: any = {}) {
	        return new AgentStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.lastUpdated = this.convertValues(source["lastUpdated"], null);
	        this.currentWorkflow = this.convertValues(source["currentWorkflow"], WorkflowInfo);
	        this.availableWorkflows = source["availableWorkflows"];
	        this.metrics = this.convertValues(source["metrics"], AgentMetrics);
	        this.recentLogs = this.convertValues(source["recentLogs"], LogEntry);
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
	
	export class WorkflowExecution {
	    id: string;
	    workflowName: string;
	    status: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime: any;
	    duration: number;
	    collaboratorsProcessed: number;
	    reportsGenerated: number;
	    anomaliesDetected: number;
	    errorMsg: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowExecution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workflowName = source["workflowName"];
	        this.status = source["status"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.duration = source["duration"];
	        this.collaboratorsProcessed = source["collaboratorsProcessed"];
	        this.reportsGenerated = source["reportsGenerated"];
	        this.anomaliesDetected = source["anomaliesDetected"];
	        this.errorMsg = source["errorMsg"];
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
	
	export class WorkflowStartRequest {
	    workflowName: string;
	    parameters: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowStartRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.workflowName = source["workflowName"];
	        this.parameters = source["parameters"];
	    }
	}

}

