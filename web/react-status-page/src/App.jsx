import { ThemeProvider } from './contexts/ThemeContext';
import ErrorBoundary from './components/common/ErrorBoundary';
import Header from './components/Layout/Header';
import Footer from './components/Layout/Footer';
import StatusDashboard from './components/Dashboard/StatusDashboard';
import IncidentManager from './components/Incidents/IncidentManager';

function App() {
  return (
    <ThemeProvider>
      <ErrorBoundary>
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
          <Header />
          <main className="flex-1">
            <StatusDashboard />
            <IncidentManager />
          </main>
          <Footer />
        </div>
      </ErrorBoundary>
    </ThemeProvider>
  );
}

export default App;