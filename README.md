# Snailrace
> 28th June 2023

Welcome to the thrilling world of snail racing! A Discord bot brings you an 
exhilarating game where you can race snails in races. But it's not just about 
winning races - you can buy, sell, and trade snails with different stats, breed 
them together to create the ultimate snail, and earn achievements and badges 
that will set you apart from the rest.

Each user has their own profile that displays their wins, virtual snail money, 
snail count, user level and experience, and recent achievements. But it's the 
snails themselves that steal the show. Each snail is its own entity with unique 
stats, personality, and history. Their owners are tracked and their number of 
wins and win rate are proudly displayed. The ranking system ensures that the 
best of the best are recognized, and even the snail's mood can be a factor in 
their performance.

With a range of general, trading, racing, and breeding commands at your 
disposal, you'll have everything you need to master the world of snail racing. 
So what are you waiting for? Get ready to feel the rush of the race and the 
thrill of creating the ultimate snail in snailrace.

> **Important:** This started out as a simple game in the [UQCS Discord Bot]
>                but people asked if it could be more indepth. It got to the 
>                point where it would make more sense to build it into it's own
>                bot so it isn't just limited to the UQCS Discord Server.

## How to Run

The first thing you will need to do is create a bot application on discord. This
can be done by following the [Discord Dev Doc]. Once you have made a bot, create
or fetch you Bot Token and put it in a `.env` file.

> **Note:** There is an example `.env` file called `.env.example`, you can 
>           rename this in your directory to `.env` and use that.

Once you have setup your app, bot and added it into a discord server, then you 
are ready to run the program. Eventually there will be binaries avaliable but
for now you can just compile it yourself. You will need `go` installed, if you
don't then go to the [Go Install] docs. Once you have the `go` then you can run
the following in the project directory:

```bash
# Build and Compile the program
go build cmd/snailrace

# Run the Program
./snailrace
```

## User Profiles

Your user profile in snailrace is your gateway to snail racing glory. Your 
profile card displays your number of wins, win rate, virtual snail money, snail 
count, user level and experience, as well as your recent achievements. But 
that's not all - the profile card also showcases a small number of badges that 
highlight your accomplishments and set you apart from other players.

Whether you're a seasoned pro or a newbie to the game, your profile is your way 
of showing off your snail racing prowess and displaying your achievements for 
all to see. So get racing, earn those achievements, and collect those badges to 
build the ultimate snail racing profile in snailrace.

## Snails

In snailrace, each snail is stored as an object with its own unique set of stats
and personality traits. The snail object tracks the following:

- `name` (randomly generated `adjetive + "-" + noun`)
- `level`
- `experience`
- `races` raced
- `wins`
- `original owner`
- `current owner`
- `ranking` based on movement stats,
- `mood` (ranging from `sad`, `happy` and `focused`).
- `personality` (either `Hasty`, `Lax`, `Quirky`, `Jolly`, `Quiet`)

One key aspect of each snail is its `Step Size Interval`, which is calculated 
based on its max `speed`, `stamina`, and `weight` which are all within the range 
of `0 to 10`. During each race step, the snail calculates how far it can move 
using a formula that takes into account its speed, stamina, weight, and mood, as 
well as its previous step size and a randomly generated bias value.

## Achievements

As alluded to earlier, there will be achievements in this. This isn't to much of
a focus currently in early implmentations, but there is a draft document with
what is planned [here](./docs/draft_achievements.md).

## Commands

The following are the active commands in the bot. Though if you are interested
there is a draft document with what is planned though not set in stone, which 
can be found [here](./docs/draft_commands.md).

> **Note:** All the commands below need to be prefixed with `/snailrace`

### General

- `ping`:
    Sends the user a `Pong @<user>!` message. This is only for testing.

- `init`:
    Initialises your account with a snail and user record. This is required to
    run any other command.

- `host`:
    Starts the race cycle, there are the following flags to customise the race:

  - `no-bets` Which removes betting from a race
  - `only-one` The race will replay up to 5 times or until the race doesn't 
    finish in a tie.
  - `dont-fill` If there are less than 4 racers, dont fill with randoms.

- `join`:
    Joins a specific race using a `race_id`. This is if you don't want to use 
    the race join buttons.

- `bet`:
    Place a bet, if you have the funds, on a specific snail in a specific race.

[UQCS Discord Bot]: https://github.com/UQComputingSociety/uqcsbot-discord
[Discord Dev Doc]: https://discord.com/developers/docs/getting-started
[Go Install]: https://go.dev/doc/install