# 1.
# curl -o heap.pprof http://localhost:56404/debug/pprof/heap
# go tool pprof heap.pprof
# 2.
# go tool pprof -http=:9999 http://localhost:56404/debug/pprof/heap
