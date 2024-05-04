import os
import json
import random
import logging
import discord

from discord.ext import commands

CLIENT_TOKEN = os.getenv('API_TOKEN')

handler = logging.FileHandler(filename='discord.log', encoding='utf-8', mode='w')

class NickNamer:
    """A class to handle IO and other naming operations."""
    def __init__(self, data_file: str=f"{os.path.expanduser('~')}/data.json"):
        self.fn = data_file
        self.names_dict = {}  # Changed from list to dict to handle server-specific lists
        try:
            self.load()
        except: 
            self.names_dict = {}
            self.save()
    
    def change_role_name(self, role_name: str, server_id):
        self.names_dict[server_id]["NICKNAMER_ROLE"] = role_name
        self.save()

    def get_role_name(self, server_id):
        if server_id in self.names_dict and "NICKNAMER_ROLE" in self.names_dict[server_id]:
            return self.names_dict[server_id]["NICKNAMER_ROLE"]
        return

    def remember(self, name, server_id) -> bool:
        """Remember a name for a specific server."""
        if server_id not in self.names_dict:
            self.names_dict[server_id] = {}
            self.names_dict[server_id]["words"] = []
            self.names_dict[server_id]["NICKNAMER_ROLE"] = "NameChanger"
        if name.lower() in self.names_dict[server_id]["words"]:
            return False
        self.names_dict[server_id]["words"].append(name)
        self.save()
        return True

    def forget(self, name, server_id):
        """Forget a name for a specific server."""
        if server_id in self.names_dict and name in self.names_dict[server_id]["words"]:
            self.names_dict[server_id]["words"].remove(name)
            self.save()

    def forget_all(self, server_id):
        """Forget all names for a specific server."""
        if server_id in self.names_dict:
            self.names_dict[server_id]["words"] = []
            self.save()

    def save(self):
        """Write remembered names to disk."""
        with open(self.fn, 'w') as f:
            json.dump(self.names_dict, f)

    def load(self):
        """Load remembered names from disk."""
        if not os.path.exists(self.fn):
            f = open(self.fn, 'x')
            f.close()
            return
        with open(self.fn, 'r') as f:
            self.names_dict = json.load(f)

    def name_list(self, server_id) -> list:
        """Return a list of names that are remembered for a specific server."""
        if server_id not in self.names_dict:
            return []
        return self.names_dict[server_id]["words"]

    def has_names(self, server_id) -> bool:
        """Check if there are any names remembered for a specific server."""
        if server_id not in self.names_dict:
            return False
        if "words" not in self.names_dict[server_id]:
            return False
        if not self.names_dict[server_id]["words"]:
            return False
        return bool(self.names_dict[server_id]["words"])

    def new_name(self, server_id, n=2) -> str:
        name_list = self.name_list(server_id)
        if not name_list:
            return []
        if len(name_list) < n:
            return []
        if n < 1: 
            return []
        name: str = ""
        i: int = 0
        while i != n:
            new_name = name_list[random.randint(0, len(name_list) - 1)]
            if not new_name in name:
                i = i + 1
                name = name + " " + new_name
        return name

async def set_name(member: discord.Member, name: str) -> bool:
    """Helper function to set the name of a member"""
    try:
        await member.edit(nick=name)
        return True
    except discord.errors.Forbidden as err:
        print(f"{type(err).__name__} was raised: {err}")
        return False

def has_role(member: discord.Member, role_name: str) -> bool:
    """Check if a member has a specific role."""
    return any(role.name == role_name for role in member.roles)

intents = discord.Intents.default()
intents.message_content = True
intents.members = True

bot = commands.Bot(command_prefix='?', intents=intents)

nick = NickNamer()

@bot.event
async def on_ready():
    print(f'We have logged in as {bot.user}')

@bot.command(name="remember", aliases=['rem'], help="[remember|rem] word_to_remember - Remember a word.")
async def remember(ctx, *name: str):
    if not name:
        await ctx.send('You must provide a word to remember...')
        return
    err: list = []
    for w in name:
        if not nick.remember(w, str(ctx.guild.id)):
            err.append(w)
    s = ' '
    if err:
        await ctx.send(f'\nI tried my best, but I could not remember _{s.join(err)}_... 🤔\n\nProbably because IT\'S ALREADY THERE. You would know that if you read the word list 🙄\nNext time open your eyes :blush:')
        return
    await ctx.send(f'Okay! I will remember {s.join(name)}!')

