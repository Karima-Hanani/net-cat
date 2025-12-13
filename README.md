# 💬 TCP Chat Server (Go)

![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![TCP](https://img.shields.io/badge/Protocol-TCP-blue)
![Status](https://img.shields.io/badge/Status-Stable-green)

A simple **multi-user TCP chat server** written in **Go**, supporting concurrent clients, message history, and real-time broadcasting.

---

## ✨ Features

- 👥 Up to **10 concurrent users**
- 🔐 Unique usernames
- 🕒 Timestamped messages
- 📜 Chat history for new users
- 📡 Real-time message broadcasting
- 🧼 UTF-8 input sanitization
- 🎨 ASCII art welcome banner
- 💬 Built-in chat commands

---

## 🧾 Commands

| Command | Description |
|-------|-------------|
| `/users` | Show all connected users |
| `/quit` | Leave the chat |

---

## 🚀 Getting Started

### 📦 Prerequisites

- Go **1.18+**
- Terminal (Linux / macOS / Windows)
- `netcat` or `telnet`

---

## 🛠 Installation

```bash
git clone https://github.com/your-username/tcp-chat.git
cd tcp-chat
go build -o TCPChat
