import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { format } from 'date-fns';
import { ArrowLeft, Clock, Calendar, Plus, DollarSign, Package, Info } from 'lucide-react';
import { auctionService } from '../services/api';
import clsx from 'clsx';

export const AuctionDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [showAddItem, setShowAddItem] = useState(false);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['auction', id],
    queryFn: () => auctionService.getAuction(parseInt(id!)),
    enabled: !!id,
  });

  const auction = data?.auction;
  const items = data?.items || [];

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  if (isError || !auction) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Auction not found</h3>
          <button
            onClick={() => navigate('/auctions')}
            className="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to auctions
          </button>
        </div>
      </div>
    );
  }

  const startTime = new Date(auction.start_date);
  const endTime = new Date(auction.end_date);

  const getStatusColor = () => {
    const statusLower = auction.status?.toLowerCase();
    switch (statusLower) {
      case 'active':
        return 'bg-green-100 text-green-800 border-green-200';
      case 'pending':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'completed':
        return 'bg-gray-100 text-gray-800 border-gray-200';
      case 'cancelled':
        return 'bg-red-100 text-red-800 border-red-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <button
          onClick={() => navigate('/auctions')}
          className="mb-6 inline-flex items-center text-gray-600 hover:text-gray-900 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 mr-2" />
          Back to auctions
        </button>

        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <div className="bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-12">
            <div className="flex justify-between items-start">
              <div>
                <h1 className="text-4xl font-bold text-white mb-3">
                  {auction.title}
                </h1>
                <p className="text-indigo-100 text-lg max-w-3xl">
                  {auction.description}
                </p>
              </div>
              <span className={clsx(
                'px-4 py-2 rounded-full text-sm font-medium border',
                getStatusColor()
              )}>
                {auction.status.charAt(0).toUpperCase() + auction.status.slice(1)}
              </span>
            </div>
          </div>

          <div className="px-8 py-6 border-b border-gray-200">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="flex items-center space-x-3">
                <div className="p-3 bg-blue-100 rounded-lg">
                  <Calendar className="w-6 h-6 text-blue-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Start Date</p>
                  <p className="font-medium text-gray-900">
                    {format(startTime, 'MMM d, yyyy h:mm a')}
                  </p>
                </div>
              </div>

              <div className="flex items-center space-x-3">
                <div className="p-3 bg-red-100 rounded-lg">
                  <Clock className="w-6 h-6 text-red-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">End Date</p>
                  <p className="font-medium text-gray-900">
                    {format(endTime, 'MMM d, yyyy h:mm a')}
                  </p>
                </div>
              </div>

              <div className="flex items-center space-x-3">
                <div className="p-3 bg-green-100 rounded-lg">
                  <Package className="w-6 h-6 text-green-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-500">Total Items</p>
                  <p className="font-medium text-gray-900">
                    {items.length} items
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="px-8 py-6">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold text-gray-900">Auction Items</h2>
              <button
                onClick={() => setShowAddItem(true)}
                className="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                <Plus className="w-4 h-4 mr-2" />
                Add Item
              </button>
            </div>

            {items.length === 0 ? (
              <div className="text-center py-12 bg-gray-50 rounded-lg">
                <Package className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500 mb-4">No items in this auction yet</p>
                <button
                  onClick={() => setShowAddItem(true)}
                  className="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Add First Item
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {items.map((item) => (
                  <div key={item.id} className="bg-gray-50 rounded-lg p-6 hover:shadow-md transition-shadow">
                    {item.image_url && (
                      <div className="w-full h-48 bg-gray-200 rounded-lg mb-4 overflow-hidden">
                        <img
                          src={item.image_url}
                          alt={item.name}
                          className="w-full h-full object-cover"
                        />
                      </div>
                    )}
                    <h3 className="font-semibold text-gray-900 mb-2">{item.name}</h3>
                    <p className="text-sm text-gray-600 mb-4 line-clamp-2">{item.description}</p>
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-500">Starting Price</span>
                        <span className="font-medium text-gray-900">
                          ${item.starting_price.toFixed(2)}
                        </span>
                      </div>
                      {item.current_price && (
                        <div className="flex justify-between items-center">
                          <span className="text-sm text-gray-500">Current Bid</span>
                          <span className="font-bold text-green-600">
                            ${item.current_price.toFixed(2)}
                          </span>
                        </div>
                      )}
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-500">Min Increment</span>
                        <span className="text-sm text-gray-700">
                          ${item.minimum_bid_increment.toFixed(2)}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {showAddItem && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl p-6 max-w-md w-full">
              <h3 className="text-xl font-bold mb-4">Add Item to Auction</h3>
              <p className="text-gray-600 mb-4">
                Item creation form would go here
              </p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => setShowAddItem(false)}
                  className="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={() => navigate(`/auctions/${id}/add-item`)}
                  className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
                >
                  Continue
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};