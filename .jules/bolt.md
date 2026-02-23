# Bolt's Performance Journal

**Learning:** `log.Printf` in high-frequency loops (like SSE stream processing) is a major performance bottleneck due to I/O locking and string allocations.
**Action:** Remove or use conditional debug logging in hot loops. Use buffered writes where possible.

## 2026-02-19 - Slice Allocations in Writes

**Learning:** `append(b, '\n')` creates a new slice for every chunk.
**Action:** Use sequential `pw.Write(b); pw.Write(newline)` to avoid allocation.

## 2026-02-19 - Gradio Client Initialization

**Learning:** `gradio_client.Client` initialization is expensive (network call to fetch config) and blocking.
**Action:** Cache `Client` instances using `functools.lru_cache` (or global) to reuse connections across requests.
## 2026-02-15 - Hot Loop Logging and Allocation in Go Streaming
**Learning:** In Go, `log.Printf` inside a hot loop (like processing SSE chunks) causes significant I/O overhead and allocations. `bufio.Reader.ReadBytes` allocates a new slice for every line, whereas `bufio.Scanner.Bytes()` reuses the internal buffer.
**Action:** For line-based streaming (like SSE), use `bufio.Scanner` (with default 64KB limit, which is sufficient for typical JSON chunks) and avoid logging every chunk. Use sequential `Write` calls instead of `append` to avoid slice reallocation.
