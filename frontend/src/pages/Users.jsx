import UserList from '../components/UserList';

const Users = () => {
  return (
    <div className="px-4 sm:px-0">
      <h1 className="text-3xl font-bold mb-6 bg-gradient-rain bg-clip-text text-transparent animate-slide-up">
        User Management
      </h1>
      <UserList />
    </div>
  );
};

export default Users;