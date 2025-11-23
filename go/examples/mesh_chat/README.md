# Mesh Chat Example

This example demonstrates how to build a group chat application using CloudBridge SDK's mesh networking capabilities.

## What This Example Shows

1. **Mesh Network Join** - How to connect to a mesh network
2. **Peer Discovery** - How to find other peers in the network
3. **Broadcasting** - How to send messages to all peers
4. **Event Handling** - How to process incoming messages
5. **Interactive Chat** - How to build a chat interface

## Running the Example

```bash
cd examples/mesh_chat
export CLOUDBRIDGE_TOKEN="your-jwt-token"
go run main.go
```

## Expected Output

```
CloudBridge SDK - Mesh Chat Example
====================================

✓ Chat client created
  Username: anonymous

Joining mesh network: chat-network
⚠ Failed to join mesh (expected without relay): not implemented
  In production, the client would join the mesh network

=== Chat Interface ===
Commands:
  /peers         - List connected peers
  /send <msg>    - Broadcast message to all peers
  /quit          - Exit chat

Type your message and press Enter to broadcast:

>
```

## Features

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/peers` | List all connected peers | `/peers` |
| `/send <msg>` | Send message to all peers | `/send Hello everyone!` |
| `/quit` or `/exit` | Exit the chat | `/quit` |
| `/help` | Show help | `/help` |

### Regular Chat

Just type your message and press Enter to broadcast to all peers:

```
> Hello, everyone!
✓ Message sent to mesh
```

## Architecture

### Chat Message Structure

```go
type ChatMessage struct {
    From      string        // Sender's username
    To        string        // Recipient (empty for broadcast)
    Content   string        // Message content
    Timestamp time.Time     // When message was sent
}
```

### Message Flow

```
Client A          Relay Server          Client B          Client C
  |                   |                    |                 |
  +---Message-------->|                    |                 |
  |                   |--Broadcast-------->|                 |
  |                   |--Broadcast--------|-------->|        |
  |                   |<---ACK------------|-----------|------>|
  |<-------ACK--------|                   |                 |
```

## Production Implementation

For a production mesh chat, you would:

### 1. Implement Message Queue

```go
type MessageQueue struct {
    messages chan ChatMessage
    buffer   []ChatMessage
    maxSize  int
}

func (q *MessageQueue) Push(msg ChatMessage) error {
    select {
    case q.messages <- msg:
        return nil
    default:
        return ErrQueueFull
    }
}
```

### 2. Add Persistence

```go
type PersistentChat struct {
    db    *sql.DB
    mesh  cloudbridge.Mesh
}

func (pc *PersistentChat) SaveMessage(msg ChatMessage) error {
    _, err := pc.db.Exec(
        "INSERT INTO messages (from, content, timestamp) VALUES (?, ?, ?)",
        msg.From, msg.Content, msg.Timestamp,
    )
    return err
}
```

### 3. Add User Management

```go
type User struct {
    ID       string
    Username string
    Status   string // online, offline, away
    PeerID   string
}

func (u *User) IsOnline() bool {
    return u.Status == "online"
}
```

### 4. Implement Message Filtering

```go
type MessageFilter struct {
    MinLength int
    MaxLength int
    Blocklist []string
}

func (f *MessageFilter) Validate(msg string) error {
    if len(msg) < f.MinLength {
        return ErrMessageTooShort
    }
    if len(msg) > f.MaxLength {
        return ErrMessageTooLong
    }
    for _, blocked := range f.Blocklist {
        if strings.Contains(msg, blocked) {
            return ErrForbiddenContent
        }
    }
    return nil
}
```

### 5. Add Encryption

```go
func (mb *MessageBroadcaster) SendEncrypted(msg ChatMessage) error {
    encrypted, err := mb.cipher.Encrypt(msg.Content)
    if err != nil {
        return err
    }

    encryptedMsg := msg
    encryptedMsg.Content = encrypted
    return mb.Send(encryptedMsg)
}
```

## Example: Building a Real Chat

Here's how you'd extend this example for production:

```go
package main

import (
    "context"
    "database/sql"
    "sync"
    "time"

    "github.com/twogc/cloudbridge-sdk/go/cloudbridge"
    _ "github.com/mattn/go-sqlite3"
)