@bot.command(name="rolename", help="[rolename] role_name - Set the role required")
async def role_name(ctx, role_name: str = None):
    # If role_name is not provided, return the current role name.
    if not role_name:
        await ctx.send(f'The current role name is {nick.get_role_name(str(ctx.guild.id))}... Can you not read?')
        return
    # Check if owner first
    if not ctx.author.guild_permissions.administrator:
        await ctx.send("You must be an administrator to use this command... Seek gainful employment instead...")
        return
    nick.change_role_name(role_name, str(ctx.guild.id))
    await ctx.send(f'Okay! I will now require the role {role_name} to change names.')

@bot.command(name="forget", help="[forget] word_to_forget - Forget a word.")
async def forget(ctx, name: str = None):
    if not name:
        await ctx.send('You must provide a word to forget... 🙄')
        return
    nick.forget(name, str(ctx.guild.id))
    await ctx.send(f'Okay! I will try to forget {name}!')

@bot.command(name="forgetall", help="[forgetall] - Forget all names.")
async def forget_all(ctx):
    nick.forget_all(str(ctx.guild.id))
    await ctx.send('Okay! Now I\'m useless 🙃')

@bot.command(name="names", aliases=["words", "n"], help="[names|words|n] - List all words remembered.")
async def get_names(ctx):
    if not nick.name_list(str(ctx.guild.id)):
        await ctx.send('🤔 I remember... nothing :(')
        return
    s :str = " "
    await ctx.send(f'🤔 I remember... \n\n```{" ".join(nick.name_list(str(ctx.guild.id)))}```')

@bot.command(name="randomizeme", aliases=['randme'], help="[randomizeme|randme] num_of_words - Randomize your nickname.")
async def randomize_me(ctx, n=2):
    if not nick.has_names(str(ctx.guild.id)):
        await ctx.send("I don't remember any words... 🤔\n\nThat's the whole point of this...")
        return
    if not has_role(ctx.author, nick.get_role_name(str(ctx.guild.id))):
        await ctx.send("Sorry, you do not have the " + nick.get_role_name(str(ctx.guild.id)) + " role... Get a job...")
        return
    name: str = nick.new_name(str(ctx.guild.id), n)
    await ctx.send(f"We will all call {ctx.author.name} _{name}_!")
    r = await set_name(ctx.author, name)
    if not r:
        await ctx.send("Something happened :(")

@bot.command(name="randomize", aliases=['rand'], help="[randomize|rand] user num_of_words - Randomize someone's nickname")
async def randomize(ctx, member: discord.Member, n=2):
    if not nick.has_names(str(ctx.guild.id)):
        await ctx.send("I don't remember any words... 🤔\n\nThat's the whole point of this...")
        return
    if not has_role(member, nick.get_role_name(str(ctx.guild.id))):
        await ctx.send("Sorry, " + member.name + " doesn't have the " + nick.get_role_name(str(ctx.guild.id)) + " role.")
        return
    name: str = nick.new_name(str(ctx.guild.id), n)
    await ctx.send(f"We will all call {member.name} _{name}_!")
    r = await set_name(member, name)
    if not r:
        await ctx.send("Something happened :(")

@bot.command(name="randomizeall", aliases=['randall'], help="[randomizeall|randall] num_of_words - Randomize everyone's nickname.")
async def randomize_all(ctx, n: int=2):
    if not nick.has_names(str(ctx.guild.id)):
        await ctx.send("I don't remember any words... 🤔\n\nThat's the whole point of this...")
        return
    guild = ctx.guild
    member_list: str = ""
    for member in guild.members:
        if not has_role(member, nick.get_role_name(str(ctx.guild.id))):
            continue
        if member.bot:
            continue
        if member.id == guild.owner.id:
            continue
        name: str = nick.new_name(str(guild.id), n)
        member_list = member_list + "\n" + (f"{member.name} ---> _{name}_!")
        r: bool = await set_name(member, name)
        if not r:
            await ctx.send(f"Something happened processing - {member.name} :(")
            return
    await ctx.send(f"Okay, it's a party! :tada::tada:\n{member_list}")

@bot.command(name="flip", help="[flip]")
async def flip(ctx, member: discord.Member=""):
    if not member:
        member = ctx.author
    if not has_role(member, nick.get_role_name(str(ctx.guild.id))):
        await ctx.send("Ain't no way I'm doing that for you.")
        return
    words = member.nick.split()
    reversed_words = words[::-1]
    old_nick = member.nick
    new_nick = ' '.join(reversed_words)
    r = await set_name(member, new_nick)
    if not r:
        await ctx.send("Something happened :(")
    await ctx.send(f"Okay! I flipped {member.name} from _{old_nick}_ to _{new_nick}_!")

@bot.command(name="reloadnames", aliases=['rl'], help="[reloadnames|rl] - Re-read data from disk.")
async def reload_names(ctx):
    nick.load()

bot.run(CLIENT_TOKEN, log_handler=handler)
