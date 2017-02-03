For each bot that ReBot has to manage, made **at least** a (BotName).bat file in /bots/ that makes the bot run, and return on crash or completion

### Autostart:

add the names of the bots to automatically execute on startup in this file, with "+" in between:

`Bot + Bot2`

if the bot supports or returns with upgrade requests, restart requests, or more, also make a (BotName).bot file in /bots/, where every bool that is returned (in order) will launch a defined sub-bat that's placed in the bots folder with (BotName)-(subcmd).bat, like this:

```
Bot:../Bot,1:restart,2:upgrade
```