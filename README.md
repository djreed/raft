# CS3700 Raft Implementation

## Node Roles

Nodes are in one of three roles:
- FOLLOWER  -> Election Timeout
- CANDIDATE -> Election Timeout
- LEADER    -> Heartbeat Timeout

---

## Design Notes

To handle quorum tracking, we have two main states for the Leader:
- Steady
- Commit

In the Steady state, GETs are responded to instantly, election messages are handled,
and we append on the heartbeat cycle. If we get a PUT, we add that entry to our
log and transition to Commit.

In Commit, we no longer accept new Put messages, and instead wait until a quorum
of nodes in our system have replicated up to the latest message in our log.
Because our leader only sends new data to replicas on each heartbeat, this means
that our system is heartbeat-clocked, so our latency is higher than if we were
to maintain individual availability for each replica, or if we handled PUT messages
separately instead of as a batch.

Once we've got a replication quorum, we bulk send responses to the PUTs to their
originating sender.

While our heartbeat is able to send bulks of entries to replicas, we transition
to the Commit state as soon as we've seen a single PUT, so the only time we make
use of our bulk send functionality is catching up nodes that have a less up-to-date
log.

---

## Timeouts

### Election Timeout

- 300ms to 500ms
- Node that times out votes for self, swaps to CANDIDATE, requests VOTE messages from peers
  - CANDIDATE that receives a VOTE RES resets its ELECTION TIMEOUT
  - FOLLOWER that receives a VOTE REQ resets its ELECTION TIMEOUT
    - Votes affirmative if not already voted
- If a CANDIDATE has received a majority vote, convert to LEADER

### Heartbeat Timeout

- Heartbeat timeout is 1/5th of the base election timeout
- LEADER sends APPEND LOG to FOLLOWERS as heartbeat message on heartbeat timeout lapse
- FOLLOWER nodes reset their ELECTION TIMEOUT on APPEND MESSAGE receipt
  - Respond with APPEND LOG RES of their own

---

## Testing

- We primarily used the simulator provided to us to test our code. 
- The simulator provided good feedback about what our code was messing up. Such as, too many duplicate messages, or a constant flip-flopping of the leader.

---

## Challenges Faced

- We ran into issues decoding messages during a partition.

It was an incredibly odd bug where the buffer we were reading from would get half a message in it. 
This led to our decoder throwing errors as it couldn't read the junk data.
Unfortunately, we didn't log this error this error at first so it just looked like nodes suddenly stopped receiving any input. 

- Interpreting the paper was a challenge because we didn't immediately understand the reasoning behind all of their decisions.

Example: 1 indexing the log of entries.
It makes sense that you need to 1 index it to get it to work with the initialized values for commitIndex and lastApplied.
If log[] weren't 1 indexed, commitIndex and lastApplied would need to be -1, which doesn't make a ton of sense.

---

## Running

With the compiled binary, `3700kvstore`:

```
./3700kvstore <your ID> <ID of second replica> [ID3 [ID4 ...]]
```

---

## Consensus Message Data Structures

### Request Vote

Request:

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "voteRequest",
  "MID": "<a unique string>",
  "term": "<TERM ID>"
	"candidateId": "<NODE_ID>",
  "lastLogIndex": "<ENTRY_INDEX>",
  "lastLogTerm": "<TERM_ID>"
}
```

Response:

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "voteRequest",
  "MID": "<a unique string>",
  "term": "<TERM ID>",
	"voteGranted": "[true | false]"
}
```

### Append Log

Request:

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "appendRequest",
  "MID": "<a unique string>",
  "term": "<TERM ID>",
  "prevLogIndex": "<ENTRY_INDEX>",
	"prevLogTerm"  "<TERM_ID>",
	"entries"      [LOG_ENTRY],
	"leaderCommit": "<ENTRY_INDEX>"
}
```

Response:

```
{
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "appendResponse",
  "term": "<TERM ID>",
  "success": "[true | false]"
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
