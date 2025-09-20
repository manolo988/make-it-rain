import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useMutation, useQuery } from '@tanstack/react-query';
import { ArrowLeft, Save, Package, Image } from 'lucide-react';
import { auctionService } from '../services/api';
import { CreateAuctionItemRequest } from '../types/auction';

export const CreateAuctionItem: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [formData, setFormData] = useState<Omit<CreateAuctionItemRequest, 'auction_id'>>({
    name: '',
    description: '',
    image_url: '',
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const { data: auction } = useQuery({
    queryKey: ['auction', id],
    queryFn: () => auctionService.getAuction(parseInt(id!)),
    enabled: !!id,
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateAuctionItemRequest) => auctionService.createAuctionItem(data),
    onSuccess: () => {
      navigate(`/auctions/${id}`);
    },
    onError: (error: any) => {
      console.error('Failed to create auction item:', error);
    },
  });

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Item name is required';
    }

    if (!formData.description.trim()) {
      newErrors.description = 'Description is required';
    }

    if (!formData.image_url.trim()) {
      newErrors.image_url = 'Image URL is required';
    } else if (!isValidUrl(formData.image_url)) {
      newErrors.image_url = 'Please enter a valid URL';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const isValidUrl = (string: string) => {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm() && id) {
      const submitData: CreateAuctionItemRequest = {
        ...formData,
        auction_id: parseInt(id),
      };
      createMutation.mutate(submitData);
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
          onClick={() => navigate(`/auctions/${id}`)}
          className="mb-6 inline-flex items-center text-gray-600 hover:text-gray-900 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 mr-2" />
          Back to auction
        </button>

        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <div className="bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-8">
            <h1 className="text-3xl font-bold text-white">Add Auction Item</h1>
            {auction?.auction && (
              <p className="mt-2 text-indigo-100">
                Adding to: {auction.auction.title}
              </p>
            )}
          </div>

          <form onSubmit={handleSubmit} className="px-8 py-8 space-y-6">
            <div className="bg-indigo-50 border border-indigo-200 rounded-lg p-4">
              <div className="flex items-start space-x-3">
                <Package className="w-5 h-5 text-indigo-600 mt-0.5" />
                <div className="flex-1">
                  <p className="text-sm text-indigo-800 font-medium">Item Details</p>
                  <p className="text-sm text-indigo-700 mt-1">
                    Each item in the auction can be bid on independently. Set a clear starting price
                    and minimum bid increment to ensure smooth bidding.
                  </p>
                </div>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Item Name
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                placeholder="e.g., Vintage Sony Walkman"
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                  errors.name ? 'border-red-300' : 'border-gray-300'
                }`}
              />
              {errors.name && (
                <p className="mt-1 text-sm text-red-600">{errors.name}</p>
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
                placeholder="Describe the item's condition, features, and any special characteristics..."
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                  errors.description ? 'border-red-300' : 'border-gray-300'
                }`}
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <Image className="w-4 h-4 inline mr-1" />
                Image URL
              </label>
              <input
                type="url"
                name="image_url"
                value={formData.image_url}
                onChange={handleChange}
                placeholder="https://example.com/image.jpg"
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors ${
                  errors.image_url ? 'border-red-300' : 'border-gray-300'
                }`}
              />
              {errors.image_url && (
                <p className="mt-1 text-sm text-red-600">{errors.image_url}</p>
              )}
              {formData.image_url && !errors.image_url && (
                <div className="mt-4 p-4 bg-gray-50 rounded-lg">
                  <p className="text-sm text-gray-600 mb-2">Preview:</p>
                  <img
                    src={formData.image_url}
                    alt="Item preview"
                    className="max-w-full h-48 object-cover rounded"
                    onError={(e) => {
                      (e.target as HTMLImageElement).style.display = 'none';
                    }}
                  />
                </div>
              )}
            </div>

            <div className="flex justify-end gap-4 pt-6 border-t">
              <button
                type="button"
                onClick={() => navigate(`/auctions/${id}`)}
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
                    Adding Item...
                  </>
                ) : (
                  <>
                    <Save className="w-5 h-5 mr-2" />
                    Add Item
                  </>
                )}
              </button>
            </div>

            {createMutation.isError && (
              <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
                <p className="text-sm text-red-800">
                  Failed to add item. Please try again.
                </p>
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
};