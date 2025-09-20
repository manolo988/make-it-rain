import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Users from './pages/Users';
import NotFound from './pages/NotFound';
import { AuctionList } from './components/AuctionList';
import { AuctionDetail } from './components/AuctionDetail';
import { CreateAuction } from './components/CreateAuction';
import { CreateAuctionItem } from './components/CreateAuctionItem';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/auctions" replace />} />
            <Route path="auctions" element={<AuctionList />} />
            <Route path="auctions/new" element={<CreateAuction />} />
            <Route path="auctions/:id" element={<AuctionDetail />} />
            <Route path="auctions/:id/add-item" element={<CreateAuctionItem />} />
            <Route path="users" element={<Users />} />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}

export default App
