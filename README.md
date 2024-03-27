# Godbase

A drop-in replacement for redis written in Go

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) 
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)


## Table of Contents
- [Introduction](#introduction)
- [Feature parity](#feature-parity)
- [Available commands](#available-commands)
    - [The SET Command](#the-set-command)
- [Installation](#installation)
- [Compatibility](#compatibility)
- [License](#license)

## Introduction
`Godbase` is a blazingly fast drop-in replacement for redis, built with golang.

It's a fun project that i started to learn more about redis and golang, taking help from this fantastic [guide](https://www.build-redis-from-scratch.dev/en/introduction) by [ahmedash95](https://github.com/ahmedash95). This project builds upon the concepts in this guide and adds features such as `TTL` options for the `SET` command.

The project is still in its early stages and is not recommended for production use. It's a fun project to learn more about redis and golang.

## Feature parity

| Feature                   | Redis | Godbase |
| ------------------------- | ----- | -------- |
| In-memory key-value store | ✅     | ✅        |
| Strings                   | ✅     | ✅        |
| Persistence               | ✅     | ✅        |
| Hashes                    | ✅     | ✅        |
| TTL                       | ✅     | ✅        |
| Lists                     | ✅     | ❌        |
| Sets                      | ✅     | ❌        |
| Sorted sets               | ✅     | ❌        |
| Streams                   | ✅     | ❌        |
| HyperLogLogs              | ✅     | ❌        |
| Bitmaps                   | ✅     | ❌        |
| Pub/Sub                   | ✅     | ❌        |
| Transactions              | ✅     | ❌        |

## Available commands

The following commands are supported by Godbase as of now:

#### MISC
`PING` 

#### Keys
`DEL` `EXISTS` `KEYS` `EXPIRE` `TTL`

#### Strings
`SET` `GET` `APPEND` `INCR` `INCRBY` `DECR` `DECRBY` `MSET` `MGET`

#### Hashes
`HSET` `HGET` `HGETALL` 

### The SET Command
```
SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | KEEPTTL]
```
The SET command supports the following options:

 - EX seconds -- Set the specified expire time, in seconds (a positive integer).
 - PX milliseconds -- Set the specified expire time, in milliseconds (a positive integer).
 - NX -- Only set the key if it does not already exist.
 - XX -- Only set the key if it already exists.
 - KEEPTTL -- Retain the time to live associated with the key.
 - GET -- Return the old string stored at key, or nil if key did not exist. An error is returned and SET aborted if the value stored at key is not a string.

**Time complexity**: O(1)

## Installation

Clone this repository and use the makefile to build or run the binary.
```
git clone https://github.com/Maniktherana/godbase.git
make build
```

## Compatibility

Godbase is compatible with existing redis clients. You can use the redis-cli to interact with godbase for the supported commands.

```
redis-cli -h localhost -p 6379
localhost:6379> set hello world
OK
localhost:6379> get hello
world
localhost:6379> set god man
OK
localhost:6379> set god base xx
OK
localhost:6379> get god
base
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.