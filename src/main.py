import os
import discord
import json
import random
import logging
from discord.ext import commands

CLIENT_TOKEN = os.getenv('API_TOKEN')

handler = logging.FileHandler(filename='discord.log', encoding='utf-8', mode='w')


class NickNamer:
    def __init__(self):
        self.fn = "/tmp/data.json"
        self.names_list = []
        try:
            self.load()
        except: 
            self.names_list = []
            self.save()

    def rem(self, name):
        self.names_list.append(name)
        self.save()

    def forget(self, name):
        self.names_list.remove(name)
        self.save()
    
    def forgetAll(self):
        self.names_list = []
        self.save()

    def save(self):
        with open(self.fn, 'w') as f:
            json.dump(self.names_list, f)

    def load(self):
        if not os.path.exists(self.fn):
            f = open(self.fn, 'x')
            f.close()
        with open(self.fn, 'r') as f:
            self.names_list = json.load(f)

    def nameList(self):
        return self.names_list

    def newName(self):
        name_list = nick.nameList()
        if not name_list:
            return
        name1 = name_list[random.randint(0, len(name_list) - 1)]
        name2 = name_list[random.randint(0, len(name_list) - 1)]
        while name1 == name2:
            name2 = name_list[random.randint(0, len(name_list) - 1)]
        return [name1, name2]

nick = NickNamer()

intents = discord.Intents.default()
intents.message_content = True
intents.members = True

bot = commands.Bot(command_prefix='?', intents=intents)

@bot.event
async def on_ready():
    print(f'We have logged in as {bot.user}')

@bot.command(aliases=['remember'])
# @commands.is_owner()
async def rem(ctx, name):
    nick.rem(name)
    await ctx.send(f'Okay! I will remember {name}!')

@bot.command()
# @commands.is_owner()
async def forget(ctx, name):
    nick.forget(name)
    await ctx.send(f'Okay! I will forget {name}!')

@bot.command()
# @commands.is_owner()
async def forgetall(ctx):
    nick.forgetAll()
    await ctx.send(f'Okay!')

@bot.command(name="names")
async def getnames(ctx):
    if not nick.nameList():
        await ctx.send(f'I remember... nothing')
        return
    await ctx.send(f'I remember... {nick.nameList()}')

@bot.command(aliases=['randme'])
async def randomizeme(ctx):
    names = nick.newName()
    name1 = names[0]
    name2 = names[1]
    try:
        await ctx.author.edit(nick=name1+" "+name2)
    except discord.errors.Forbidden as err:
        print(f"{type(err).__name__} was raised: {err}")
        await ctx.send(f'Something happened! :(\n{err}')
        return   
    await ctx.send(f'I found {name1} {name2}!')

@bot.command(aliases=['rand'])
async def randomize(ctx, member: discord.Member):
    names = nick.newName()
    name1 = names[0]
    name2 = names[1]
    try:
        await member.edit(nick=name1+" "+name2)
    except discord.errors.Forbidden as err:
        print(f"{type(err).__name__} was raised: {err}")
        await ctx.send(f'Something happened! :(\n{err}')
        return   
    await ctx.send(f'{member.mention} is now {name1} {name2}!')

@bot.command(aliases=['randall'])
async def randomizeall(ctx):
    name_list = nick.nameList()
    if not name_list:
        await ctx.send('I found nothing :(')
        return
    guild = ctx.guild
    member_list = ""
    for member in guild.members:
        if member.bot:
            continue
        names = nick.newName()
        name1 = names[0]
        name2 = names[1]
        member_list = member_list + "\n" + (f"I would have set {member.name} to {name1} {name2}")
    await ctx.send(f'This is what I would have done:{member_list}!')

    # name1 = name_list[random.randint(0, len(name_list) - 1)]
    # name2 = name_list[random.randint(0, len(name_list) - 1)]
    # while name1 == name2:
    #     name2 = name_list[random.randint(0, len(name_list) - 1)]
    
    # try:
    #     await member.edit(nick=name1+" "+name2)
    #     await ctx.send(f'I found {name1} {name2}!')
    # except Exception as err: 
    #     await ctx.send(f"I can't touch {member.mention} because {err=}  :(")
    # except discord.ext.commands.errors.MemberNotFound as err:
    #     await ctx.send(f"I can't touch {member.mention} because {err=}  :(")
bot.run(CLIENT_TOKEN, log_handler=handler)
