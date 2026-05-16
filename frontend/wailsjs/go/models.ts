export namespace domain {
	
	export class Alert {
	    id: string;
	    serverId: string;
	    serverName: string;
	    kind: string;
	    value: number;
	    threshold: number;
	    message: string;
	    createdAt: number;
	    read: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serverId = source["serverId"];
	        this.serverName = source["serverName"];
	        this.kind = source["kind"];
	        this.value = source["value"];
	        this.threshold = source["threshold"];
	        this.message = source["message"];
	        this.createdAt = source["createdAt"];
	        this.read = source["read"];
	    }
	}

}

export namespace main {
	
	export class AlertDTO {
	    id: string;
	    serverId: string;
	    serverName: string;
	    kind: string;
	    value: number;
	    threshold: number;
	    message: string;
	    createdAt: number;
	    read: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AlertDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serverId = source["serverId"];
	        this.serverName = source["serverName"];
	        this.kind = source["kind"];
	        this.value = source["value"];
	        this.threshold = source["threshold"];
	        this.message = source["message"];
	        this.createdAt = source["createdAt"];
	        this.read = source["read"];
	    }
	}
	export class ClusterDTO {
	    id: string;
	    name: string;
	    description: string;
	    serverIds: string[];
	    createdAt: number;
	    updatedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new ClusterDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.serverIds = source["serverIds"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class ClusterInputDTO {
	    id: string;
	    name: string;
	    description: string;
	    serverIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new ClusterInputDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.serverIds = source["serverIds"];
	    }
	}
	export class ClusterSummaryDTO {
	    id: string;
	    name: string;
	    description: string;
	    serverIds: string[];
	    createdAt: number;
	    updatedAt: number;
	    totalServers: number;
	    onlineCount: number;
	    offlineCount: number;
	
	    static createFrom(source: any = {}) {
	        return new ClusterSummaryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.serverIds = source["serverIds"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.totalServers = source["totalServers"];
	        this.onlineCount = source["onlineCount"];
	        this.offlineCount = source["offlineCount"];
	    }
	}
	export class DashboardStatsDTO {
	    totalServers: number;
	    onlineServers: number;
	    totalClusters: number;
	    unreadAlerts: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalServers = source["totalServers"];
	        this.onlineServers = source["onlineServers"];
	        this.totalClusters = source["totalClusters"];
	        this.unreadAlerts = source["unreadAlerts"];
	    }
	}
	export class MetricsSampleDTO {
	    serverId: string;
	    timestamp: number;
	    cpuPercent: number;
	    memPercent: number;
	    diskPercent: number;
	    memAvailBytes: number;
	    diskFreeBytes: number;
	
	    static createFrom(source: any = {}) {
	        return new MetricsSampleDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serverId = source["serverId"];
	        this.timestamp = source["timestamp"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memPercent = source["memPercent"];
	        this.diskPercent = source["diskPercent"];
	        this.memAvailBytes = source["memAvailBytes"];
	        this.diskFreeBytes = source["diskFreeBytes"];
	    }
	}
	export class ServerDTO {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    role: string;
	    status: string;
	    tags: string[];
	    checkType: string;
	    checkPath: string;
	    lastCheck: number;
	    latencyMs: number;
	    clusterId: string;
	    notes: string;
	    useAgent: boolean;
	    agentPort: number;
	    agentToken: string;
	    cpuThreshold: number;
	    memThreshold: number;
	    diskThreshold: number;
	    cpuPercent: number;
	    memPercent: number;
	    diskPercent: number;
	    memAvailBytes: number;
	    diskFreeBytes: number;
	
	    static createFrom(source: any = {}) {
	        return new ServerDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.role = source["role"];
	        this.status = source["status"];
	        this.tags = source["tags"];
	        this.checkType = source["checkType"];
	        this.checkPath = source["checkPath"];
	        this.lastCheck = source["lastCheck"];
	        this.latencyMs = source["latencyMs"];
	        this.clusterId = source["clusterId"];
	        this.notes = source["notes"];
	        this.useAgent = source["useAgent"];
	        this.agentPort = source["agentPort"];
	        this.agentToken = source["agentToken"];
	        this.cpuThreshold = source["cpuThreshold"];
	        this.memThreshold = source["memThreshold"];
	        this.diskThreshold = source["diskThreshold"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memPercent = source["memPercent"];
	        this.diskPercent = source["diskPercent"];
	        this.memAvailBytes = source["memAvailBytes"];
	        this.diskFreeBytes = source["diskFreeBytes"];
	    }
	}
	export class ServerInputDTO {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    role: string;
	    tags: string[];
	    checkType: string;
	    checkPath: string;
	    clusterId: string;
	    notes: string;
	    useAgent: boolean;
	    agentPort: number;
	    agentToken: string;
	    cpuThreshold: number;
	    memThreshold: number;
	    diskThreshold: number;
	
	    static createFrom(source: any = {}) {
	        return new ServerInputDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.role = source["role"];
	        this.tags = source["tags"];
	        this.checkType = source["checkType"];
	        this.checkPath = source["checkPath"];
	        this.clusterId = source["clusterId"];
	        this.notes = source["notes"];
	        this.useAgent = source["useAgent"];
	        this.agentPort = source["agentPort"];
	        this.agentToken = source["agentToken"];
	        this.cpuThreshold = source["cpuThreshold"];
	        this.memThreshold = source["memThreshold"];
	        this.diskThreshold = source["diskThreshold"];
	    }
	}

}

