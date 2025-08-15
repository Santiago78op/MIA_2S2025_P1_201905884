# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **MIA (Sistemas Operativos 2) Project 1 - EXT2 File System Simulator** developed as a university assignment. The project simulates an EXT2 file system with disk management capabilities using a Go backend and React TypeScript frontend.

**⚠️ Important:** This is an educational file system simulator that creates virtual disk files (`.mia`) and simulates disk operations - it does not modify actual system disks.

## Development Commands

### Backend (Go)
```bash
# Navigate to backend directory
cd backend

# Install dependencies
go mod tidy

# Run in development mode
go run main.go

# Build executable
go build -o backend main.go

# Run tests
go test ./...
```

### Frontend (React + TypeScript)
```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm start

# Build for production
npm run build

# Run tests
npm test
```

### Full Application
```bash
# Run both backend and frontend simultaneously
./run-all.sh

# Run backend only
./run-backend.sh

# Run frontend only  
./run-frontend.sh
```

### Docker
```bash
# Run with docker-compose
docker-compose up

# Production build
docker-compose --profile production up
```

## Architecture Overview

### Project Structure
```
MIA_2S2025_P1_201905884/
├── backend/                    # Go backend server
│   ├── main.go                # Main server with REST API and WebSocket
│   ├── command/               # Command parsing and execution
│   │   └── disk/              # Disk management commands (mkdisk, fdisk, etc.)
│   ├── struct/                # Data structures (MBR, Partition, EBR)
│   ├── Utils/                 # Logging and utilities
│   ├── action/                # Disk operations and file management
│   └── Discos/                # Virtual disk storage directory
├── frontend/                  # React TypeScript frontend
│   └── src/
│       ├── components/        # React components (CommandExecutor, FileSystemList)
│       ├── services/          # API client and command parser
│       └── hooks/             # Custom React hooks
└── docs/                      # Project documentation
```

### Core Components

**Backend (Go 1.21+):**
- **HTTP Server:** Gorilla Mux router with CORS support
- **WebSocket Server:** Real-time logging using Gorilla WebSocket
- **Command System:** Modular command parser supporting MIA script execution
- **Disk Simulation:** Binary file operations to simulate disk I/O
- **Data Structures:** MBR, Partition, EBR structures for disk organization

**Frontend (React 18 + TypeScript):**
- **Command Interface:** Text area for entering MIA commands
- **Real-time Logs:** WebSocket connection for live command output
- **File System Browser:** Lists virtual disks and mounted partitions
- **Script Execution:** Support for loading and executing `.smia` script files

### Key Features
- **Real-time Communication:** WebSocket + Server-Sent Events for live logging
- **Command Execution:** Support for MIA command language (mkdisk, fdisk, mount, etc.)
- **Virtual Disk Management:** Create and manage `.mia` virtual disk files
- **Partition Management:** MBR-based partitioning with primary/extended/logical support
- **EXT2 Simulation:** File system structures and operations

## API Endpoints

The backend exposes REST API endpoints on `localhost:8080`:

- `GET /api/health` - Server health check
- `GET /api/filesystems` - List virtual disks and file systems
- `POST /api/execute` - Execute MIA commands or scripts
- `POST /api/validate` - Validate command syntax
- `GET /api/commands` - Get list of supported commands
- `WebSocket /api/ws` - Real-time logging connection

## MIA Command System

### Currently Implemented Commands
- **mkdisk:** Create virtual disk files with MBR structure
  ```bash
  mkdisk -size 10 -unit M -fit FF -path /path/to/disk.mia
  ```

### Commands in Development
- **rmdisk:** Remove virtual disk files
- **fdisk:** Partition management (create/delete partitions)
- **mount:** Mount partitions for file system operations
- **mkfs:** Format partitions with EXT2 file system

### Command Syntax
Commands follow MIA specification with dash-prefixed parameters:
- Parameters: `-size`, `-path`, `-unit`, `-fit`, `-name`, etc.
- Units: `K` (Kilobytes), `M` (Megabytes)
- Fit algorithms: `BF` (Best Fit), `FF` (First Fit), `WF` (Worst Fit)

