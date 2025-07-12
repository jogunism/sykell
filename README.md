# Sykell Test task - Hyunwoo Cho

## Tech Stack

- **Backend**: Go
  - **Framework**: Gin
  - **Database**: MySQL
  - **Web Scraping**: `goquery`

- **Frontend**: React
  - **UI Framework**: Tailwind CSS
  - **HTTP Client**: `axios`

- **Containerization**: Docker, Docker Compose

## Getting Started

### Prerequisites

- Node 20.x above
- Docker

### Running the Application
   ```bash
   git clone https://github.com/jogunism/sykell
   cd sykell

   #Project start
   ./start.sh

   #Project stop
   ./stop.sh

   #logs
   docker logs -f sykell-frontend-1
   docker logs -f sykell-backend-1
   ```

3. **Access the application**:
   - **Frontend**: `http://localhost:5173`
   - **Backend API**: `http://localhost:8080`

## Project Structure

```
sykell/
├── backend/                # Go Backend Application
│   ├── application/        # Application layer (services, commands, queries)
│   ├── domain/             # Core domain models and business logic
│   ├── infrastructure/     # Persistence, database connections
│   ├── presentation/       # API handlers and routing
│   ├── Dockerfile
│   └── main.go             # Application entry point
│
├── frontend/               # React Frontend Application
│   ├── src/                # source code
│   │   ├── components/     # Reusable UI components
│   │   └── App.tsx         # Main application component
│   ├── Dockerfile
│   └── package.json
│
├── docker-compose.yml      # Docker Compose for development
├── docker-compose.prod.yml # Docker Compose for production
└── README.md
```