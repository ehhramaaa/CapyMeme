[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/the_capybara_meme_bot/start?startapp=c749201405a471872c338164f3727bdc)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/CapyMeme/main/demo/demo.png)

# 🔥🔥 Capybara Meme Bot 🔥🔥

### Tested on Windows and Docker Alpine Os with a 4-core CPU using 5 threads.

**Go Version Tested 1.23.1**

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) Must >= 1.23.

- #### Rename config.yml.example to config.yml.
- #### Rename query.txt.example to query.txt and place your query data.
- #### Rename proxy.txt.example to proxy.txt and place your query data.
- #### If you don’t have a query data, you can obtain it from [Telegram Web Tools](https://github.com/ehhramaaa/telegram-web-tools)
- #### It is recommended to use an IP info token to improve request efficiency when checking IPs.

## Features

|        Feature         | Supported |
| :--------------------: | :-------: |
|     Use Query Data     |    ✅     |
|      Auto Staking      |    ✅     |
|    Auto Claim Task     |    ✅     |
| Auto Claim Achievement |    ✅     |
|     Multithreading     |    ✅     |
|    Auto Spin Wheel     |    ✅     |
|         Proxy          |    ✅     |
|   Random User Agent    |    ✅     |

## [Settings](https://github.com/ehhramaaa/CapyMeme/blob/main/configs/config.yml.example)

|        Settings        |                       Description                       |
| :--------------------: | :-----------------------------------------------------: |
|    **IPINFO_TOKEN**    | For Increase Check Ip Efficiency Put Your Token If Have |
|     **AUTO_SPIN**      |             Auto Spin Wheel If Have Ticket              |
|    **AUTO_STAKING**    |   Auto Spin Staking With Custom Amount Of Your Score    |
| **STAKING_PERCENTAGE** |   Percentage Amount To Stake Your Score (e.g. 1-100)    |
|    **RANDOM_SLEEP**    |      Delay before the next lap (e.g. [1800, 3600])      |
|     **MAX_THREAD**     |             Max Thread Worker Run Parallel              |

## Installation

```shell
git clone https://github.com/ehhramaaa/CapyMeme.git
cd CapyMeme
go run .
```

## Usage

```shell
go run .
```

Or

```shell
go run main.go
```

## Or you can do build application by typing:

Windows:

```shell
go build -o CapyMeme.exe
```

Linux:

```shell
go build -o CapyMeme
chmod +x CapyMeme
./CapyMeme
```
