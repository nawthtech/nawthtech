// backend/worker/src/cloudflare/index.ts

export interface Env {
    CLOUDFLARE_ZONE_ID: string;
    CLOUDFLARE_API_KEY: string;
    CLOUDFLARE_EMAIL: string;

    // D1
    DB: D1Database;
}

// ==========================
// Cloudflare API Wrapper
// ==========================

export class CloudflareService {
    private zoneId: string;
    private apiKey: string;
    private email: string;
    private baseUrl: string;

    constructor(env: Env) {
        this.zoneId = env.CLOUDFLARE_ZONE_ID;
        this.apiKey = env.CLOUDFLARE_API_KEY;
        this.email = env.CLOUDFLARE_EMAIL;

        this.baseUrl = `https://api.cloudflare.com/client/v4/zones/${this.zoneId}`;
    }

    // ---------------------------
    // Generic Request Handler
    // ---------------------------
    private async request(method: string, endpoint: string, body?: any) {
        const url = `${this.baseUrl}${endpoint}`;

        const req: RequestInit = {
            method,
            headers: {
                "Content-Type": "application/json",
                "X-Auth-Key": this.apiKey,
                "X-Auth-Email": this.email
            }
        };

        if (body) req.body = JSON.stringify(body);

        const response = await fetch(url, req);
        const data = await response.json();

        if (!response.ok || !data.success) {
            throw new Error(`Cloudflare API Error: ${JSON.stringify(data.errors)}`);
        }

        return data.result;
    }

    // ============================
    // Cache Purge functions
    // ============================

    async purgeFiles(files: string[]) {
        return this.request("POST", "/purge_cache", { files });
    }

    async purgeEverything() {
        return this.request("POST", "/purge_cache", { purge_everything: true });
    }

    async purgeTags(tags: string[]) {
        return this.request("POST", "/purge_cache", { tags });
    }

    // ============================
    // Analytics
    // ============================

    async getZoneAnalytics(start: string, end: string) {
        const endpoint = `/analytics/dashboard?since=${start}&until=${end}`;
        return this.request("GET", endpoint);
    }

    async getZoneDetails() {
        return this.request("GET", "");
    }

    // ============================
    // Security Events
    // ============================

    async getSecurityEvents(start: string, end: string) {
        const endpoint = `/security/events?since=${start}&until=${end}`;
        return this.request("GET", endpoint);
    }

    // ============================
    // Firewall Rules
    // ============================

    async createFirewallRule(rule: any) {
        return this.request("POST", "/firewall/rules", rule);
    }

    // ============================
    // Health Check
    // ============================

    async health() {
        try {
            await this.getZoneDetails();
            return {
                service: "cloudflare",
                status: "healthy",
                enabled: true,
            };
        } catch (error: any) {
            return {
                service: "cloudflare",
                status: "error",
                enabled: true,
                error: error.message
            };
        }
    }
}