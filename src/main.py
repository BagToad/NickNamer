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
        self.names_list = []
        try:
            self.load()
        except: 
            self.names_list = []
            self.save()

    def remember(self, name) -> bool:
        """Remember a name."""
        if name.lower() in self.names_list:
            return False
        self.names_list.append(name)
        self.save()
        return True

    def forget(self, name):
        """Forget a name."""
        self.names_list.remove(name)
        self.save()

    def forget_all(self):
        """Forget all names."""
        self.names_list = []
        self.save()

    def save(self):
        """Write remembered names to disk."""
        with open(self.fn, 'w') as f:
            json.dump(self.names_list, f)

    def load(self):
        """Load remembered names from disk."""
        if not os.path.exists(self.fn):
            f = open(self.fn, 'x')
            f.close()
            return
        with open(self.fn, 'r') as f:
            self.names_list = json.load(f)

    def name_list(self) -> list:
        """Return a list of names that are remembered."""
        return self.names_list

    def new_name(self, n=2) -> str:
        name_list = self.name_list()
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

intents = discord.Intents.default()
intents.message_content = True
intents.members = True

bot = commands.Bot(command_prefix='?', intents=intents)

nick = NickNamer()

@bot.event
async def on_ready():
    print(f'We have logged in as {bot.user}')

@bot.command(name="remember", aliases=['rem'], help="[remember|rem] word_to_remember - Remember a word.")
# @commands.is_owner()
async def remember(ctx, *name: str):
    err: list = []
    for w in name:
        if not nick.remember(w):
            err.append(w)
    s = ' '
    await ctx.send(f'Okay! I will remember {s.join(name)}!')
    if err:
        await ctx.send(f'I could not remember {s.join(err)}')

@bot.command(name="forget", help="[forget] word_to_forget - Forget a word.")
# @commands.is_owner()
async def forget(ctx, name: str):
    nick.forget(name)
    await ctx.send(f'Okay! I will forget {name}!')

@bot.command(name="forgetall", help="[forgetall] - Forget all names.")
# @commands.is_owner()
async def forget_all(ctx):
    nick.forget_all()
    await ctx.send('Okay!')

@bot.command(name="names", aliases=["words", "n"], help="[names|words|n] - List all words remembered.")
async def get_names(ctx):
    if not nick.name_list():
        await ctx.send('I remember... nothing')
        return
    s :str = " "
    await ctx.send(f'I remember... {s.join(nick.name_list())}')

@bot.command(name="randomizeme", aliases=['randme'], help="[randomizeme|randme] num_of_words - Randomize your nickname.")
async def randomize_me(ctx, n=2):
    name: str = nick.new_name(n)
    await ctx.send(f"Okay! I'll call {ctx.author.name} {name}!")
    r = await set_name(ctx.author, name)
    if not r:
        await ctx.send("Something happened :(")

@bot.command(name="randomize", aliases=['rand'], help="[randomize|rand] user num_of_words - Randomize someone's nickname")
async def randomize(ctx, member: discord.Member, n=2):
    name: str = nick.new_name(n)
    await ctx.send(f"Okay! I'll call {member.name} {name}!")
    r = await set_name(member, name)
    if not r:
        await ctx.send("Something happened :(")

@bot.command(name="randomizeall", aliases=['randall'], help="[randomizeall|randall] num_of_words - Randomize everyone's nickname.")
async def randomize_all(ctx, n: int=2):
    guild = ctx.guild
    member_list: str = ""
    for member in guild.members:
        if member.bot:
            continue
        if member.id == guild.owner.id:
            continue
        name: str = nick.new_name(n)
        member_list = member_list + "\n" + (f"I'll call {member.name} {name}!")
        r: bool = await set_name(member, name)
        if not r:
            await ctx.send(f"Something happened processing - {member.name} :(")
            return
    await ctx.send(f"Okay, I'm doing it!\n{member_list}")

@bot.command(name="flip", help="[flip]")
async def flip(ctx, member: discord.Member=""):
    if not member:
        member = ctx.author
    words = member.nick.split()
    reversed_words = words[::-1]
    r = await set_name(member, ' '.join(reversed_words))
    if not r:
        await ctx.send("Something happened :(")

@bot.command(name="reloadnames", aliases=['rl'], help="[reloadnames|rl] - Re-read data from disk.")
async def reload_names(ctx):
    nick.load()

bot.run(CLIENT_TOKEN, log_handler=handler)
