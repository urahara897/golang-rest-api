# Event Planner API

This is a RESTful API backend for an event planner application built with Go.

## Live Demo

The API is hosted on Vercel:

```
https://eventplanner-eight.vercel.app
```

## Features

### User Management

- User signup and login
- Admin signup and login
- JWT-based authentication
- Password hashing for security

### Event Management

- Create, read, update, and delete events
- List all events
- Get single event details
- Event registration system
- Cancel event registrations

### Admin Features

- View all users
- Delete users
- Manage events

## Database Setup

The application uses SQLite with two different configurations:

### Local Development

- Uses file-based SQLite database (`./data/app.db`)
- Data persists between server restarts
- Full CRUD functionality

### Vercel Deployment

- Uses in-memory SQLite database
- Data is not persistent between API calls due to serverless architecture
- Each function invocation creates a fresh database
- Suitable for demo purposes

## Why Data Persistence Differs

### Local Development

- Database file is stored on disk
- Data persists between server restarts
- All CRUD operations are permanent

### Vercel (Production)

- Serverless functions run in isolated environments
- Each API call starts with a fresh in-memory database
- Data resets between function calls because:
  - Vercel's filesystem is read-only
  - Serverless functions are stateless
  - In-memory database exists only during function execution

## API Routes

- `/signup`, `/login` - User authentication
- `/admin/signup`, `/admin/login` - Admin authentication
- `/events` - Event management
- `/events/:id/register` - Event registration
- `/users` - User management (Admin only)

## Technologies Used

- Go
- SQLite
- JWT Authentication
- Vercel Serverless Functions
- Bcrypt for password hashing
- Graceful Shutdown

## Server Features

### Graceful Shutdown

The server implements graceful shutdown to ensure clean termination:

- Handles SIGINT (Ctrl+C) and SIGTERM signals
- Allows in-progress requests to complete
- Closes database connections properly
- Uses 5-second timeout for shutdown operations
- Prevents new requests during shutdown

## Note

This project is designed for demonstration purposes. For a production environment, consider using a persistent cloud database service instead of SQLite.
