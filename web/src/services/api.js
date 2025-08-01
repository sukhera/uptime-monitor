class StatusAPI {
    constructor(baseURL = '/api') {
        this.baseURL = baseURL;
    }

    async fetchStatus() {
        try {
            const response = await fetch(`${this.baseURL}/status`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error fetching status:', error);
            throw error;
        }
    }

    async healthCheck() {
        try {
            const response = await fetch(`${this.baseURL}/health`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error checking health:', error);
            throw error;
        }
    }
}

export default StatusAPI;