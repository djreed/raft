# CS3700 Raft Implementation

## State Machine Structure

Each node is a state machine consisting of:
- Log of messages received
- Log of messages committed
- Current key/value state
- Current node role



### Node Role

Nodes are in one of three roles:
- FOLLOWER  -> Election Timeout
- CANDIDATE -> Election Timeout
- LEADER    -> Heartbeat Timeout

## Timeouts

### Election Timeout

- 150ms to 300ms
- Node that times out votes for self, swaps to CANDIDATE, requests VOTE messages from peers
  - CANDIDATE that receives a VOTE RES resets its ELECTION TIMEOUT (TODO, is this true?)
  - FOLLOWER that receives a VOTE REQ resets its ELECTION TIMEOUT, votes affirmative if not already voted

- If a CANDIDATE has received a majority vote, convert to LEADER

### Heartbeat Timeout

- LEADER sends APPEND LOG to FOLLOWERS as heartbeat message on heartbeat timeout lapse
- FOLLOWER nodes reset their ELECTION TIMEOUT on APPEND MESSAGE receipt
  - Respond with APPEND LOG of their own

---

## Message Types (Client)

### GET (request / response)

**request**:
```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "get",
  "MID": "<a unique string>",
  "key": "<some key>"
}
```

**response**:
```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "ok" | "fail",
  "MID": "<a unique string>",
  "value": "<value of the key>" _(if OK_MSG)_
}
```

### PUT (request / response)

**request**:
```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "put",
  "MID": "<a unique string>",
  "key": "<some key>",
  "value": "<value of the key>"
}
```

**response**:
```
{
  "src": "<ID>" 
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "ok" | "fail", 
  "MID": "<a unique string>"
}
```

### REDIRECT (response only)

If the client sends any message (get() or put()) to a replica that is not the leader, it should respond with a redirect

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "redirect",
  "MID": "<a unique string>"
}
```

---

## Consensus Message

### Request Vote

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "vote",
  "MID": "<a unique string>",
  "term": "<TERM ID>"
  "vote": "<ID>"
}
```

### Append Log

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "append",
  "MID": "<a unique string>",
  "term": "<TERM ID>",
  "updates": [<LOG_UPDATE>]
}
```

#### LOG_UPDATE -- Promise

```
{
  "id": "<APPEND_ID>"
  "key": "<KEY>",
  "value": "<VALUE>",
  "type": "promise"
}
```

#### LOG_UPDATE -- Commit

```
{
  "id": "<APPEND_ID>"
  "key": "<KEY>",
  "value": "<VALUE>",
  "type": "commit"
}
```

---

## Running

With the compiled binary, `3700kvstore`:

```
./3700kvstore <your ID> <ID of second replica> [ID3 [ID4 ...]]
```


