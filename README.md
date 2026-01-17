# NickNamer

A Discord bot designed to randomize nicknames within your server. It offers a fun and interactive way to mix things up and keep your server lively.

## Overview

NickNamer is written in Go and uses Discord slash commands for a modern bot experience.

### Setup Instructions

1. Ensure Go 1.22 or higher is installed.
2. Clone this repository to your server.
3. Set the `API_TOKEN` environment variable to your Discord bot token.
4. Build and run:
   ```bash
   make build
   ./bin/nicknamer
   ```
   Or run directly:
   ```bash
   make run
   ```

## Bot Commands

NickNamer implements slash commands to manage nicknames within your server:

- `/remember [words]`: Adds words to the pool of nickname parts.
- `/forget [word]`: Removes a word from the pool of nickname parts.
- `/forgetall`: Clears the pool of nickname parts.
- `/names`: Lists all words currently in the pool.
- `/randomizeme [words]`: Randomizes the nickname of the command issuer using the specified number of words from the pool.
- `/randomize [user] [words]`: Randomizes the nickname of the specified user.
- `/randomizeall [words]`: Randomizes the nicknames of all server members with the specified role.
- `/rolename [role]`: Views or sets the role name required to use the randomization commands.
- `/flip [user]`: Reverses the order of words in the user's nickname.
- `/reloadnames`: Reloads nickname data from disk.

### Requirements for Running Commands

- The bot will only change the names of users with the specified role. Change the role with `/rolename`.

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `API_TOKEN` | Discord bot token (required) | - |
| `DATA_FILE` | Path to data file | `data.json` |

## Building

```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

Remember to keep your `API_TOKEN` secure and never share it publicly.

