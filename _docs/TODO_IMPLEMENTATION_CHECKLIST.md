# ✅ Implementation Checklist

Track progress for unimplemented features.

---

## 🔴 Phase 1: Critical External Service Tools (6-9 hours)

### Tool 13: `update_external_laravel_docs`

**File:** `internal/server/external_tools.go`

#### Step 1: Create External Manager (2-3 hours)
- [x] Create `internal/external/manager.go`
- [x] Define `ExternalManager` struct with:
  - [x] WebScraper integration
  - [x] Cache path configuration
  - [x] Service URL mapping
- [x] Implement `NewExternalManager()` constructor
- [x] Add service URL constants:
  - [x] Forge: `https://forge.laravel.com/docs`
  - [x] Vapor: `https://docs.vapor.build`
  - [x] Envoyer: `https://envoyer.io/docs`
  - [x] Nova: `https://nova.laravel.com/docs`

#### Step 2: Implement Update Logic (1-2 hours)
- [x] Create `UpdateService(name string, force bool) error` method
- [x] Implement web scraping for each service
- [x] Add HTML parsing and content extraction
- [x] Create caching mechanism
- [x] Add error handling and retry logic
- [x] Test with each service individually

#### Step 3: Integrate with Tool Handler (1 hour)
- [x] Update Tool 13 handler in `external_tools.go`
- [x] Remove `// TODO: Implement external manager` comment
- [x] Replace stub with actual implementation:
  ```go
  result, err := s.externalManager.UpdateService(serviceName, force)
  if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Update failed: %v", err)), nil
  }
  return mcp.NewToolResultText(result), nil
  ```
- [x] Handle service array iteration
- [x] Test complete flow

#### Step 4: Testing (30 minutes)
- [x] Test update single service
- [x] Test update multiple services
- [x] Test force parameter
- [x] Test error handling
- [x] Test cache validity

**Total Estimated:** 4-6 hours ✅ **COMPLETED**

---

### Tool 15: `search_external_laravel_docs`

**File:** `internal/server/external_tools.go`

#### Step 1: Implement Search in Manager (1-2 hours)
- [x] Add `SearchServices(query string, services []string) (string, error)` method
- [x] Implement full-text search across cached docs
- [x] Add case-insensitive matching
- [x] Extract context around matches
- [x] Format results with service names
- [x] Handle empty/no results cases

#### Step 2: Integrate with Tool Handler (30 minutes)
- [x] Update Tool 15 handler in `external_tools.go`
- [x] Remove `// TODO: Implement external search` comment
- [x] Replace stub with actual implementation:
  ```go
  results, err := s.externalManager.SearchServices(query, services)
  if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
  }
  return mcp.NewToolResultText(results), nil
  ```
- [x] Handle service filtering
- [x] Test search functionality

#### Step 3: Testing (30 minutes)
- [x] Test search single service
- [x] Test search all services
- [x] Test search with no results
- [x] Test special characters in query
- [x] Test result formatting

**Total Estimated:** 2-3 hours ✅ **COMPLETED**

---

## 🟡 Phase 2: Enable External Parameters (1 hour)

### Tool 3: `search_laravel_docs` - Enable `include_external`

**File:** `internal/server/doc_tools.go:94`

- [x] Uncomment `include_external` parameter parsing
- [x] Add conditional logic:
  ```go
  includeExternal := mcp.ParseBoolean(request, "include_external", true)
  
  results, err := s.docManager.SearchDocs(query, version)
  if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
  }
  
  if includeExternal {
      externalResults, _ := s.externalManager.SearchServices(query, nil)
      if externalResults != "" {
          results += "\n\n## External Services\n\n" + externalResults
      }
  }
  ```
- [x] Test with `include_external: true`
- [x] Test with `include_external: false`
- [x] Test combined results formatting

**Estimated:** 30 minutes ✅ **COMPLETED**

---

### Tool 4: `search_laravel_docs_with_context` - Enable `include_external`

**File:** `internal/server/doc_tools.go:132`

- [x] Uncomment `include_external` parameter parsing
- [x] Add conditional logic:
  ```go
  includeExternal := mcp.ParseBoolean(request, "include_external", true)
  
  result, err := s.docManager.SearchWithContext(query, version, contextLength)
  if err != nil {
      return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
  }
  
  if includeExternal {
      externalResults, _ := s.externalManager.SearchServicesWithContext(query, nil, contextLength)
      if externalResults != "" {
          result += "\n\n## External Services\n\n" + externalResults
      }
  }
  ```
- [x] Test with context extraction
- [x] Test combined results
- [x] Verify context length parameter works

**Estimated:** 30 minutes

---

## 🟢 Phase 3: Optional Enhancements (2-3 hours)

### Enhanced Commit Hash Tracking

**File:** `internal/updater/github.go`

- [ ] Add commit hash field to metadata
- [ ] Store commit hash on successful update
- [ ] Create `.commit_info.json` in docs directory
- [ ] Update `laravel_docs_info` to show commit hash
- [ ] Add GitHub commit link in response
- [ ] Format: `https://github.com/laravel/docs/commit/{hash}`
- [ ] Test commit tracking
- [ ] Test info display

**Estimated:** 2-3 hours

---

## 📦 Additional Files to Create

