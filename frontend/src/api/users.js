import client from './client';

const USERS_ENDPOINT = '/api/v1/users';

export const usersApi = {
  getAll: async (page = 1, pageSize = 10) => {
    const response = await client.get(USERS_ENDPOINT, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getById: async (id) => {
    const response = await client.get(`${USERS_ENDPOINT}/${id}`);
    return response.data;
  },

  create: async (userData) => {
    const response = await client.post(USERS_ENDPOINT, userData);
    return response.data;
  },

  update: async (id, userData) => {
    const response = await client.put(`${USERS_ENDPOINT}/${id}`, userData);
    return response.data;
  },

  delete: async (id) => {
    const response = await client.delete(`${USERS_ENDPOINT}/${id}`);
    return response.data;
  },
};

export default usersApi;