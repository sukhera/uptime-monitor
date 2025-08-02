import ThemeToggle from './ThemeToggle';

const Header = () => {
  return (
    <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <h1 className="text-xl font-bold text-gray-900 dark:text-white">
              Status Page
            </h1>
          </div>
          
          <div className="flex items-center space-x-4">
            <ThemeToggle />
            <button className="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300">
              Subscribe to Updates
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;