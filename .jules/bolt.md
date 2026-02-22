# Bolt's Performance Journal

## 2026-02-15 - Hot Loop Logging and Allocation in Go Streaming
**Learning:** In Go, `log.Printf` inside a hot loop (like processing SSE chunks) causes significant I/O overhead and allocations. `bufio.Reader.ReadBytes` allocates a new slice for every line, whereas `bufio.Scanner.Bytes()` reuses the internal buffer.
**Action:** For line-based streaming (like SSE), use `bufio.Scanner` (with default 64KB limit, which is sufficient for typical JSON chunks) and avoid logging every chunk. Use sequential `Write` calls instead of `append` to avoid slice reallocation.