type ChatServer struct {
    client      *cloudbridge.Client
    mesh        cloudbridge.Mesh
    db          *sql.DB
    users       map[string]*User
    messageQ    chan ChatMessage
    mu          sync.RWMutex
    ctx         context.Context
    cancel      context.CancelFunc
}

func NewChatServer(token string) (*ChatServer, error) {
    client, err := cloudbridge.NewClient(
        cloudbridge.WithToken(token),
    )
    if err != nil {
        return nil, err
    }

    // Create database
    db, err := sql.Open("sqlite3", "./chat.db")
    if err != nil {
        client.Close()
        return nil, err
    }

    ctx, cancel := context.WithCancel(context.Background())

    return &ChatServer{
        client:   client,
        db:       db,
        users:    make(map[string]*User),
        messageQ: make(chan ChatMessage, 100),
        ctx:      ctx,
        cancel:   cancel,
    }, nil
}

func (cs *ChatServer) Start() error {
    mesh, err := cs.client.JoinMesh(cs.ctx, "chat-network")
    if err != nil {
        return err
    }
    cs.mesh = mesh

    // Start message processor
    go cs.processMessages()

    return nil
}

func (cs *ChatServer) Stop() error {
    cs.cancel()
    if cs.mesh != nil {
        cs.mesh.Leave()
    }
    if cs.client != nil {
        cs.client.Close()
    }
    if cs.db != nil {
        cs.db.Close()
    }
    return nil
}

func (cs *ChatServer) AddUser(user *User) {
    cs.mu.Lock()
    defer cs.mu.Unlock()
    cs.users[user.ID] = user
}

func (cs *ChatServer) SendMessage(msg ChatMessage) error {
    msg.Timestamp = time.Now()
    return cs.mesh.Send(msg)
}

func (cs *ChatServer) processMessages() {
    for {
        select {
        case msg := <-cs.messageQ:
            // Save to database
            cs.saveMessage(msg)
            // Broadcast to peers
            cs.SendMessage(msg)

        case <-cs.ctx.Done():
            return
        }
    }
}

func (cs *ChatServer) saveMessage(msg ChatMessage) error {
    _, err := cs.db.Exec(
        "INSERT INTO messages (username, content, timestamp) VALUES (?, ?, ?)",
        msg.From, msg.Content, msg.Timestamp,
    )
    return err
}
```

## Testing the Chat

### Local Testing

```bash
# Terminal 1
export CLOUDBRIDGE_USER="alice"
go run main.go

# Terminal 2
export CLOUDBRIDGE_USER="bob"
go run main.go
```

### Simulated Multi-User

```bash
# Create a test harness
go test -v -count=5 ./examples/mesh_chat
```

## Performance Considerations

1. **Message Rate** - Limit messages per user per second
2. **Broadcast Efficiency** - Use tree-based forwarding
3. **Memory** - Implement message cleanup
4. **Persistence** - Use batched writes to database

## Security Considerations

1. **Message Validation** - Check content before broadcasting
2. **Rate Limiting** - Prevent spam
3. **User Authentication** - Verify peer identity
4. **Data Encryption** - Encrypt messages in transit
5. **Access Control** - Limit who can join/message

## Troubleshooting

### No peers visible

```bash
# Try discovering peers manually
export CLOUDBRIDGE_TOKEN="your-token"
./cloudbridge discover
```

### Messages not being received

- Check peer is online: `/peers`
- Verify network connectivity
- Check logs: `export CLOUDBRIDGE_LOG_LEVEL=debug`

### Performance issues

- Reduce message frequency
- Use batching for bulk sends
- Monitor memory usage
- Check relay server metrics

## Next Steps

- See [Echo Server Example](../echo_server/) for simple server
- See [Tunnel Example](../tunnel/) for port forwarding
- Read [Architecture](../../docs/ARCHITECTURE.md) for design details
- Explore [API Reference](../../docs/API_REFERENCE.md)

## Learn More

- [CloudBridge Documentation](https://github.com/twogc/cloudbridge-sdk)
- [Building P2P Applications](https://github.com/twogc/cloudbridge-sdk/wiki)
- [Mesh Networking Patterns](https://github.com/twogc/cloudbridge-sdk/wiki/Mesh-Patterns)
