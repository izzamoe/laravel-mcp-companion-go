# 📝 Quick TODO List

> **TL;DR:** 14/16 tools complete, 2 tools need full implementation

---

## 🔴 Critical (Must Implement)

### 1. Tool 13: `update_external_laravel_docs` 
**File:** `internal/server/external_tools.go:79-100`  
**Status:** Returns stub message only  
**Effort:** 4-6 hours  

**Missing:**
```
- Web scraping for Forge/Vapor/Envoyer/Nova
- External documentation caching
- Update mechanism
```

**TODO Comment:**
```go
// TODO: Implement external manager
```

---

### 2. Tool 15: `search_external_laravel_docs`
**File:** `internal/server/external_tools.go:136-165`  
**Status:** Returns stub message only  
**Effort:** 2-3 hours (depends on Tool 13)

**Missing:**
```
- Search across external service docs
- Integration with external manager
```

**TODO Comment:**
```go
// TODO: Implement external search
```

---

## 🟡 Medium Priority (After Critical)

### 3. Tool 3: Enable `include_external` parameter
**File:** `internal/server/doc_tools.go:94`  
**Effort:** 30 minutes  
**Depends on:** Tool 15

**TODO Comment:**
```go
// includeExternal := mcp.ParseBoolean(request, "include_external", true) // TODO: implement external search
```

---

### 4. Tool 4: Enable `include_external` parameter
**File:** `internal/server/doc_tools.go:132`  
**Effort:** 30 minutes  
**Depends on:** Tool 15

**TODO Comment:**
```go
// includeExternal := mcp.ParseBoolean(request, "include_external", true) // TODO: implement external search
```

---

## 🟢 Nice to Have

### 5. Enhanced Commit Hash Tracking
**File:** `internal/updater/github.go`  
**Effort:** 2-3 hours  
**Status:** Basic implementation exists, can be enhanced

---

## 📊 Summary

| Priority | Item | File | Effort | Status |
|----------|------|------|--------|--------|
| 🔴 High | Tool 13 | `external_tools.go` | 4-6h | Stub |
| 🔴 High | Tool 15 | `external_tools.go` | 2-3h | Stub |
| 🟡 Medium | Tool 3 param | `doc_tools.go` | 30m | Commented |
| 🟡 Medium | Tool 4 param | `doc_tools.go` | 30m | Commented |
| 🟢 Low | Commit tracking | `github.go` | 2-3h | Optional |

**Total Critical Work:** 6-9 hours  
**Total All Work:** 9-13 hours

---

## 🎯 Implementation Order

```
1. Tool 13: update_external_laravel_docs (4-6h)
   └─> Create internal/external/manager.go
   └─> Implement web scraping for 4 services
   └─> Add caching mechanism

2. Tool 15: search_external_laravel_docs (2-3h)
   └─> Use ExternalManager from Tool 13
   └─> Implement search logic

3. Tool 3 & 4: Enable include_external (1h)
   └─> Uncomment and implement parameter
   └─> Combine local + external results

4. Optional: Enhanced commit tracking (2-3h)
   └─> Add commit metadata storage
   └─> Update info display
```

---

## 🛠️ Files to Create

- [ ] `internal/external/manager.go` - External doc manager
- [ ] `internal/external/cache.go` (optional) - Caching logic
- [ ] `internal/external/types.go` (optional) - Type definitions

---

## 📝 Files to Modify

- [ ] `internal/server/external_tools.go` - Implement Tool 13 & 15
- [ ] `internal/server/doc_tools.go` - Enable external parameters
- [ ] `internal/updater/github.go` (optional) - Enhanced tracking

---

## ✅ Already Available (No Work Needed)

- ✅ `internal/external/web.go` - WebScraper ready
- ✅ HTTP client with proper timeout
- ✅ Content size limiting
- ✅ HTML cleanup utilities
- ✅ 14 out of 16 tools fully working

---

## 🚀 Quick Start

```bash
# 1. Check current TODOs in code
grep -r "TODO" internal/

# 2. Start with Tool 13
vim internal/external/manager.go

# 3. Test implementation
go build -o bin/server ./cmd/server
./bin/server --docs-path ./docs --version 12.x

# 4. Test Tool 13
./test_tools.sh  # or manual test
```

---

## 📌 Key Points

1. **Core functionality (14 tools) is 100% complete** ✅
2. **External service tools (2 tools) are stubs** ⚠️
3. **Project is production-ready for main use cases** ✅
4. **External features are nice-to-have, not critical** 💡
5. **Estimated 6-9 hours for full completion** ⏱️

---

**See:** `TODO_UNIMPLEMENTED.md` for detailed implementation guide

**Last Updated:** October 8, 2025
