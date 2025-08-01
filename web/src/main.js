import StatusAPI from './services/api.js';
import StatusDashboard from './components/StatusDashboard.js';
import IncidentManager from './components/IncidentManager.js';

document.addEventListener('DOMContentLoaded', () => {
    const api = new StatusAPI();
    const dashboard = new StatusDashboard('status-container', api);
    const incidentManager = new IncidentManager();
    
    window.statusDashboard = dashboard;
    window.incidentManager = incidentManager;
    
    dashboard.init();
    incidentManager.init();

    window.addEventListener('beforeunload', () => {
        dashboard.destroy();
    });

    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.register('/sw.js').catch(err => {
            console.log('Service Worker registration failed:', err);
        });
    }
});