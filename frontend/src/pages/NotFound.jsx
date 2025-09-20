import { Link } from 'react-router-dom';

const NotFound = () => {
  return (
    <div className="flex flex-col items-center justify-center py-20">
      <h1 className="text-6xl font-bold bg-gradient-rain bg-clip-text text-transparent animate-fade-in">
        404
      </h1>
      <p className="mt-4 text-xl text-gray-400">Page not found</p>
      <Link
        to="/"
        className="mt-6 btn-primary"
      >
        Go back to Dashboard
      </Link>
    </div>
  );
};

export default NotFound;