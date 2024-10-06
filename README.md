# astronomer

![Latest release](https://img.shields.io/github/v/release/QuokkaStake/astronomer)
[![Actions Status](https://github.com/QuokkaStake/astronomer/workflows/test/badge.svg)](https://github.com/QuokkaStake/astronomer/actions)
[![codecov](https://codecov.io/gh/QuokkaStake/astronomer/graph/badge.svg?token=JhR7t6G1s6)](https://codecov.io/gh/QuokkaStake/astronomer)

astronomer is your pocket Telegram cosmos-sdk multichain explorer and wallet!

## Why is it cool?
- Can work with multiple chain
- Uses PostgreSQL as a database to store all data on one place
- Allows you to set it up in runtime and avoid patching configuration file every time you need to update something
- Allows you to fetch proposals, chain params, wallet balances, validators info and many more without leaving Telegram
- Allows working with it in both chats and in private DMs
- Allows binding specific chains for a specific chat
- (TODO) Comes with Prometheus metrics, so you can observe if something is wrong
- (TODO) Includes authz-based non-custodial wallet that allows you to interact with the blockchain while owning your wallet keys

## How can I set it up?

Download the latest release from [the releases page](https://github.com/QuokkaStake/astronomer/releases/). After that, you should unzip it, and you are ready to go:

```sh
wget <the link from the releases page>
tar <downloaded file>
./astronomer --config <path to config>
```

Alternatively, install `golang` (>1.18), clone the repo and build it:
```
git clone https://github.com/QuokkaStake/astronomer
cd astronomer
# This will generate a `astronomer` binary file in the repository folder
make build
# This will generate a `astronomer` binary file in $GOPATH/bin
```

To run it detached, first we have to copy the file to the system apps folder:

```sh
sudo cp ./astronomer /usr/bin
```

Then we need to create a systemd service for our app:

```sh
sudo nano /etc/systemd/system/astronomer.service
```

You can use this template (change the user to whatever user you want this to be executed from.
It's advised to create a separate user for that instead of running it from root):

```
[Unit]
Description=Astronomer
After=network-online.target

[Service]
User=<username>
TimeoutStartSec=0
CPUWeight=95
IOWeight=95
ExecStart=astronomer --config <config path>
Restart=always
RestartSec=2
LimitNOFILE=800000
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

Before starting, do not forget to run migrations (actually, better to do it every time you update the app):
```sh
astronomer migrate --config path/to/config.toml
```

Then we'll add this service to autostart and run it:

```sh
sudo systemctl daemon-reload # reload config to reflect changed
sudo systemctl enable astronomer # put service to autostart
sudo systemctl start astronomer # start the service
sudo systemctl status astronomer # validate it's running
```

If you need to, you can also see the logs of the process:

```sh
sudo journalctl -u astronomer -f --output cat
```

## How does it work?

It runs a bunch of Interacters (currently, Telegram only) and once it receives a query,
it will fetch the data from both chain and the database and return the answer to the user.
Some commands are allowed for admins only, for example creating chains/denoms/explorers, others
are free to use for everybody.

## How can I configure it?

All configuration is done via `.toml` config file, which is mandatory. Run the app with `--config <path/to/config.toml>`
to specify config. Check out `config.example.toml` to see the params that can be set.

## Notifiers

Currently, this program supports the following notifications channels:
1) Telegram

Go to @BotFather in Telegram and create a bot. After that, there are three options:
- you want to send messages to a user. This user should write a message to @getmyid_bot, then copy
the `Your user ID` number. Also keep in mind that the bot won't be able to send messages unless you contact it first,
so write a message to a bot before proceeding.
- you want to send messages to a channel. Write something to a channel, then forward it to @getmyid_bot and copy
the `Forwarded from chat` number. Then add the bot as an admin.
- you want to send message to a chat. Add @raw_data_bot to this chat, write something, then copy a channel_id
from bot response (starts with a minus), then you can remove @raw_data_bot from the channel.

To have fancy commands auto-suggestion, go to @BotFather again, select your bot -> Edit bot -> Edit description
and paste the following:
```
start - Displays bot info
help - Displays bot info
balance - Display your wallets' balance, delegations, rewards etc.
validator - Search for a validator
validators - Display info on validators you are subscribed to
params - Display chain(s) params
proposals - Display all active proposals
proposal - Display a proposal by ID
wallet_link - Link a wallet
wallet_link - Unlink a wallet
validator_link - Link a validator
validator_unlink - Unlink a validator
chains - Display all chains and the chains bound to this chat
supply - See total chain supply, bonded ratio and community pool
```

Then add a Telegram config to your config file (see `config.example.toml` for reference).

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
