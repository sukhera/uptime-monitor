class StatusDashboard {
    constructor(containerId, api) {
        this.container = document.getElementById(containerId);
        this.api = api;
        this.refreshInterval = null;
        this.retryCount = 0;
        this.maxRetries = 3;
        this.services = [];
        this.theme = localStorage.getItem('theme') || 'light';
        
        this.bindEvents();
        this.initTheme();
    }

    bindEvents() {
        const refreshBtn = document.getElementById('refresh-btn');
        const themeToggle = document.getElementById('theme-toggle');
        const subscribeBtn = document.getElementById('subscribe-btn');

        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => this.loadStatus(true));
        }

        if (themeToggle) {
            themeToggle.addEventListener('click', () => this.toggleTheme());
        }

        if (subscribeBtn) {
            subscribeBtn.addEventListener('click', () => this.showSubscribeModal());
        }

        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.stopAutoRefresh();
            } else {
                this.startAutoRefresh();
                this.loadStatus();
            }
        });
    }

    initTheme() {
        document.documentElement.setAttribute('data-theme', this.theme);
    }

    toggleTheme() {
        this.theme = this.theme === 'light' ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', this.theme);
        localStorage.setItem('theme', this.theme);
    }

    init() {
        this.loadStatus();
        this.startAutoRefresh();
    }

    async loadStatus(isManualRefresh = false) {
        try {
            if (!isManualRefresh) {
                this.showLoading();
            } else {
                this.showRefreshing();
            }
            
            const services = await this.api.fetchStatus();
            this.services = services;
            this.updateUI(services);
            this.updateOverallStatus(services);
            this.retryCount = 0;
            
        } catch (error) {
            console.error('Failed to load status:', error);
            this.handleLoadError(error);
        }
    }

    handleLoadError(error) {
        this.retryCount++;
        
        if (this.retryCount <= this.maxRetries) {
            const retryDelay = Math.min(1000 * Math.pow(2, this.retryCount - 1), 10000);
            this.showError(`Failed to load service status. Retrying in ${retryDelay / 1000} seconds...`);
            
            setTimeout(() => this.loadStatus(), retryDelay);
        } else {
            this.showError('Failed to load service status after multiple attempts. Please check your connection and try again.');
        }
    }

    updateOverallStatus(services) {
        const overallStatusEl = document.getElementById('overall-status');
        if (!overallStatusEl) return;

        const statusIndicator = overallStatusEl.querySelector('.status-indicator');
        const statusText = overallStatusEl.querySelector('.status-text');

        if (services.length === 0) {
            statusText.textContent = 'No Services';
            return;
        }

        const downServices = services.filter(s => s.status === 'down').length;
        const degradedServices = services.filter(s => s.status === 'degraded').length;

        if (downServices > 0) {
            statusIndicator.style.background = 'var(--status-down)';
            statusText.textContent = `${downServices} Service${downServices > 1 ? 's' : ''} Down`;
        } else if (degradedServices > 0) {
            statusIndicator.style.background = 'var(--status-degraded)';
            statusText.textContent = `${degradedServices} Service${degradedServices > 1 ? 's' : ''} Degraded`;
        } else {
            statusIndicator.style.background = 'var(--status-operational)';
            statusText.textContent = 'All Systems Operational';
        }
    }

    updateUI(services) {
        if (services.length === 0) {
            this.container.innerHTML = `
                <div class="loading">
                    <div style="font-size: 2rem; margin-bottom: 1rem;">🔍</div>
                    <span>No services found</span>
                </div>
            `;
            return;
        }

        this.container.innerHTML = `
            <div class="status-grid">
                ${services.map(service => this.renderServiceCard(service)).join('')}
            </div>
        `;
    }

    renderServiceCard(service) {
        const statusClass = `status-${service.status}`;
        const responseTime = service.latency_ms;
        const responseTimeClass = responseTime > 1000 ? 'slow' : responseTime > 500 ? 'medium' : 'fast';
        
        return `
            <div class="service-card ${statusClass}" data-service="${this.escapeHtml(service.name)}">
                <div class="service-header">
                    <div class="service-name">${this.escapeHtml(service.name)}</div>
                    <div class="status-badge ${statusClass}">
                        ${this.formatStatus(service.status)}
                    </div>
                </div>
                <div class="service-details">
                    <div class="latency ${responseTimeClass}">
                        <span>Response time: ${responseTime}ms</span>
                    </div>
                    <div class="last-updated">
                        <span>Last updated: ${this.formatTime(service.updated_at)}</span>
                    </div>
                    ${service.description ? `<div class="service-description">${this.escapeHtml(service.description)}</div>` : ''}
                </div>
            </div>
        `;
    }

    formatStatus(status) {
        const statusMap = {
            'operational': 'Operational',
            'degraded': 'Degraded',
            'down': 'Down'
        };
        return statusMap[status] || status;
    }

    showLoading() {
        this.container.innerHTML = `
            <div class="loading" role="status" aria-label="Loading service status">
                <div class="loading-spinner"></div>
                <span>Loading service status...</span>
            </div>
        `;
    }

    showRefreshing() {
        const refreshBtn = document.getElementById('refresh-btn');
        if (refreshBtn) {
            refreshBtn.style.transform = 'rotate(360deg)';
            setTimeout(() => {
                refreshBtn.style.transform = '';
            }, 500);
        }
    }

    showError(message) {
        this.container.innerHTML = `
            <div class="error" role="alert">
                <span>${this.escapeHtml(message)}</span>
                <button onclick="window.statusDashboard.loadStatus(true)" style="margin-left: 1rem; padding: 0.5rem 1rem; background: var(--status-down); color: white; border: none; border-radius: 0.5rem; cursor: pointer;">
                    Retry Now
                </button>
            </div>
        `;
    }

    showSubscribeModal() {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.innerHTML = `
            <div class="modal">
                <div class="modal-header">
                    <h3>Get Status Updates</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <p>Subscribe to receive notifications about service outages and maintenance.</p>
                    <form class="subscribe-form">
                        <input type="email" placeholder="Enter your email" required>
                        <button type="submit">Subscribe</button>
                    </form>
                </div>
            </div>
        `;

        document.body.appendChild(modal);

        modal.querySelector('.modal-close').addEventListener('click', () => {
            document.body.removeChild(modal);
        });

        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                document.body.removeChild(modal);
            }
        });

        modal.querySelector('.subscribe-form').addEventListener('submit', (e) => {
            e.preventDefault();
            const email = e.target.querySelector('input').value;
            this.handleSubscription(email);
            document.body.removeChild(modal);
        });
    }

    handleSubscription(email) {
        console.log('Subscription request for:', email);
        
        const notification = document.createElement('div');
        notification.className = 'notification success';
        notification.innerHTML = `
            <span>✅ Successfully subscribed to status updates!</span>
            <button onclick="this.parentElement.remove()">&times;</button>
        `;
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            if (notification.parentElement) {
                notification.remove();
            }
        }, 5000);
    }

    formatTime(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();
        const diff = now - date;
        const minutes = Math.floor(diff / 60000);
        
        if (minutes < 1) {
            return 'Just now';
        } else if (minutes < 60) {
            return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
        } else if (minutes < 1440) {
            const hours = Math.floor(minutes / 60);
            return `${hours} hour${hours > 1 ? 's' : ''} ago`;
        } else {
            return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
    }

    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    startAutoRefresh(interval = 60000) {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
        }
        this.refreshInterval = setInterval(() => this.loadStatus(), interval);
    }

    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }

    destroy() {
        this.stopAutoRefresh();
    }
}

export default StatusDashboard;