## Data Structures

### MBR (Master Boot Record)
```go
type MBR struct {
    MbrTamanio       int64        // Total disk size in bytes
    MbrFechaCreacion int64        // Creation timestamp
    MbrDiskSignature int64        // Unique disk identifier
    MbrFit           byte         // Partition fit algorithm (B/F/W)
    MbrParticiones   [4]Partition // Partition table (max 4)
}
```

### Partition Structure
```go
type Partition struct {
    PartStatus byte      // Status: 0=inactive, 1=active
    PartType   byte      // Type: P=Primary, E=Extended, L=Logical
    PartFit    byte      // Fit algorithm
    PartStart  int64     // Starting byte position
    PartSize   int64     // Size in bytes
    PartName   [16]byte  // Partition name
}
```

## Development Guidelines

### Code Style
- **Go:** Follow standard Go formatting with `gofmt`
- **TypeScript:** Use ESLint configuration in `frontend/.eslintrc`
- **Error Handling:** Always handle errors appropriately with logging
- **Logging:** Use the centralized logging system in `backend/Utils/logger.go`

### File Naming Conventions
- **Go files:** `camelCase.go` (e.g., `strMBR.go`, `mkdisk.go`)
- **TypeScript:** `PascalCase.tsx` for components, `camelCase.ts` for services
- **Virtual disks:** `.mia` extension required

### Testing
- **Backend:** Use Go's built-in testing framework
- **Frontend:** Jest and React Testing Library
- **Integration:** Test command execution through API endpoints

### Virtual Disk Management
- **Location:** Virtual disks stored in `backend/Discos/` directory
- **Format:** Binary files with `.mia` extension
- **Structure:** MBR at offset 0, followed by partition data
- **Size:** Support for KB and MB units, calculated as multiples of 1024

## Debugging and Logging

### Log Levels
- **INFO:** General information and successful operations
- **WARNING:** Non-critical issues
- **ERROR:** Critical errors that prevent operation
- **SUCCESS:** Successful completion of operations

### Real-time Monitoring
- **WebSocket:** Connect to `ws://localhost:8080/api/ws` for live logs
- **Server-Sent Events:** Available at `/api/logs/stream`
- **HTTP Polling:** GET `/api/logs` for batch log retrieval

## Important Notes

### Academic Project Requirements
- **Backend Language:** Must use Go (required by course)
- **File System:** EXT2/EXT3 simulation only
- **Disk Files:** Must not exceed specified size limits
- **Command Compatibility:** Must support MIA command specification
- **Documentation:** Extensive technical and user documentation required

### Limitations
- **No Real Disk Access:** Only works with virtual `.mia` files
- **Memory Usage:** Structures should not hold file/directory data in memory
- **File Growth:** Virtual disk files must maintain fixed size
- **Command Set:** Limited to MIA-specified commands only

### Development Phases
1. **Phase 1:** Disk management (mkdisk, fdisk, mount) ✅ Partially complete
2. **Phase 2:** EXT2 file system (mkfs, login/logout) 🔄 In progress
3. **Phase 3:** User/group management 🔄 Pending
4. **Phase 4:** File operations (mkdir, mkfile, cat) 🔄 Pending
5. **Phase 5:** Reports and documentation 🔄 Pending

## Common Tasks

### Adding New Commands
1. Create command file in `backend/command/disk/` or appropriate subdirectory
2. Implement command parsing in command parser
3. Add validation and error handling
4. Update API endpoints if needed
5. Add frontend support for new command

### Testing Commands
```bash
# Via command line tool (when available)
./backend -command "mkdisk -size 10 -path test.mia"

# Via API
curl -X POST http://localhost:8080/api/execute \
  -H "Content-Type: application/json" \
  -d '{"command": "mkdisk -size 10 -path test.mia"}'
```

### Troubleshooting
- **Port conflicts:** Backend uses port 8080, frontend uses 3000
- **CORS issues:** Backend configured for localhost:3000 origin
- **File permissions:** Ensure write access to `backend/Discos/` directory
- **Go version:** Requires Go 1.21 or higher
- **Node version:** Compatible with Node.js 16+