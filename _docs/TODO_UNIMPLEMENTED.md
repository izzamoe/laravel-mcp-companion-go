# 📋 TODO: Unimplemented Features

> **Last Updated:** October 8, 2025  
> **Project:** Laravel MCP Companion (Go)  
> **Status:** 14/16 tools fully implemented, 2 tools stub only

---

## 🔴 High Priority - Stub Implementations

### 1. Tool 13: `update_external_laravel_docs`

**Status:** ⚠️ Stub Implementation Only  
**File:** `internal/server/external_tools.go` (Lines 79-100)  
**Function:** Updates documentation for external Laravel services

#### Current State:
```go
// TODO: Implement external manager
_ = services
_ = force

return mcp.NewToolResultText("External service documentation update is not yet implemented. Available services: forge, vapor, envoyer, nova"), nil
```

#### What's Missing:
- [ ] Web scraping implementation for Laravel Forge docs
- [ ] Web scraping implementation for Laravel Vapor docs  
- [ ] Web scraping implementation for Laravel Envoyer docs
- [ ] Web scraping implementation for Laravel Nova docs
- [ ] Caching mechanism for external documentation
- [ ] Update manager for external services
- [ ] Force update logic
- [ ] Service-specific URL mapping
- [ ] Content extraction and formatting
- [ ] Error handling for failed downloads

#### Implementation Requirements:
1. **Create External Manager:**
   - File: `internal/external/manager.go`
   - Struct: `ExternalManager`
   - Methods: `UpdateService()`, `GetServiceDocs()`, `CacheService()`

2. **Service URLs:**
   - Forge: `https://forge.laravel.com/docs`
   - Vapor: `https://docs.vapor.build`
   - Envoyer: `https://envoyer.io/docs`
   - Nova: `https://nova.laravel.com/docs`

3. **Use Existing WebScraper:**
   - Already available: `internal/external/web.go`
   - Method: `FetchResource()`

#### Estimated Effort:
- **Time:** 4-6 hours
- **Complexity:** Medium
- **Files to Create:** 1-2 new files
- **Files to Modify:** 1 file (`external_tools.go`)

---

### 2. Tool 15: `search_external_laravel_docs`

**Status:** ⚠️ Stub Implementation Only  
**File:** `internal/server/external_tools.go` (Lines 136-165)  
**Function:** Searches through external Laravel service documentation

#### Current State:
```go
// TODO: Implement external search
_ = services

return mcp.NewToolResultText(fmt.Sprintf("External service search for '%s' is not yet implemented. Try using the main documentation search or visit the service websites directly.", query)), nil
```

#### What's Missing:
- [ ] Search implementation across cached external docs
- [ ] Integration with External Manager (from Tool 13)
- [ ] Multi-service search capability
- [ ] Result ranking and relevance scoring
- [ ] Context extraction for search results
- [ ] Service filtering logic

#### Implementation Requirements:
1. **Depends on:** Tool 13 implementation (External Manager)
2. **Search Logic:**
   - Full-text search across cached external docs
   - Case-insensitive matching
   - Result formatting with context
   - Service name in results

3. **Integration:**
   - Use `ExternalManager.SearchServices(query, services)`
   - Similar to `docManager.SearchDocs()` pattern

#### Estimated Effort:
- **Time:** 2-3 hours (after Tool 13 is done)
- **Complexity:** Low-Medium
- **Dependency:** Requires Tool 13 completion
- **Files to Modify:** 1 file (`external_tools.go`)

---

## 🟡 Medium Priority - Incomplete Parameters

### 3. Tool 3: `search_laravel_docs` - Missing Parameter

**Status:** ⚠️ Parameter Not Used  
**File:** `internal/server/doc_tools.go` (Line 94)  
**Function:** Search Laravel documentation

#### Current State:
```go
// includeExternal := mcp.ParseBoolean(request, "include_external", true) // TODO: implement external search
```

#### What's Missing:
- [ ] Use `include_external` parameter
- [ ] Combine local + external search when `true`
- [ ] Integration with Tool 15 search functionality

#### Implementation Requirements:
1. **After Tool 15 is implemented:**
   ```go
   includeExternal := mcp.ParseBoolean(request, "include_external", true)
   
   if includeExternal {
       // Also search external services
       externalResults := s.externalManager.SearchAll(query)
       // Combine with local results
   }
   ```

#### Estimated Effort:
- **Time:** 30 minutes
- **Complexity:** Low
- **Dependency:** Requires Tool 15 completion

---

### 4. Tool 4: `search_laravel_docs_with_context` - Missing Parameter

**Status:** ⚠️ Parameter Not Used  
**File:** `internal/server/doc_tools.go` (Line 132)  
**Function:** Search Laravel docs with context

#### Current State:
```go
// includeExternal := mcp.ParseBoolean(request, "include_external", true) // TODO: implement external search
```

#### What's Missing:
- [ ] Use `include_external` parameter
- [ ] Combine local + external search results with context
- [ ] Integration with Tool 15 search functionality

#### Implementation Requirements:
1. **After Tool 15 is implemented:**
   ```go
   includeExternal := mcp.ParseBoolean(request, "include_external", true)
   
   if includeExternal {
       // Also search external services with context
       externalResults := s.externalManager.SearchWithContext(query, contextLength)
       // Combine with local results
   }
   ```

#### Estimated Effort:
- **Time:** 30 minutes
- **Complexity:** Low
- **Dependency:** Requires Tool 15 completion

---

## 🟢 Low Priority - Nice to Have

### 5. Enhanced Commit Hash Tracking

**Status:** ✅ Basic Implementation Exists  
**File:** `internal/updater/github.go`  
**Function:** Track GitHub commit hashes for documentation updates

#### Current State:
- Basic update tracking is implemented
- Commit information is available from GitHub API

