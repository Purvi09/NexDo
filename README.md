# NexDo 🌱

A **local-first to-do app** with a forgiving archive and a **Chrome new-tab** surface — powered by a database engine written from scratch in **Go**.

No server. No login. Works offline. Syncs peer-to-peer.

## Why
The to-do app is the surface; the engineering is the point. NexDo's core is a hand-written storage engine (an append-only log), a CRDT layer for conflict-free merges, and peer-to-peer sync — the machinery that usually hides inside Postgres or Firebase. Here it's built from the ground up.

## The idea in one line
NexDo doesn't store your task list — it stores an **append-only log of changes** (add / complete / archive / restore). Your list, and your archive, are just different *reads* of that log. Durability, undo, the archive, and conflict-free sync all fall out of that one idea.

## Roadmap
- [x] Running Go server
- [ ] Storage engine (durable, crash-safe log)
- [ ] Web UI + live updates
- [ ] Indexing & queries
- [ ] CRDT layer
- [ ] Peer-to-peer sync
- [ ] Chrome new-tab extension

## Run it
```sh
go run ./cmd/nexdo
# then open http://localhost:7777
```

## Layout
```
cmd/nexdo/   entrypoint — starts the engine + server
internal/    the engine: storage, crdt, sync, server
web/         the browser UI (todo app + new-tab page)
```

## Tech
Pure Go standard library where possible. No external database — NexDo *is* the database.
