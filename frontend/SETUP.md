# O2C Graph Intelligence Dashboard - Setup Guide

## ✅ Complete Project Structure

Your full-stack O2C Graph Intelligence Dashboard is now ready!

```
FDE-DodgeAI/
├── backend/                   (Go API Server - Ready)
├── frontend/                  (React Dashboard - Ready)
└── sap-o2c-data/             (Test datasets)
```

## 🚀 Quick Start

### Step 1: Start the Go Backend

Open Terminal 1 (PowerShell):

```powershell
cd c:\Users\kingk\OneDrive\Documents\FDE-DodgeAI\backend
go run ./cmd
```

Expected output:
```
[TIMESTAMP] [INFO] [SERVER] Starting O2C Graph API server
[TIMESTAMP] [INFO] [DATABASE] ✓ Connected to PostgreSQL
[TIMESTAMP] [INFO] [INGESTION] ✓ Data ingestion complete
Server running on http://localhost:8089
```

### Step 2: Start the React Frontend

Open Terminal 2 (PowerShell):

```powershell
cd c:\Users\kingk\OneDrive\Documents\FDE-DodgeAI\frontend
npm install  # (First time only)
npm run dev
```

Expected output:
```
VITE v5.0.0  ready in 123 ms

➜  Local:   http://localhost:5173/
➜  Press h to show help
```

### Step 3: Open Browser

Navigate to: **http://localhost:5173**

---

## 📦 What's Included

### Frontend Components (TypeScript + React 18)

#### Layout
- **Header.tsx** - Logo, title, dark/light toggle, connection status
- **StatusBar.tsx** - Node/edge counts, dataset info, last updated time

#### Graph Visualization
- **GraphPanel.tsx** - React Flow canvas with full O2C graph
- **NodeDetailPanel.tsx** - Slide-in sidebar with node properties & neighbors
- **Node Components** (7 types):
  - SalesOrderNode.tsx
  - DeliveryNode.tsx
  - BillingNode.tsx
  - PaymentNode.tsx
  - CustomerNode.tsx
  - ProductNode.tsx
  - PlantNode.tsx

#### Chat Interface
- **ChatPanel.tsx** - Message history with example queries
- **ChatMessage.tsx** - User/AI messages with SQL blocks & data tables
- **ChatInput.tsx** - Textarea with auto-resize, character counter, Send button

#### Hooks (API Integration)
- **useGraph.ts** - Fetches all nodes & edges
- **useNodeDetail.ts** - Fetches single node with neighbors
- **useChat.ts** - Manages chat messages and Gemini calls
- **useLayoutedElements.ts** - Dagre layout algorithm

#### Styling
- **index.css** - CSS variables, dark/light theme, animations
- **tailwind.config.ts** - Custom theme configuration
- **package.json** - All dependencies

### API Integration

Proxy configured in vite.config.ts:
- React calls `/api/graph`, `/api/node/:type/:id`, `/api/chat`
- Vite dev server forwards to `http://localhost:8089`
- In production, update `VITE_API_URL` environment variable

---

## 🎨 Design Features

