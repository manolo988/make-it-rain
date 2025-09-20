export enum AuctionStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  COMPLETED = 'completed',
  CANCELLED = 'cancelled'
}

export interface Auction {
  id: number;
  user_id: number;
  title: string;
  description: string;
  start_date: string;
  end_date: string;
  start_price: number;
  status: AuctionStatus;
  created_at: string;
  updated_at: string;
  items?: AuctionItem[];
}

export interface AuctionItem {
  id: number;
  auction_id: number;
  name: string;
  description: string;
  image_url: string;
  created_at: string;
  updated_at: string;
}

export interface CreateAuctionRequest {
  user_id?: number;
  title: string;
  description: string;
  start_date: string;
  end_date: string;
  start_price: number;
  status?: AuctionStatus;
}

export interface CreateAuctionItemRequest {
  auction_id: number;
  name: string;
  description: string;
  image_url: string;
}

export interface Bid {
  id: number;
  auction_item_id: number;
  user_id: number;
  amount: number;
  created_at: string;
  updated_at: string;
}

export interface CreateBidRequest {
  auction_item_id: number;
  user_id: number;
  amount: number;
}

export interface BackendPaginatedAuctions {
  auctions: Auction[];
  total_count: number;
  page: number;
  page_size: number;
}

export interface AuctionResponse {
  auction: Auction;
  items: AuctionItem[];
}