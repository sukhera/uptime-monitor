class IncidentManager {
    constructor() {
        this.incidentsContainer = document.getElementById('incidents-container');
        this.maintenanceContainer = document.getElementById('maintenance-container');
        this.incidents = [];
        this.maintenanceSchedule = [];
    }

    init() {
        this.loadIncidents();
        this.loadMaintenance();
    }

    async loadIncidents() {
        try {
            const incidents = this.generateSampleIncidents();
            this.incidents = incidents;
            this.renderIncidents();
        } catch (error) {
            console.error('Failed to load incidents:', error);
        }
    }

    async loadMaintenance() {
        try {
            const maintenance = this.generateSampleMaintenance();
            this.maintenanceSchedule = maintenance;
            this.renderMaintenance();
        } catch (error) {
            console.error('Failed to load maintenance:', error);
        }
    }

    generateSampleIncidents() {
        const now = new Date();
        const sampleIncidents = [
            {
                id: '1',
                title: 'API Response Delays',
                description: 'Some users experienced slower API response times due to increased traffic.',
                status: 'resolved',
                severity: 'minor',
                startTime: new Date(now - 2 * 24 * 60 * 60 * 1000),
                endTime: new Date(now - 2 * 24 * 60 * 60 * 1000 + 3 * 60 * 60 * 1000),
                affectedServices: ['API', 'Dashboard']
            },
            {
                id: '2',
                title: 'Database Connection Issues',
                description: 'Brief database connectivity issues affecting user authentication.',
                status: 'resolved',
                severity: 'major',
                startTime: new Date(now - 7 * 24 * 60 * 60 * 1000),
                endTime: new Date(now - 7 * 24 * 60 * 60 * 1000 + 45 * 60 * 1000),
                affectedServices: ['Authentication', 'User Management']
            }
        ];

        return sampleIncidents.slice(0, Math.random() > 0.5 ? 2 : 0);
    }

    generateSampleMaintenance() {
        const now = new Date();
        const sampleMaintenance = [
            {
                id: '1',
                title: 'Database Maintenance',
                description: 'Scheduled database optimization and security updates.',
                status: 'scheduled',
                startTime: new Date(now.getTime() + 3 * 24 * 60 * 60 * 1000),
                endTime: new Date(now.getTime() + 3 * 24 * 60 * 60 * 1000 + 2 * 60 * 60 * 1000),
                affectedServices: ['All Services'],
                impact: 'Brief downtime expected'
            }
        ];

        return sampleMaintenance.slice(0, Math.random() > 0.7 ? 1 : 0);
    }

    renderIncidents() {
        if (!this.incidentsContainer) return;

        if (this.incidents.length === 0) {
            this.incidentsContainer.innerHTML = `
                <div class="no-incidents">
                    <svg class="check-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="20,6 9,17 4,12"></polyline>
                    </svg>
                    <p>No recent incidents reported</p>
                </div>
            `;
            return;
        }

        this.incidentsContainer.innerHTML = `
            <div class="incidents-timeline">
                ${this.incidents.map(incident => this.renderIncident(incident)).join('')}
            </div>
        `;
    }

    renderIncident(incident) {
        const duration = this.formatDuration(incident.startTime, incident.endTime);
        const timeAgo = this.formatTimeAgo(incident.startTime);
        
        return `
            <div class="incident-item severity-${incident.severity}">
                <div class="incident-header">
                    <div class="incident-title">
                        <h4>${this.escapeHtml(incident.title)}</h4>
                        <span class="incident-severity severity-${incident.severity}">${incident.severity}</span>
                    </div>
                    <div class="incident-status status-${incident.status}">
                        ${this.formatStatus(incident.status)}
                    </div>
                </div>
                <div class="incident-details">
                    <p>${this.escapeHtml(incident.description)}</p>
                    <div class="incident-meta">
                        <span class="incident-time">Started ${timeAgo} • Duration: ${duration}</span>
                        <span class="affected-services">
                            Affected: ${incident.affectedServices.join(', ')}
                        </span>
                    </div>
                </div>
            </div>
        `;
    }

    renderMaintenance() {
        if (!this.maintenanceContainer) return;

        if (this.maintenanceSchedule.length === 0) {
            this.maintenanceContainer.innerHTML = `
                <div class="no-maintenance">
                    <svg class="calendar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                        <line x1="16" y1="2" x2="16" y2="6"></line>
                        <line x1="8" y1="2" x2="8" y2="6"></line>
                        <line x1="3" y1="10" x2="21" y2="10"></line>
                    </svg>
                    <p>No scheduled maintenance</p>
                </div>
            `;
            return;
        }

        this.maintenanceContainer.innerHTML = `
            <div class="maintenance-timeline">
                ${this.maintenanceSchedule.map(maintenance => this.renderMaintenanceItem(maintenance)).join('')}
            </div>
        `;
    }

    renderMaintenanceItem(maintenance) {
        const duration = this.formatDuration(maintenance.startTime, maintenance.endTime);
        const timeUntil = this.formatTimeUntil(maintenance.startTime);
        
        return `
            <div class="maintenance-item status-${maintenance.status}">
                <div class="maintenance-header">
                    <div class="maintenance-title">
                        <h4>${this.escapeHtml(maintenance.title)}</h4>
                    </div>
                    <div class="maintenance-status status-${maintenance.status}">
                        ${this.formatStatus(maintenance.status)}
                    </div>
                </div>
                <div class="maintenance-details">
                    <p>${this.escapeHtml(maintenance.description)}</p>
                    <div class="maintenance-meta">
                        <span class="maintenance-time">
                            ${timeUntil} • Duration: ${duration}
                        </span>
                        <span class="affected-services">
                            Affected: ${maintenance.affectedServices.join(', ')}
                        </span>
                        ${maintenance.impact ? `<span class="maintenance-impact">${this.escapeHtml(maintenance.impact)}</span>` : ''}
                    </div>
                </div>
            </div>
        `;
    }

    formatStatus(status) {
        const statusMap = {
            'resolved': 'Resolved',
            'investigating': 'Investigating',
            'identified': 'Identified',
            'monitoring': 'Monitoring',
            'scheduled': 'Scheduled',
            'in-progress': 'In Progress',
            'completed': 'Completed'
        };
        return statusMap[status] || status;
    }

    formatDuration(startTime, endTime) {
        const diff = endTime - startTime;
        const hours = Math.floor(diff / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        
        if (hours > 0) {
            return `${hours}h ${minutes}m`;
        }
        return `${minutes}m`;
    }

    formatTimeAgo(timestamp) {
        const now = new Date();
        const diff = now - timestamp;
        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor(diff / (1000 * 60 * 60));
        const minutes = Math.floor(diff / (1000 * 60));
        
        if (days > 0) {
            return `${days} day${days > 1 ? 's' : ''} ago`;
        } else if (hours > 0) {
            return `${hours} hour${hours > 1 ? 's' : ''} ago`;
        } else if (minutes > 0) {
            return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
        }
        return 'Just now';
    }

    formatTimeUntil(timestamp) {
        const now = new Date();
        const diff = timestamp - now;
        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor(diff / (1000 * 60 * 60));
        
        if (days > 0) {
            return `Scheduled in ${days} day${days > 1 ? 's' : ''}`;
        } else if (hours > 0) {
            return `Scheduled in ${hours} hour${hours > 1 ? 's' : ''}`;
        }
        return 'Starting soon';
    }

    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

export default IncidentManager;