✅ **Dark-first enterprise theme** (light mode toggle ready)
✅ **Teal accent color** (#4f98a3) for CTAs and highlights
✅ **Node-type colors**: Blue (SO), Purple (Delivery), Amber (Billing), Green (Payment), Teal (Customer), Orange (Product), Gray (Plant)
✅ **Split-panel layout**: 60% graph, 40% chat
✅ **Responsive**: Desktop (side-by-side) → Tablet (stacked) → Mobile (tabs)
✅ **Smooth animations**: Slide-in panels, fade-in messages, pulsing loaders
✅ **Monospace fonts**: JetBrains Mono for IDs and amounts
✅ **Indian Rupee formatting**: ₹9,966.10

---

## 🔗 API Endpoints

Backend running on `http://localhost:8089` provides:

| Method | Endpoint | Response |
|--------|----------|----------|
| GET | `/api/graph` | All nodes & edges |
| GET | `/api/node/:type/:id` | Single node + neighbors + edges |
| POST | `/api/chat` | Natural language → SQL → Answer |

### Example Requests

```bash
# Get all graph data
curl http://localhost:8089/api/graph

# Get a sales order node
curl http://localhost:8089/api/node/SalesOrder/100000001

# Ask a question
curl -X POST http://localhost:8089/api/chat \
  -H "Content-Type: application/json" \
  -d '{"query": "Show me high-value orders"}'
```

---

## 📂 Frontend File Structure

```
frontend/src/
├── main.tsx                           # Entry point
├── App.tsx                            # Main layout
├── index.css                          # Global styles + theme
├── lib/
│   └── api.ts                        # Axios client & types
├── hooks/
│   ├── useGraph.ts
│   ├── useNodeDetail.ts
│   ├── useChat.ts
│   └── useLayoutedElements.ts
├── components/
│   ├── layout/
│   │   ├── Header.tsx
│   │   └── StatusBar.tsx
│   ├── graph/
│   │   ├── GraphPanel.tsx
│   │   ├── NodeDetailPanel.tsx
│   │   ├── Skeleton.tsx
│   │   └── nodes/
│   │       ├── NodeCard.tsx
│   │       ├── SalesOrderNode.tsx
│   │       ├── DeliveryNode.tsx
│   │       ├── BillingNode.tsx
│   │       ├── PaymentNode.tsx
│   │       ├── CustomerNode.tsx
│   │       ├── ProductNode.tsx
│   │       └── PlantNode.tsx
│   └── chat/
│       ├── ChatPanel.tsx
│       ├── ChatMessage.tsx
│       └── ChatInput.tsx
```

---

## 🔧 Key Technologies

- **React 18** - UI framework
- **Vite** - Fast build tool
- **React Flow** - Graph visualization
- **TanStack React Query** - API data fetching
- **Axios** - HTTP client
- **Tailwind CSS** - Utility-first styling
- **Lucide React** - Icon library
- **Dagre** - Hierarchical graph layout
- **TypeScript** - Type safety

---

## 🎯 Feature Highlights

### Graph Panel
- ✅ Hierarchical left-to-right (LR) layout via Dagre
- ✅ 7 node types with distinct colors
- ✅ Animated dashed edges between nodes
- ✅ Node legend (top-left)
- ✅ Minimap (bottom-left)
- ✅ Zoom controls (bottom-right)
- ✅ Fit view button
- ✅ Click node → highlight neighbors & show detail panel
- ✅ Click empty canvas → deselect node

### Node Detail Panel
- ✅ Slides in from right when node selected
- ✅ Shows all node properties (formatted dates/amounts)
- ✅ Lists connected nodes as clickable chips
- ✅ "View in Graph" button to fit view
- ✅ Close button (X)

### Chat Panel
- ✅ Empty state with 3 example queries
- ✅ User messages (right-aligned, teal)
- ✅ AI messages (left-aligned, dark)
- ✅ Collapsible SQL blocks
- ✅ Mini data table (if ≤20 rows)
- ✅ Node reference chips: Click to navigate
- ✅ Animated loading dots
- ✅ Auto-scroll to latest message
- ✅ Character counter (>200 chars)
- ✅ Keyboard: Enter to send, Shift+Enter for newline

---

## 🛠️ Build & Deployment

### Development

```bash
cd frontend
npm install
npm run dev  # Starts on http://localhost:5173
```

### Production Build

```bash
npm run build  # Outputs to dist/
npm run preview  # Preview production build locally
```

### Docker (Optional)

Frontend can be containerized:

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package.json .
RUN npm install
COPY src src
RUN npm run build
```

---

## 🧪 Testing the Full Stack

1. **Backend Health Check**:
   ```bash
   curl http://localhost:8089/health
   ```

2. **Graph Data**:
   Click "Fit View" in React app

3. **Node Details**:
   Click any node → Detail panel slides in

4. **Chat Examples**:
   - "Which products have the most billing documents?"
   - "Show me sales orders with missing deliveries"
   - "Trace the full flow of billing document 90504248"

5. **Logs**:
   - Backend: Console output (categorized)
   - Backend file: `backend/logs/o2c-graph-YYYYMMDD.log`
   - Frontend: Browser DevTools Console

---

## 🐛 Troubleshooting

### Frontend won't connect to backend

1. Check backend is running: `http://localhost:8089/health`
2. Check frontend dev server: `http://localhost:5173`
3. Check vite.config.ts proxy is correct
4. Clear browser cache: Ctrl+Shift+Delete

### Graph not loading

1. Check browser console for API errors
2. Check backend logs for database issues
3. Verify PostgreSQL is running
4. Verify SAP data in `sap-o2c-data/` directory

### Chat not working

1. Check Gemini API key in backend .env
2. Check backend `/api/chat` endpoint responds
3. Check Gemini API quota
4. Check browser console for errors

### Performance issues

1. Reduce graph size (backend auto-limits to reasonable node count)
2. Check browser DevTools Performance tab
3. Verify network/API response times
4. Clear cache and rebuild: `rm -rf node_modules && npm install`

---

## 📝 Notes

- **Dark mode is default** - Toggle via sun/moon icon in header
- **No authentication** - Add auth layer as needed
- **Indian Rupee formatting** - Adjust currency in code if needed
- **Responsive layout** - Stack vertically on screens <1200px

---

## 🎉 You're All Set!

Your O2C Graph Intelligence Dashboard is ready for:
- ✅ Interactive graph visualization
- ✅ Natural language queries via Gemini
- ✅ Real-time data exploration
- ✅ Enterprise analytics

Start both servers and navigate to **http://localhost:5173** 🚀
