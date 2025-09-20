# React Frontend - Interview Boilerplate

## Quick Start

```bash
# Install dependencies
make frontend-install

# Start development server (port 3000)
make frontend-dev

# Build for production
make frontend-build

# Run both backend and frontend together
# Terminal 1:
make air

# Terminal 2:
make frontend-dev
```

## Project Structure

```
frontend/
├── src/
│   ├── api/           # API client and endpoints
│   │   ├── client.js    # Axios instance with interceptors
│   │   ├── users.js     # User CRUD operations
│   │   └── health.js    # Health check endpoints
│   ├── components/    # Reusable UI components
│   │   ├── Layout.jsx      # Main app layout with navigation
│   │   ├── UserList.jsx    # User management table
│   │   ├── UserForm.jsx    # User create/edit form
│   │   ├── HealthStatus.jsx # API health indicators
│   │   ├── LoadingSpinner.jsx
│   │   └── ErrorMessage.jsx
│   ├── pages/         # Page components
│   │   ├── Dashboard.jsx   # Main dashboard
│   │   ├── Users.jsx       # User management page
│   │   └── NotFound.jsx    # 404 page
│   ├── hooks/         # Custom React hooks (ready for additions)
│   ├── utils/         # Helper functions (ready for additions)
│   ├── App.jsx        # Main app with routing
│   ├── main.jsx       # Entry point
│   └── index.css      # Tailwind CSS imports
├── public/            # Static assets
├── package.json       # Dependencies and scripts
├── vite.config.js     # Vite config with API proxy
├── tailwind.config.js # Tailwind configuration
└── postcss.config.js  # PostCSS configuration
```

## Key Features for Interviews

### 1. Pre-configured API Client
- Axios with interceptors for auth tokens
- Centralized error handling
- Environment-based configuration

### 2. Component Library
- UserList with CRUD operations
- Form with validation
- Loading and error states
- Pagination support

### 3. Tailwind CSS
- Rapid UI development
- No custom CSS needed
- Responsive by default

### 4. Development Proxy
- No CORS issues in development
- Configured in `vite.config.js`
- Proxies `/api`, `/health`, `/ready` to backend

## Adding New Features

### Quick Resource Implementation (e.g., Products)

1. **Create API module** (`src/api/products.js`):
```javascript
import client from './client';

const PRODUCTS_ENDPOINT = '/api/v1/products';

export const productsApi = {
  getAll: async (page = 1, pageSize = 10) => {
    const response = await client.get(PRODUCTS_ENDPOINT, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },
  // ... other CRUD methods
};
```

2. **Copy UserList component** → ProductList
3. **Create Products page** using the new component
4. **Add route** in App.jsx
5. **Add navigation** in Layout.jsx

### Common Patterns

#### API Calls with Error Handling
```javascript
const [loading, setLoading] = useState(false);
const [error, setError] = useState(null);

const fetchData = async () => {
  setLoading(true);
  setError(null);
  try {
    const data = await api.getSomething();
    // handle data
  } catch (err) {
    setError('Failed to load');
  } finally {
    setLoading(false);
  }
};
```

#### Form Validation
```javascript
const validate = () => {
  const errors = {};
  if (!formData.name) errors.name = 'Name is required';
  if (!formData.email) errors.email = 'Email is required';
  return errors;
};
```

## Production Deployment

### Build Process
```bash
# Build frontend
make frontend-build

# Creates frontend/dist directory with:
# - Optimized JS/CSS bundles
# - index.html
# - Static assets
```

### Integration with Go Backend
- Backend serves `frontend/dist` in production
- SPA routing handled in `routes/routes.go`
- Single binary deployment possible

### Environment Variables
- Development: Uses Vite proxy
- Production: Set `VITE_API_URL` if needed
- Default: Uses same origin (no CORS)

## Common Commands

```bash
# Development
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build

# With Make
make frontend-install   # Install dependencies
make frontend-dev       # Start dev server
make frontend-build     # Build for production
make frontend-clean     # Remove build & node_modules
```

## Interview Tips

1. **Quick Setup**: Everything works out of the box
2. **Copy Components**: UserList/Form are templates for other resources
3. **Use Tailwind Classes**: No custom CSS needed
4. **API Client Ready**: Just add new endpoints to `src/api/`
5. **Error Handling**: Patterns already established
6. **Loading States**: Components ready to use

## Technologies

- **React 18** - Latest React features
- **Vite** - Fast build tool with HMR
- **React Router v6** - Client-side routing
- **Axios** - HTTP client with interceptors
- **Tailwind CSS** - Utility-first CSS framework

## Notes for Interview

- Frontend runs on port 3000, backend on 8080
- API proxy configured - no CORS issues
- Hot reload enabled for rapid development
- Component structure mirrors backend pattern
- Ready for additional features with minimal setup
