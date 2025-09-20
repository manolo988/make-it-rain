import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Plus, Search, Filter, RefreshCw } from 'lucide-react';
import { auctionService } from '../services/api';
import { AuctionCard } from './AuctionCard';
import { useNavigate } from 'react-router-dom';

export const AuctionList: React.FC = () => {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [searchTerm, setSearchTerm] = useState('');

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['auctions', page],
    queryFn: () => auctionService.getAuctions(page, 12),
  });

  const filteredAuctions = data?.auctions?.filter((auction) =>
    auction.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    auction.description.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  const totalPages = data ? Math.ceil(data.total_count / data.page_size) : 0;

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-900 mb-2">Failed to load auctions</h3>
          <button
            onClick={() => refetch()}
            className="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            <RefreshCw className="w-4 h-4 mr-2" />
            Try again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h1 className="text-4xl font-bold text-gray-900 bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                Live Auctions
              </h1>
              <p className="mt-2 text-gray-600">
                Discover and bid on exclusive items
              </p>
            </div>
            <button
              onClick={() => navigate('/auctions/new')}
              className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-medium rounded-lg shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200"
            >
              <Plus className="w-5 h-5 mr-2" />
              Create Auction
            </button>
          </div>

          <div className="mt-6 flex flex-col sm:flex-row gap-4">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search auctions..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 bg-white border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              />
            </div>
            <button className="inline-flex items-center px-4 py-3 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
              <Filter className="w-5 h-5 mr-2 text-gray-500" />
              <span className="text-gray-700">Filters</span>
            </button>
          </div>
        </div>

        {filteredAuctions.length === 0 ? (
          <div className="text-center py-16 bg-white rounded-xl shadow-sm">
            <div className="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
              <Search className="w-10 h-10 text-gray-400" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No auctions found</h3>
            <p className="text-gray-500">Try adjusting your search or create a new auction</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredAuctions.map((auction) => (
                <AuctionCard
                  key={auction.id}
                  auction={auction}
                  onClick={() => navigate(`/auctions/${auction.id}`)}
                />
              ))}
            </div>

            {data && totalPages > 1 && (
              <div className="mt-8 flex justify-center gap-2">
                <button
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="px-4 py-2 bg-white border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  Previous
                </button>
                <div className="flex items-center gap-2">
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    const pageNum = i + 1;
                    return (
                      <button
                        key={pageNum}
                        onClick={() => setPage(pageNum)}
                        className={`px-4 py-2 rounded-lg transition-colors ${
                          page === pageNum
                            ? 'bg-indigo-600 text-white'
                            : 'bg-white border border-gray-300 hover:bg-gray-50'
                        }`}
                      >
                        {pageNum}
                      </button>
                    );
                  })}
                </div>
                <button
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                  className="px-4 py-2 bg-white border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  Next
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};