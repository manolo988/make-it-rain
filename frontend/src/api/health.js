import client from './client';

export const healthApi = {
  check: async () => {
    const response = await client.get('/health');
    return response.data;
  },

  ready: async () => {
    const response = await client.get('/ready');
    return response.data;
  },
};

export default healthApi;