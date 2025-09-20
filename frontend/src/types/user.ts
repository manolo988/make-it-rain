export interface User {
  id: number;
  email: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  email: string;
  name: string;
  password: string;
}

export interface UpdateUserRequest {
  email?: string;
  name?: string;
  is_active?: boolean;
}

export interface UsersResponse {
  users: User[];
  total: number;
  page: number;
  page_size: number;
}