#### What Could Be Enhanced:
- [ ] Store detailed commit hash in metadata
- [ ] Show commit history for docs
- [ ] Compare versions by commit
- [ ] Link to specific GitHub commits in responses
- [ ] Track when each file was last updated

#### Implementation Requirements:
1. **Add commit tracking to DocManager:**
   - Store commit hash per version
   - Display in `laravel_docs_info` tool
   - Add metadata file: `.commit_info.json`

2. **Update Response Format:**
   ```
   Last Updated: 2024-10-08
   Commit: abc123def456
   GitHub: https://github.com/laravel/docs/commit/abc123def456
   ```

#### Estimated Effort:
- **Time:** 2-3 hours
- **Complexity:** Low
- **Priority:** Nice to have

---

## 📊 Implementation Summary

### Completion Status:

| Category | Status | Count |
|----------|--------|-------|
| ✅ Fully Implemented | 100% | 14 tools |
| ⚠️ Stub Only | Partial | 2 tools |
| 🔴 Missing Parameters | Inactive | 2 params |
| 🟢 Enhancements | Optional | 1 feature |

### Priority Roadmap:

```
Phase 1 (Critical): 🔴
├── Tool 13: update_external_laravel_docs (4-6h)
└── Tool 15: search_external_laravel_docs (2-3h)
    Total: 6-9 hours

Phase 2 (Integration): 🟡
├── Tool 3: Enable include_external parameter (30m)
└── Tool 4: Enable include_external parameter (30m)
    Total: 1 hour

Phase 3 (Optional): 🟢
└── Enhanced commit tracking (2-3h)
    Total: 2-3 hours

TOTAL ESTIMATED TIME: 9-13 hours
```

---

## 🛠️ Technical Details

### Files That Need Creation:
1. `internal/external/manager.go` - External documentation manager
2. `internal/external/cache.go` (optional) - Caching for external docs
3. `internal/external/types.go` (optional) - Types for external services

### Files That Need Modification:
1. `internal/server/external_tools.go` - Implement Tool 13 & 15
2. `internal/server/doc_tools.go` - Enable external search parameters
3. `internal/updater/github.go` (optional) - Enhanced commit tracking

### Dependencies Already Available:
- ✅ `internal/external/web.go` - WebScraper ready to use
- ✅ HTTP client with timeout and redirect handling
- ✅ Content size limiting (5MB max)
- ✅ HTML cleanup utilities

---

## 🎯 Implementation Guide

### For Tool 13 & 15:

```go
// 1. Create ExternalManager
type ExternalManager struct {
    scraper    *WebScraper
    cachePath  string
    services   map[string]ServiceConfig
}

// 2. Service URLs
var ServiceURLs = map[string]string{
    "forge":   "https://forge.laravel.com/docs",
    "vapor":   "https://docs.vapor.build",
    "envoyer": "https://envoyer.io/docs",
    "nova":    "https://nova.laravel.com/docs",
}

// 3. Implement methods
func (m *ExternalManager) UpdateService(name string, force bool) error
func (m *ExternalManager) SearchServices(query string, services []string) (string, error)
func (m *ExternalManager) GetCachedDocs(service string) (string, error)
```

### Testing Commands:

```bash
# After implementation, test with:
./bin/server --docs-path ./docs --version 12.x

# Test Tool 13
echo '{"method":"tools/call","params":{"name":"update_external_laravel_docs","arguments":{"services":["forge"],"force":true}}}' | ./bin/server

# Test Tool 15
echo '{"method":"tools/call","params":{"name":"search_external_laravel_docs","arguments":{"query":"deployment","services":["forge","vapor"]}}}' | ./bin/server
```

---

## 📝 Notes

### Why These Features Are Stub:
1. **External service docs require web scraping** - More complex than GitHub API
2. **Different doc structures** - Each service has unique HTML structure
3. **Rate limiting concerns** - Need to be careful with external requests
4. **Maintenance overhead** - External sites may change structure

### Current Workaround:
- Tool 14 (`list_laravel_services`) provides service information
- Tool 16 (`get_laravel_service_info`) gives detailed service info
- Users can visit service websites directly
- Main Laravel docs (Tools 1-6) fully functional

### Production Status:
✅ **Project is production-ready for:**
- Complete Laravel documentation access (11.x, 12.x)
- Package recommendations (100+ packages)
- Documentation updates from GitHub
- Service information and discovery

⚠️ **Not ready for:**
- Automated external service documentation updates
- Searching through Forge/Vapor/Envoyer/Nova docs
- Integrated external service documentation access

---

## 🚀 Quick Start for Contributors

### To Implement Tool 13:

1. Create `internal/external/manager.go`
2. Use existing `WebScraper` from `internal/external/web.go`
3. Implement caching mechanism
4. Update `internal/server/external_tools.go` Tool 13 handler
5. Test with each service

### To Implement Tool 15:

1. Depend on Tool 13 completion
2. Add search method to `ExternalManager`
3. Update `internal/server/external_tools.go` Tool 15 handler
4. Test search across services

### To Enable External Search Parameters:

1. Depend on Tool 15 completion
2. Update Tool 3 handler in `internal/server/doc_tools.go`
3. Update Tool 4 handler in `internal/server/doc_tools.go`
4. Combine local + external results

---

## 📞 Contact & References

- **Main Spec:** `MCP_GO_IMPLEMENTATION_COMPLETE.md`
- **Verification:** `TOOLS_VERIFICATION.md`
- **Implementation:** `IMPLEMENTATION_SUMMARY.md`
- **Test Script:** `test_tools.sh`

---

**Generated:** October 8, 2025  
**For:** Laravel MCP Companion Go Implementation  
**Status:** Ready for contribution 🎉
