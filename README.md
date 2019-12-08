# CS3700 Raft Implementation

### Node Roles

Nodes are in one of three roles:
- FOLLOWER  -> Election Timeout
- CANDIDATE -> Election Timeout
- LEADER    -> Heartbeat Timeout

## Timeouts

### Election Timeout

- 300ms to 500ms
- Node that times out votes for self, swaps to CANDIDATE, requests VOTE messages from peers
  - CANDIDATE that receives a VOTE RES resets its ELECTION TIMEOUT (TODO, is this true?)
  - FOLLOWER that receives a VOTE REQ resets its ELECTION TIMEOUT, votes affirmative if not already voted
- If a CANDIDATE has received a majority vote, convert to LEADER

### Heartbeat Timeout

- Heartbeat timeout is 1/5th of the base election timeout
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

### Redirect (response only)

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

## Consensus Messages

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

#### Log Entry

```
{
  "term": "<TERM_ID>"
  "key": "<KEY>",
  "value": "<VALUE>",
  "type": "put"
  "sender": "<NODE_ID>"
  "mid": "<MESSAGE_ID>"
}
```

---


## Testing

- We primarily used the simulator provided to us to test our code. 
- The simulator provided good feedback about what our code was messing up. Such as, too many duplicate messages, or a constant flip-flopping of the leader.

## Challenges Faced

- We ran into issues decoding messages during a partition.
 It was an incredibly odd bug where the buffer we were reading from would get half a message in it. 
This led to our decoder erroring as it couldn't read the junk data.
 Unfortunately, we didn't log this error this error at first so it just looked like nodes suddenly stopped receiving any input. 
- Interpreting the paper was a challenge because we didn't immediately understand the reasoning behind all of their decisions.
Example: 1 indexing the log of entries. It makes sense that you need to 1 index it to get it to work with the intialized values for commitIndex and lastApplied.
If log[] weren't 1 indexed, commitIndex and lastApplied would need to be -1, which doesn't make a ton of sense.

## Running

With the compiled binary, `3700kvstore`:

```
./3700kvstore <your ID> <ID of second replica> [ID3 [ID4 ...]]
```


