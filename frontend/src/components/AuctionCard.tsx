import React from 'react';
import { Clock, Calendar, Tag, TrendingUp } from 'lucide-react';
import { format, formatDistanceToNow, isPast, isFuture } from 'date-fns';
import { Auction } from '../types/auction';
import clsx from 'clsx';

interface AuctionCardProps {
  auction: Auction;
  onClick: () => void;
}

export const AuctionCard: React.FC<AuctionCardProps> = ({ auction, onClick }) => {
  const startTime = new Date(auction.start_date);
  const endTime = new Date(auction.end_date);

  const getStatusBadge = () => {
    const statusLower = auction.status?.toLowerCase();

    if (statusLower === 'cancelled') {
      return (
        <span className="px-2 py-1 text-xs font-medium bg-red-100 text-red-700 rounded-full">
          Cancelled
        </span>
      );
    }

    if (statusLower === 'completed' || isPast(endTime)) {
      return (
        <span className="px-2 py-1 text-xs font-medium bg-gray-100 text-gray-700 rounded-full">
          Completed
        </span>
      );
    }

    if (statusLower === 'pending' || isFuture(startTime)) {
      return (
        <span className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-700 rounded-full">
          Upcoming
        </span>
      );
    }

    return (
      <span className="px-2 py-1 text-xs font-medium bg-green-100 text-green-700 rounded-full animate-pulse">
        Live Now
      </span>
    );
  };

  const getTimeDisplay = () => {
    if (isPast(endTime)) {
      return `Ended ${formatDistanceToNow(endTime, { addSuffix: true })}`;
    }

    if (isFuture(startTime)) {
      return `Starts ${formatDistanceToNow(startTime, { addSuffix: true })}`;
    }

    return `Ends ${formatDistanceToNow(endTime, { addSuffix: true })}`;
  };

  return (
    <div
      onClick={onClick}
      className={clsx(
        'group relative bg-white rounded-xl shadow-sm border border-gray-200 p-6',
        'hover:shadow-xl hover:border-indigo-500 transition-all duration-300 cursor-pointer',
        'hover:-translate-y-1'
      )}
    >
      <div className="absolute inset-0 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-xl opacity-0 group-hover:opacity-5 transition-opacity" />

      <div className="flex justify-between items-start mb-4">
        <div className="flex-1">
          <h3 className="text-xl font-bold text-gray-900 group-hover:text-indigo-600 transition-colors">
            {auction.title}
          </h3>
          <p className="mt-2 text-sm text-gray-600 line-clamp-2">
            {auction.description}
          </p>
        </div>
        <div className="ml-4">
          {getStatusBadge()}
        </div>
      </div>

      <div className="space-y-3 mt-6">
        <div className="flex items-center text-sm text-gray-500">
          <Calendar className="w-4 h-4 mr-2 text-gray-400" />
          <span>
            {format(startTime, 'MMM d, yyyy')} - {format(endTime, 'MMM d, yyyy')}
          </span>
        </div>

        <div className="flex items-center text-sm font-medium">
          <Clock className="w-4 h-4 mr-2 text-indigo-500" />
          <span className={clsx(
            isPast(endTime) ? 'text-gray-500' :
            isFuture(startTime) ? 'text-blue-600' : 'text-green-600'
          )}>
            {getTimeDisplay()}
          </span>
        </div>

        <div className="flex items-center text-sm text-gray-600">
          <Tag className="w-4 h-4 mr-2 text-gray-400" />
          <span>Starting at ${auction.start_price?.toFixed(2) || '0.00'}</span>
        </div>
      </div>

      <div className="mt-6 pt-4 border-t border-gray-100">
        <div className="flex items-center justify-between">
          <div className="flex items-center text-sm text-gray-500">
            <TrendingUp className="w-4 h-4 mr-1 text-green-500" />
            <span>View details</span>
          </div>
          <svg
            className="w-5 h-5 text-gray-400 group-hover:text-indigo-600 group-hover:translate-x-1 transition-all"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </div>
  );
};