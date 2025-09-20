import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { ArrowLeft, Save, Calendar, Clock, Info } from 'lucide-react';
import { auctionService } from '../services/api';
import { CreateAuctionRequest, AuctionStatus } from '../types/auction';

export const CreateAuction: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState<CreateAuctionRequest>({
    title: '',
    description: '',
    start_date: '',
    end_date: '',
    start_price: 0,
    user_id: 1,
    status: AuctionStatus.ACTIVE
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const createMutation = useMutation({
    mutationFn: auctionService.createAuction,
    onSuccess: (data) => {
      navigate(`/auctions/${data.id}`);
    },
    onError: (error: any) => {
      console.error('Failed to create auction:', error);
    },
  });

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.title.trim()) {
      newErrors.title = 'Title is required';
    }

    if (!formData.description.trim()) {
      newErrors.description = 'Description is required';
    }

    if (!formData.start_date) {
      newErrors.start_date = 'Start date is required';
    }

    if (!formData.end_date) {
      newErrors.end_date = 'End date is required';
    } else if (formData.start_date && new Date(formData.end_date) <= new Date(formData.start_date)) {
      newErrors.end_date = 'End date must be after start date';
    }

    if (formData.start_price < 0) {
      newErrors.start_price = 'Starting price must be positive';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm()) {
      const formattedData = {
        ...formData,
        start_date: new Date(formData.start_date).toISOString(),
        end_date: new Date(formData.end_date).toISOString(),
      };
      createMutation.mutate(formattedData);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    const processedValue = type === 'number' ? parseFloat(value) || 0 : value;
    setFormData(prev => ({ ...prev, [name]: processedValue }));
    if (errors[name]) {
      setErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <button
          onClick={() => navigate('/auctions')}
          className="mb-6 inline-flex items-center text-gray-600 hover:text-gray-900 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 mr-2" />
          Back to auctions
        </button>

        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <div className="bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-8">
            <h1 className="text-3xl font-bold text-white">Create New Auction</h1>
            <p className="mt-2 text-indigo-100">Set up your auction details and timing</p>
          </div>

          <form onSubmit={handleSubmit} className="px-8 py-8 space-y-6">
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 flex items-start space-x-3">
              <Info className="w-5 h-5 text-blue-600 mt-0.5" />
              <div className="flex-1">
                <p className="text-sm text-blue-800 font-medium">Before you begin</p>
                <p className="text-sm text-blue-700 mt-1">
                  Once an auction is created, you can add items to it. Make sure to set appropriate
                  start and end times as they cannot be changed once the auction goes live.
                </p>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Auction Title
              </label>
              <input
                type="text"
                name="title"
                value={formData.title}
                onChange={handleChange}
                placeholder="e.g., Vintage Electronics Collection"
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                  errors.title ? 'border-red-300' : 'border-gray-300'
                }`}
              />
              {errors.title && (
                <p className="mt-1 text-sm text-red-600">{errors.title}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Description
              </label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleChange}
                rows={4}
                placeholder="Provide a detailed description of your auction..."
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                  errors.description ? 'border-red-300' : 'border-gray-300'
                }`}
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description}</p>
              )}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Starting Price
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                  <input
                    type="number"
                    name="start_price"
                    value={formData.start_price}
                    onChange={handleChange}
                    step="0.01"
                    min="0"
                    placeholder="0.00"
                    className={`w-full pl-8 pr-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                      errors.start_price ? 'border-red-300' : 'border-gray-300'
                    }`}
                  />
                </div>
                {errors.start_price && (
                  <p className="mt-1 text-sm text-red-600">{errors.start_price}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Status
                </label>
                <select
                  name="status"
                  value={formData.status || AuctionStatus.ACTIVE}
                  onChange={(e) => setFormData(prev => ({ ...prev, status: e.target.value as AuctionStatus }))}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                >
                  <option value={AuctionStatus.ACTIVE}>Active</option>
                  <option value={AuctionStatus.INACTIVE}>Inactive</option>
                  <option value={AuctionStatus.COMPLETED}>Completed</option>
                  <option value={AuctionStatus.CANCELLED}>Cancelled</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  <Calendar className="w-4 h-4 inline mr-1" />
                  Start Date & Time
                </label>
                <input
                  type="datetime-local"
                  name="start_date"
                  value={formData.start_date}
                  onChange={handleChange}
                  min={new Date().toISOString().slice(0, 16)}
                  className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                    errors.start_date ? 'border-red-300' : 'border-gray-300'
                  }`}
                />
                {errors.start_date && (
                  <p className="mt-1 text-sm text-red-600">{errors.start_date}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  <Clock className="w-4 h-4 inline mr-1" />
                  End Date & Time
                </label>
                <input
                  type="datetime-local"
                  name="end_date"
                  value={formData.end_date}
                  onChange={handleChange}
                  min={formData.start_date || new Date().toISOString().slice(0, 16)}
                  className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                    errors.end_date ? 'border-red-300' : 'border-gray-300'
                  }`}
                />
                {errors.end_date && (
                  <p className="mt-1 text-sm text-red-600">{errors.end_date}</p>
                )}
              </div>
            </div>

            <div className="flex justify-end gap-4 pt-6 border-t">
              <button
                type="button"
                onClick={() => navigate('/auctions')}
                className="px-6 py-3 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={createMutation.isPending}
                className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-medium rounded-lg shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {createMutation.isPending ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Creating...
                  </>
                ) : (
                  <>
                    <Save className="w-5 h-5 mr-2" />
                    Create Auction
                  </>
                )}
              </button>
            </div>

            {createMutation.isError && (
              <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
                <p className="text-sm text-red-800">
                  Failed to create auction. Please try again.
                </p>
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
};