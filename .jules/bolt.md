## 2026-02-19 - Logging in Hot Loops

**Learning:** `log.Printf` in high-frequency loops (like SSE stream processing) is a major performance bottleneck due to I/O locking and string allocations.
**Action:** Remove or use conditional debug logging in hot loops. Use buffered writes where possible.

## 2026-02-19 - Slice Allocations in Writes

**Learning:** `append(b, '\n')` creates a new slice for every chunk.
**Action:** Use sequential `pw.Write(b); pw.Write(newline)` to avoid allocation.

## 2026-02-19 - Gradio Client Initialization

**Learning:** `gradio_client.Client` initialization is expensive (network call to fetch config) and blocking.
**Action:** Cache `Client` instances using `functools.lru_cache` (or global) to reuse connections across requests.
