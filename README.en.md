## Counter-Strike Online 2 Server 

[![Build status](https://ci.appveyor.com/api/projects/status/a4pj1il9li5s08k5?svg=true)](https://ci.appveyor.com/project/KouKouChan/cso2-server)
[![](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![](https://img.shields.io/badge/version-v0.8.0-blue)](https://github.com/NightElff/CSO2-Server/releases)

### 0x01 Description

Counter-Strike Online 2 Server

***Based on [cso2-master-server](https://github.com/L-Leite/cso2-master-server) by l-leite.***

***Localization files are welcome！See part 3.***

### 0x02 Plan

    1. Basic game play √
    2. Refactoring code ...

### 0x03 Localize

```
1. Create a localization file in CSO2-Server\locales\ folder,such as en-us.ini
2. Edit you localization file , you can see zh-cn.ini
3. Edit server.conf , and set LocaleFile to your file name
```

### 0x04 Play

    1. You should have a game client with korea version.
    2. Download a launcher from l-leite's github page.
    3. Download lastest game server file from  ( https://github.com/KouKouChan/CSO2-Server/releases )
    4. Start game server and use bat file to launch your game.
    5. Have fun!

**Attention**!

- If you want to enable registration , you should modify server.conf file , and set EnableRegister to 1 , and you must set your email smtp server and your email code.then you can open localhost:1314 with your browser.

### 0x05 Build

    1. open the folder
    2. enter "go build" command to build
    3. run it

### 0x06 Build Env

    Go 1.15.6
    Use port:30001-TCP、30002-UDP

***If you want to set up a LAN or Internet Server, please open the port of firewall.***

### 0x07 Screenshots

![Image](./photos/main.png)

![Image](./photos/intro.png)

![Image](./photos/channel.png)

![Image](./photos/ingame.jpg)

![Image](./photos/result.jpg)
