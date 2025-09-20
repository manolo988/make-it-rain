import { User, CreateUserRequest, UpdateUserRequest, UsersResponse } from '../types/user';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

class UserService {
  async getUsers(page = 1, pageSize = 10, sortBy = 'created_at', sortOrder = 'desc'): Promise<UsersResponse> {
    const response = await fetch(
      `${API_URL}/users?page=${page}&page_size=${pageSize}&sort_by=${sortBy}&sort_order=${sortOrder}`
    );
    if (!response.ok) throw new Error('Failed to fetch users');
    return response.json();
  }

  async getUser(id: number): Promise<User> {
    const response = await fetch(`${API_URL}/users/${id}`);
    if (!response.ok) {
      if (response.status === 404) throw new Error('User not found');
      throw new Error('Failed to fetch user');
    }
    return response.json();
  }

  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await fetch(`${API_URL}/users`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create user');
    }
    return response.json();
  }

  async updateUser(id: number, data: UpdateUserRequest): Promise<void> {
    const response = await fetch(`${API_URL}/users/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      if (response.status === 404) throw new Error('User not found');
      throw new Error('Failed to update user');
    }
  }

  async deleteUser(id: number): Promise<void> {
    const response = await fetch(`${API_URL}/users/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      if (response.status === 404) throw new Error('User not found');
      throw new Error('Failed to delete user');
    }
  }
}

export default new UserService();