# NickNamer

A Discord bot designed to randomize nicknames within your server. It offers a fun and interactive way to mix things up and keep your server lively. 

## Overview

Setting up NickNamer is straightforward. You'll need to have Python installed on your server, along with the discord.py library. The bot uses a `.env` file to securely store the `API_TOKEN` required for Discord bot authentication.

### Setup Instructions

1. Ensure Python 3.8 or higher is installed.
2. Install discord.py using pip: `pip install discord.py`.
3. Clone this repository to your server.
4. Create an env var and add your Discord bot token as `API_TOKEN=your_token_here`.
5. Run the bot using `python src/main.py`.

## Bot Commands

NickNamer implements several commands to manage nicknames within your server:

- `?remember [word]`: Adds a word to the pool of nickname parts.
- `?forget [word]`: Removes a word from the pool of nickname parts.
- `?forgetall`: Clears the pool of nickname parts.
- `?names`: Lists all words currently in the pool.
- `?randomizeme [num_of_words]`: Randomizes the nickname of the command issuer using the specified number of words from the pool.
- `?randomize [user] [num_of_words]`: Randomizes the nickname of the specified user.
- `?randomizeall [num_of_words]`: Randomizes the nicknames of all server members with the specified role.
- `?rolename [role_name]`: Sets the role name required to use the randomization commands.
- `?flip [user]`: Reverses the order of words in the user's nickname.

### Requirements for Running Commands

- The bot will only change the names of users with the specified role. Change the role with `?rolename`.

## Requirements

- Python 3.8 or higher.
- discord.py library.
- A Discord bot token.

Remember to keep your `API_TOKEN` secure and never share it publicly.

