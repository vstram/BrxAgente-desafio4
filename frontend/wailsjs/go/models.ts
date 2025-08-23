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

export namespace intelligence {
	
	export class Alternative {
	    title: string;
	    description: string;
	    impact: ImpactAssessment;
	    effort: number;
	    risk: number;
	
	    static createFrom(source: any = {}) {
	        return new Alternative(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.description = source["description"];
	        this.impact = this.convertValues(source["impact"], ImpactAssessment);
	        this.effort = source["effort"];
	        this.risk = source["risk"];
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
	export class RequiredResource {
	    type: string;
	    description: string;
	    quantity: number;
	    unit: string;
	    cost?: number;
	
	    static createFrom(source: any = {}) {
	        return new RequiredResource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.description = source["description"];
	        this.quantity = source["quantity"];
	        this.unit = source["unit"];
	        this.cost = source["cost"];
	    }
	}
	export class ImplementationStep {
	    id: string;
	    title: string;
	    description: string;
	    duration: string;
	    order: number;
	    dependencies: string[];
	    resources: RequiredResource[];
	    status: string;
	    progress: number;
	
	    static createFrom(source: any = {}) {
	        return new ImplementationStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.duration = source["duration"];
	        this.order = source["order"];
	        this.dependencies = source["dependencies"];
	        this.resources = this.convertValues(source["resources"], RequiredResource);
	        this.status = source["status"];
	        this.progress = source["progress"];
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
	export class ImplementationGuide {
	    steps: ImplementationStep[];
	    resources: RequiredResource[];
	    timeline: string;
	    risks: string[];
	    success_criteria: string[];
	    rollback_plan: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImplementationGuide(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.steps = this.convertValues(source["steps"], ImplementationStep);
	        this.resources = this.convertValues(source["resources"], RequiredResource);
	        this.timeline = source["timeline"];
	        this.risks = source["risks"];
	        this.success_criteria = source["success_criteria"];
	        this.rollback_plan = source["rollback_plan"];
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
	export class ImpactDetail {
	    value: number;
	    unit: string;
	    description: string;
	    confidence: number;
	
	    static createFrom(source: any = {}) {
	        return new ImpactDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.unit = source["unit"];
	        this.description = source["description"];
	        this.confidence = source["confidence"];
	    }
	}
	export class ImpactAssessment {
	    financial: ImpactDetail;
	    operational: ImpactDetail;
	    strategic: ImpactDetail;
	    risk: ImpactDetail;
	    timeline: string;
	    probability: number;
	
	    static createFrom(source: any = {}) {
	        return new ImpactAssessment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.financial = this.convertValues(source["financial"], ImpactDetail);
	        this.operational = this.convertValues(source["operational"], ImpactDetail);
	        this.strategic = this.convertValues(source["strategic"], ImpactDetail);
	        this.risk = this.convertValues(source["risk"], ImpactDetail);
	        this.timeline = source["timeline"];
	        this.probability = source["probability"];
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
	export class Evidence {
	    type: string;
	    source: string;
	    description: string;
	    value: number;
	    confidence: number;
	    metadata: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Evidence(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.source = source["source"];
	        this.description = source["description"];
	        this.value = source["value"];
	        this.confidence = source["confidence"];
	        this.metadata = source["metadata"];
	    }
	}
	export class RecommendationContext {
	    source_analysis: string;
	    trigger_condition: string;
	    affected_entities: string[];
	    time_horizon: string;
	    seasonality?: predicoes.SeasonalityInfo;
	    trend_info?: predicoes.VRTrend;
	
	    static createFrom(source: any = {}) {
	        return new RecommendationContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source_analysis = source["source_analysis"];
	        this.trigger_condition = source["trigger_condition"];
	        this.affected_entities = source["affected_entities"];
	        this.time_horizon = source["time_horizon"];
	        this.seasonality = this.convertValues(source["seasonality"], predicoes.SeasonalityInfo);
	        this.trend_info = this.convertValues(source["trend_info"], predicoes.VRTrend);
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
	export class EnhancedRecommendation {
	    id: string;
	    priority: string;
	    category: string;
	    title: string;
	    description: string;
	    due_date?: time.Time;
	    owner?: string;
	    status: string;
	    context: RecommendationContext;
	    evidence: Evidence[];
	    expected_impact: ImpactAssessment;
	    implementation: ImplementationGuide;
	    dependencies: string[];
	    alternatives: Alternative[];
	    confidence: number;
	    auto_approved: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EnhancedRecommendation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.priority = source["priority"];
	        this.category = source["category"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.due_date = this.convertValues(source["due_date"], time.Time);
	        this.owner = source["owner"];
	        this.status = source["status"];
	        this.context = this.convertValues(source["context"], RecommendationContext);
	        this.evidence = this.convertValues(source["evidence"], Evidence);
	        this.expected_impact = this.convertValues(source["expected_impact"], ImpactAssessment);
	        this.implementation = this.convertValues(source["implementation"], ImplementationGuide);
	        this.dependencies = source["dependencies"];
	        this.alternatives = this.convertValues(source["alternatives"], Alternative);
	        this.confidence = source["confidence"];
	        this.auto_approved = source["auto_approved"];
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
	export class ActionPlan {
	    quick_wins: EnhancedRecommendation[];
	    short_term: EnhancedRecommendation[];
	    medium_term: EnhancedRecommendation[];
	    long_term: EnhancedRecommendation[];
	    dependencies: Record<string, Array<string>>;
	    timeline: string;
	
	    static createFrom(source: any = {}) {
	        return new ActionPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.quick_wins = this.convertValues(source["quick_wins"], EnhancedRecommendation);
	        this.short_term = this.convertValues(source["short_term"], EnhancedRecommendation);
	        this.medium_term = this.convertValues(source["medium_term"], EnhancedRecommendation);
	        this.long_term = this.convertValues(source["long_term"], EnhancedRecommendation);
	        this.dependencies = source["dependencies"];
	        this.timeline = source["timeline"];
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
	
	export class ConsumptionPattern {
	    sindicato: string;
	    pattern_type: string;
	    characteristics: Record<string, number>;
	    trend?: predicoes.VRTrend;
	    seasonality?: predicoes.SeasonalityInfo;
	    stability: number;
	    predictability: number;
	
	    static createFrom(source: any = {}) {
	        return new ConsumptionPattern(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sindicato = source["sindicato"];
	        this.pattern_type = source["pattern_type"];
	        this.characteristics = source["characteristics"];
	        this.trend = this.convertValues(source["trend"], predicoes.VRTrend);
	        this.seasonality = this.convertValues(source["seasonality"], predicoes.SeasonalityInfo);
	        this.stability = source["stability"];
	        this.predictability = source["predictability"];
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
	
	
	
	
	
	
	export class Opportunity {
	    id: string;
	    title: string;
	    description: string;
	    value: number;
	    effort: number;
	    timeline: string;
	    actions: predicoes.ActionItem[];
	
	    static createFrom(source: any = {}) {
	        return new Opportunity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.value = source["value"];
	        this.effort = source["effort"];
	        this.timeline = source["timeline"];
	        this.actions = this.convertValues(source["actions"], predicoes.ActionItem);
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
	
	export class RiskFactor {
	    id: string;
	    type: string;
	    description: string;
	    probability: number;
	    impact: string;
	    mitigation: predicoes.ActionItem[];
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new RiskFactor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.description = source["description"];
	        this.probability = source["probability"];
	        this.impact = source["impact"];
	        this.mitigation = this.convertValues(source["mitigation"], predicoes.ActionItem);
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
	export class RecommendationSummary {
	    total_recommendations: number;
	    by_priority: Record<string, number>;
	    by_category: Record<string, number>;
	    expected_value: number;
	    total_effort: number;
	    average_confidence: number;
	    top_opportunities: string[];
	    critical_risks: string[];
	
	    static createFrom(source: any = {}) {
	        return new RecommendationSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_recommendations = source["total_recommendations"];
	        this.by_priority = source["by_priority"];
	        this.by_category = source["by_category"];
	        this.expected_value = source["expected_value"];
	        this.total_effort = source["total_effort"];
	        this.average_confidence = source["average_confidence"];
	        this.top_opportunities = source["top_opportunities"];
	        this.critical_risks = source["critical_risks"];
	    }
	}
	export class RecommendationSuite {
	    summary: RecommendationSummary;
	    recommendations: EnhancedRecommendation[];
	    action_plan: ActionPlan;
	    risk_factors: RiskFactor[];
	    opportunities: Opportunity[];
	    generated_at: time.Time;
	    valid_until: time.Time;
	
	    static createFrom(source: any = {}) {
	        return new RecommendationSuite(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.summary = this.convertValues(source["summary"], RecommendationSummary);
	        this.recommendations = this.convertValues(source["recommendations"], EnhancedRecommendation);
	        this.action_plan = this.convertValues(source["action_plan"], ActionPlan);
	        this.risk_factors = this.convertValues(source["risk_factors"], RiskFactor);
	        this.opportunities = this.convertValues(source["opportunities"], Opportunity);
	        this.generated_at = this.convertValues(source["generated_at"], time.Time);
	        this.valid_until = this.convertValues(source["valid_until"], time.Time);
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
	
	
	
	export class TrendAnalysisMetadata {
	    data_points: number;
	    analysis_period: string;
	    methods: string[];
	    parameters: Record<string, any>;
	    quality_score: number;
	
	    static createFrom(source: any = {}) {
	        return new TrendAnalysisMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data_points = source["data_points"];
	        this.analysis_period = source["analysis_period"];
	        this.methods = source["methods"];
	        this.parameters = source["parameters"];
	        this.quality_score = source["quality_score"];
	    }
	}
	export class TrendRecommendation {
	    type: string;
	    priority: string;
	    title: string;
	    description: string;
	    actions: predicoes.ActionItem[];
	    impact: number;
	    timeframe: string;
	
	    static createFrom(source: any = {}) {
	        return new TrendRecommendation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.priority = source["priority"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.actions = this.convertValues(source["actions"], predicoes.ActionItem);
	        this.impact = source["impact"];
	        this.timeframe = source["timeframe"];
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
	export class TrendAnalysisResult {
	    primary_trend?: predicoes.VRTrend;
	    secondary_trends: predicoes.VRTrend[];
	    seasonality?: predicoes.SeasonalityInfo;
	    volatility: number;
	    confidence: number;
	    recommendations: TrendRecommendation[];
	    metadata: TrendAnalysisMetadata;
	
	    static createFrom(source: any = {}) {
	        return new TrendAnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.primary_trend = this.convertValues(source["primary_trend"], predicoes.VRTrend);
	        this.secondary_trends = this.convertValues(source["secondary_trends"], predicoes.VRTrend);
	        this.seasonality = this.convertValues(source["seasonality"], predicoes.SeasonalityInfo);
	        this.volatility = source["volatility"];
	        this.confidence = source["confidence"];
	        this.recommendations = this.convertValues(source["recommendations"], TrendRecommendation);
	        this.metadata = this.convertValues(source["metadata"], TrendAnalysisMetadata);
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
	    timestamp: time.Time;
	    level: string;
	    message: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], time.Time);
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
	    startTime?: time.Time;
	    endTime?: time.Time;
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
	        this.startTime = this.convertValues(source["startTime"], time.Time);
	        this.endTime = this.convertValues(source["endTime"], time.Time);
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
	    startTime: time.Time;
	    endTime?: time.Time;
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
	        this.startTime = this.convertValues(source["startTime"], time.Time);
	        this.endTime = this.convertValues(source["endTime"], time.Time);
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
	    lastUpdated: time.Time;
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
	        this.lastUpdated = this.convertValues(source["lastUpdated"], time.Time);
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
	    startTime: time.Time;
	    endTime: time.Time;
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
	        this.startTime = this.convertValues(source["startTime"], time.Time);
	        this.endTime = this.convertValues(source["endTime"], time.Time);
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

export namespace predicoes {
	
	export class ActionItem {
	    id: string;
	    priority: string;
	    category: string;
	    title: string;
	    description: string;
	    due_date?: time.Time;
	    owner?: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ActionItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.priority = source["priority"];
	        this.category = source["category"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.due_date = this.convertValues(source["due_date"], time.Time);
	        this.owner = source["owner"];
	        this.status = source["status"];
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
	export class VRTrend {
	    type: string;
	    strength: number;
	    period: number;
	    confidence: number;
	    start_date: time.Time;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new VRTrend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.strength = source["strength"];
	        this.period = source["period"];
	        this.confidence = source["confidence"];
	        this.start_date = this.convertValues(source["start_date"], time.Time);
	        this.description = source["description"];
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
	export class SeasonalityInfo {
	    is_detected: boolean;
	    period: number;
	    amplitude: number;
	    peak_months: number[];
	    trough_months: number[];
	    confidence: number;
	    pattern: Record<number, number>;
	
	    static createFrom(source: any = {}) {
	        return new SeasonalityInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_detected = source["is_detected"];
	        this.period = source["period"];
	        this.amplitude = source["amplitude"];
	        this.peak_months = source["peak_months"];
	        this.trough_months = source["trough_months"];
	        this.confidence = source["confidence"];
	        this.pattern = source["pattern"];
	    }
	}
	export class ForecastRange {
	    lower: number;
	    upper: number;
	
	    static createFrom(source: any = {}) {
	        return new ForecastRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lower = source["lower"];
	        this.upper = source["upper"];
	    }
	}
	export class ConsumptionForecast {
	    sindicato: string;
	    month: time.Time;
	    predicted_vr: number;
	    confidence: number;
	    range: ForecastRange;
	    factors: string[];
	    seasonality?: SeasonalityInfo;
	    trend?: VRTrend;
	    assumptions: string[];
	
	    static createFrom(source: any = {}) {
	        return new ConsumptionForecast(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sindicato = source["sindicato"];
	        this.month = this.convertValues(source["month"], time.Time);
	        this.predicted_vr = source["predicted_vr"];
	        this.confidence = source["confidence"];
	        this.range = this.convertValues(source["range"], ForecastRange);
	        this.factors = source["factors"];
	        this.seasonality = this.convertValues(source["seasonality"], SeasonalityInfo);
	        this.trend = this.convertValues(source["trend"], VRTrend);
	        this.assumptions = source["assumptions"];
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
	
	export class HistoricalVRData {
	    month: time.Time;
	    sindicato: string;
	    total_vr: number;
	    num_colaboradores: number;
	    media_por_pessoa: number;
	    days_processed: number;
	    anomalies: string[];
	    metadata: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new HistoricalVRData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = this.convertValues(source["month"], time.Time);
	        this.sindicato = source["sindicato"];
	        this.total_vr = source["total_vr"];
	        this.num_colaboradores = source["num_colaboradores"];
	        this.media_por_pessoa = source["media_por_pessoa"];
	        this.days_processed = source["days_processed"];
	        this.anomalies = source["anomalies"];
	        this.metadata = source["metadata"];
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
	export class PredictionMeta {
	    model: string;
	    features: string[];
	    data_points: number;
	    training_period: string;
	    accuracy: number;
	    parameters: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new PredictionMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.features = source["features"];
	        this.data_points = source["data_points"];
	        this.training_period = source["training_period"];
	        this.accuracy = source["accuracy"];
	        this.parameters = source["parameters"];
	    }
	}
	export class Prediction {
	    id: string;
	    type: string;
	    target: string;
	    value: any;
	    confidence: number;
	    timeframe: string;
	    created_at: time.Time;
	    valid_until: time.Time;
	    method: string;
	    description: string;
	    metadata: PredictionMeta;
	    actions: ActionItem[];
	
	    static createFrom(source: any = {}) {
	        return new Prediction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.target = source["target"];
	        this.value = source["value"];
	        this.confidence = source["confidence"];
	        this.timeframe = source["timeframe"];
	        this.created_at = this.convertValues(source["created_at"], time.Time);
	        this.valid_until = this.convertValues(source["valid_until"], time.Time);
	        this.method = source["method"];
	        this.description = source["description"];
	        this.metadata = this.convertValues(source["metadata"], PredictionMeta);
	        this.actions = this.convertValues(source["actions"], ActionItem);
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

export namespace time {
	
	export class Time {
	
	
	    static createFrom(source: any = {}) {
	        return new Time(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

