

# Car Inventory System — Premium Automotive Frontend

## Design Theme
- **Dark mode by default** with an automotive luxury aesthetic
- **Glassmorphism** effects on cards and modals (backdrop-blur, semi-transparent backgrounds)
- **Inter font** for clean, premium typography
- **Smooth animations** (fade-in, scale, hover effects) on cards and page transitions
- **Color palette**: Deep blacks/charcoals with accent colors (e.g., electric blue or gold highlights)

## Pages & Features

### 1. Home Page (Public)
- **Hero section** with bold "Premium Inventory" headline, subtle gradient background, and a call-to-action
- **Car grid** displaying all vehicles fetched via GraphQL (`cars` query)
- Each **car card** shows: Make, Model, Year, Price (formatted as currency), Mileage, and Color accent
- Hover effect reveals a "View Details" overlay with smooth animation
- Loading skeletons while data is fetched

### 2. Login Page (Public)
- Centered glassmorphism card on a dark background
- **Two-step flow**:
  - Step 1: Email input → "Send Code" button (calls `requestLogin` mutation)
  - Step 2: 6-digit OTP input → "Verify" button (calls `verifyLogin` mutation)
- On success: stores JWT in localStorage, redirects to Dashboard
- Toast notifications for success/error states

### 3. Admin Dashboard (Protected — requires JWT)
- **Route guard**: redirects to Login if no valid token
- **Navbar** with app logo, "Dashboard" title, and Sign Out button
- **Data table** showing all cars with sortable columns
- **Row actions**: Edit and Delete buttons per car
- **"Add Car" button** opens a modal form (Make, Model, Year, Price, Color, Mileage)
- **Edit** opens the same modal pre-filled with existing data
- **Delete** shows a confirmation dialog before calling `deleteCar` mutation
- All mutations send JWT in `Authorization: Bearer <token>` header

## Architecture

### Auth Context (React Context API)
- Stores JWT token and login state
- Provides `login`, `logout` helpers
- Wraps the app to make auth state available everywhere

### GraphQL Client
- Utility function using `fetch` API to call `http://localhost:8000/graphql`
- Automatically attaches `Authorization` header only for protected mutations
- Vite proxy config (`/graphql` → `localhost:8000`) to avoid CORS issues

### Navigation
- React Router with routes: `/` (Home), `/login`, `/dashboard`
- Protected route wrapper for Dashboard

### UX Polish
- Loading spinners/skeletons for all data fetches
- Toast notifications (via Sonner) for all mutation results and errors
- Responsive design for desktop and tablet

