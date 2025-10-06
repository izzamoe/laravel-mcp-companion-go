package models

import "time"

// DocMetadata represents metadata about documentation
type DocMetadata struct {
	Version   string    `json:"version"`
	CommitSHA string    `json:"commit_sha"`
	SyncTime  time.Time `json:"sync_time"`
	FileCount int       `json:"file_count"`
}

// SupportedVersions list of Laravel versions
var SupportedVersions = []string{
	"12.x", "11.x", "10.x", "9.x",
	"8.x", "7.x", "6.x",
}

// DefaultVersion is the latest stable version
const DefaultVersion = "12.x"

// CategoryMappings maps categories to their documentation files
var CategoryMappings = map[string][]string{
	"getting-started": {
		"installation.md",
		"configuration.md",
		"structure.md",
		"deployment.md",
	},
	"architecture": {
		"lifecycle.md",
		"container.md",
		"providers.md",
		"facades.md",
	},
	"basics": {
		"routing.md",
		"middleware.md",
		"controllers.md",
		"requests.md",
		"responses.md",
		"views.md",
	},
	"frontend": {
		"blade.md",
		"vite.md",
		"frontend.md",
	},
	"security": {
		"authentication.md",
		"authorization.md",
		"verification.md",
		"encryption.md",
		"hashing.md",
		"passwords.md",
	},
	"database": {
		"database.md",
		"queries.md",
		"pagination.md",
		"migrations.md",
		"seeding.md",
		"redis.md",
	},
	"eloquent": {
		"eloquent.md",
		"eloquent-relationships.md",
		"eloquent-collections.md",
		"eloquent-mutators.md",
		"eloquent-resources.md",
		"eloquent-serialization.md",
	},
	"testing": {
		"testing.md",
		"http-tests.md",
		"console-tests.md",
		"dusk.md",
		"database-testing.md",
		"mocking.md",
	},
}

// Package represents a Laravel package
type Package struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	ComposerName      string   `json:"composer_name"`
	UseCase           []string `json:"use_case"`
	Alternatives      []string `json:"alternatives"`
	MinLaravelVersion string   `json:"min_laravel_version"`
	Tags              []string `json:"tags"`
	PopularityScore   int      `json:"popularity_score"`
	Maintained        bool     `json:"maintained"`
}

// PackageCategory represents a category of packages
type PackageCategory struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Packages    []Package `json:"packages"`
}

// PackageCatalog represents the full package catalog
type PackageCatalog struct {
	Categories map[string]PackageCategory `json:"categories"`
}