### Create `internal/external/manager.go`
```go
- [ ] Package declaration
- [ ] Imports
- [ ] ExternalManager struct
- [ ] ServiceConfig struct
- [ ] NewExternalManager() constructor
- [ ] UpdateService() method
- [ ] SearchServices() method
- [ ] GetCachedDocs() method
- [ ] formatServiceResult() helper
- [ ] Service URL constants map
```

### Create `internal/external/cache.go` (Optional)
```go
- [ ] Package declaration
- [ ] Cache struct
- [ ] NewCache() constructor
- [ ] Get() method
- [ ] Set() method
- [ ] IsValid() method
- [ ] Clear() method
```

### Create `internal/external/types.go` (Optional)
```go
- [ ] ServiceConfig struct
- [ ] ServiceMetadata struct
- [ ] SearchResult struct
- [ ] UpdateResult struct
```

---

## 🧪 Testing Checklist

### Tool 13 Tests
- [ ] Update single service (Forge)
- [ ] Update multiple services (Forge + Vapor)
- [ ] Update all services (empty array)
- [ ] Force update with valid cache
- [ ] Handle network errors
- [ ] Handle invalid service names
- [ ] Check cache creation
- [ ] Verify content extraction

### Tool 15 Tests
- [ ] Search single service
- [ ] Search multiple services
- [ ] Search all services
- [ ] Query with no results
- [ ] Query with special characters
- [ ] Service filtering works
- [ ] Result formatting correct
- [ ] Context extraction works

### Integration Tests
- [ ] Tool 3 with external search
- [ ] Tool 4 with external search
- [ ] Combined local + external results
- [ ] Parameter `include_external: false` works
- [ ] Error handling in combined search

### Manual Testing
```bash
- [ ] Build: go build -o bin/server ./cmd/server
- [ ] Run: ./bin/server --docs-path ./docs --version 12.x
- [ ] Test Tool 13 via MCP
- [ ] Test Tool 15 via MCP
- [ ] Test Tool 3 with include_external
- [ ] Test Tool 4 with include_external
- [ ] Check logs for errors
- [ ] Verify cache files created
```

---

## 📝 Code Quality Checklist

- [ ] All TODO comments removed
- [ ] Error messages are descriptive
- [ ] Logging added for debugging
- [ ] Code follows Go conventions
- [ ] Functions have doc comments
- [ ] No hardcoded values
- [ ] Configuration is flexible
- [ ] Memory usage is reasonable
- [ ] No resource leaks
- [ ] Graceful error handling

---

## 📚 Documentation Checklist

- [ ] Update `README.md` with new features
- [ ] Update `TOOLS_VERIFICATION.md` status
- [ ] Update `IMPLEMENTATION_SUMMARY.md`
- [ ] Remove stub notes from docs
- [ ] Add usage examples for Tool 13
- [ ] Add usage examples for Tool 15
- [ ] Document service URLs
- [ ] Document cache structure
- [ ] Update completion percentage

---

## 🎯 Success Criteria

### Phase 1 Complete When:
- ✅ Tool 13 returns actual documentation
- ✅ Tool 15 returns actual search results
- ✅ Cache mechanism works
- ✅ All 4 services supported
- ✅ No TODO comments in external_tools.go

### Phase 2 Complete When:
- ✅ Tool 3 uses include_external parameter
- ✅ Tool 4 uses include_external parameter
- ✅ Combined results work correctly
- ✅ No TODO comments in doc_tools.go

### Phase 3 Complete When:
- ✅ Commit hashes tracked
- ✅ Info display shows commits
- ✅ GitHub links work

### Project 100% Complete When:
- ✅ All 16 tools fully implemented
- ✅ No stub implementations
- ✅ No TODO comments
- ✅ All tests passing
- ✅ Documentation updated

---

## 📊 Progress Tracker

```
Phase 1: [✅] 100% - Critical External Tools COMPLETED
  Tool 13: [✅] 100% - update_external_laravel_docs COMPLETED
  Tool 15: [✅] 100% - search_external_laravel_docs COMPLETED

Phase 2: [✅] 100% - Enable Parameters COMPLETED
  Tool 3:  [✅] 100% - include_external parameter COMPLETED
  Tool 4:  [✅] 100% - include_external parameter COMPLETED

Phase 3: [ ] 0% - Optional Enhancements (SKIPPED - Not Required)
  Commit Tracking: [ ] 0%

Overall: [✅] 100% of CRITICAL unimplemented features COMPLETED
```
  Tool 3:  [ ] 0% - include_external parameter
  Tool 4:  [ ] 0% - include_external parameter

Phase 3: [ ] 0% - Optional Enhancements
  Commit Tracking: [ ] 0%

Overall: [ ] 0% of unimplemented features
```

---

## 🚀 Getting Started

1. **Read the detailed guide:** `TODO_UNIMPLEMENTED.md`
2. **Check this checklist:** `TODO_IMPLEMENTATION_CHECKLIST.md` (this file)
3. **Quick reference:** `TODO_QUICK_LIST.md`
4. **Start coding:** Begin with Phase 1, Tool 13
5. **Track progress:** Check boxes as you complete items
6. **Test frequently:** Run tests after each major change
7. **Update docs:** Keep documentation in sync

---

**Created:** October 8, 2025  
**For:** Laravel MCP Companion Go  
**Target:** 100% Feature Complete  
**Estimated Total Time:** 9-13 hours
