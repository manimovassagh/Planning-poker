
# Planning Poker API

Planning Poker API is a collaborative tool designed to help teams estimate the complexity and time of tasks during sprint planning sessions. This API allows users to create planning sessions, join as participants, submit estimates (votes) for tasks, and reveal the votes collectively for a consensus-based estimate.

## Project Overview

The Planning Poker API provides endpoints to:
- Register users
- Log in users and manage authentication
- Create planning sessions and tasks
- Allow users to join sessions, submit votes, and reveal final votes

This project is built to streamline sprint planning for Agile teams, making it easier to estimate tasks with a structured process.

## Technologies Used

- **Backend**: Go, Fiber
- **Database**: SQLite (for development)
- **Testing**: `.http` files for API endpoint testing

## Getting Started

### Prerequisites

Ensure the following are installed on your local machine:

- [Go](https://golang.org/dl/) (latest version)
- [SQLite](https://www.sqlite.org/download.html)
- HTTP client or REST Client extension in VS Code (for running `.http` files)

### Project Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/planning-poker.git
   cd planning-poker
   ```

2. **Install dependencies**:

   Run the following command to install necessary Go modules:

   ```bash
   go mod tidy
   ```

3. **Run Database Migrations**:

   Ensure the SQLite database is set up and run migrations if required.

4. **Start the Server**:

   Run the following command to start the server:

   ```bash
   go run main.go
   ```

   The server will start at `http://localhost:3000`.

### Folder Structure

- **cmd/** - Contains application entry points.
- **internal/** - Contains core modules, services, and handlers for API requests.
- **test/** - Contains `.http` files for running API tests (see below).

## Testing the API

The `test` folder includes `.http` files to test various API endpoints sequentially.

### Testing Workflow Overview

1. **Register Users**: Registers the admin and five regular users.
2. **Login Users**: Logs in each user and captures their JWT tokens for use in the next steps.
3. **Session and Voting**: Covers session creation, task creation, user joining, voting, and vote revealing.

### Setup Instructions for Testing

1. **Open the `test` folder**: Navigate to the `test` folder in your project.
2. **Open the `.http` files**: Files are organized as follows:
   - **RegisterUsers.http**: For registering the admin and all users.
   - **LoginUsers.http**: For logging in all users and capturing tokens.
   - **SessionAndVoting.http**: For creating sessions, adding tasks, joining sessions, voting, and revealing votes.

3. **Run Each File Sequentially**: Execute each `.http` file in order for a seamless workflow.

### HTTP Files Details

- **RegisterUsers.http**: Registers each user individually, including the admin and five regular users.
- **LoginUsers.http**: Logs in each user and captures their tokens.
- **SessionAndVoting.http**: Includes:
  - **Session Creation**: The admin creates a session.
  - **Task Creation**: A task is created within the session by the admin.
  - **User Joins**: Each user joins the session.
  - **Voting**: Each user submits a vote.
  - **Vote Reveal**: The admin reveals the votes, concluding the session.

## Example Usage

1. Open the `.http` files sequentially from the `test` folder.
2. Follow the comments in each file to complete the workflow.

## Tips

- Ensure your server is running and accessible at `http://localhost:3000`.
- Using the **REST Client** in VS Code is recommended for running `.http` files efficiently.

## Final Notes

These `.http` files provide a clear and direct approach to testing, without the need for external collections. Each step is designed to be run sequentially, ensuring a smooth and easy-to-follow testing process.

Happy testing with Planning Poker API!
