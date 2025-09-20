import axios from 'axios';
import type { Auction, AuctionItem, CreateAuctionRequest, CreateAuctionItemRequest, BackendPaginatedAuctions, AuctionResponse } from '../types/auction';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const auctionService = {
  async getAuctions(page = 1, pageSize = 10): Promise<BackendPaginatedAuctions> {
    const { data } = await api.get<BackendPaginatedAuctions>('/auctions', {
      params: { page, page_size: pageSize },
    });
    return data;
  },

  async getAuction(id: number): Promise<AuctionResponse> {
    const { data } = await api.get<AuctionResponse>(`/auctions/${id}`);
    return data;
  },

  async createAuction(auction: CreateAuctionRequest): Promise<Auction> {
    const { data } = await api.post<Auction>('/auctions', auction);
    return data;
  },

  async createAuctionItem(item: CreateAuctionItemRequest): Promise<AuctionItem> {
    const { data } = await api.post<AuctionItem>('/auctions/items', item);
    return data;
  },
};

export default api;