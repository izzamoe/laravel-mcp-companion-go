package packages

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/models"
)

// Catalog manages Laravel package recommendations
type Catalog struct {
	data      models.PackageCatalog
	indexPath string
}

// NewCatalog creates a new package catalog
func NewCatalog(indexPath string) (*Catalog, error) {
	catalog := &Catalog{
		indexPath: indexPath,
	}

	if err := catalog.load(); err != nil {
		return nil, fmt.Errorf("failed to load catalog: %w", err)
	}

	return catalog, nil
}

// load reads the package index from JSON file
func (c *Catalog) load() error {
	data, err := os.ReadFile(c.indexPath)
	if err != nil {
		return fmt.Errorf("failed to read catalog file: %w", err)
	}

	if err := json.Unmarshal(data, &c.data); err != nil {
		return fmt.Errorf("failed to parse catalog JSON: %w", err)
	}

	return nil
}

// ListCategories returns all package categories
func (c *Catalog) ListCategories() []string {
	categories := make([]string, 0, len(c.data.Categories))
	for name := range c.data.Categories {
		categories = append(categories, name)
	}
	sort.Strings(categories)
	return categories
}

// GetCategory returns packages in a specific category
func (c *Catalog) GetCategory(categoryName string) (*models.PackageCategory, error) {
	category, ok := c.data.Categories[categoryName]
	if !ok {
		return nil, fmt.Errorf("category not found: %s", categoryName)
	}
	return &category, nil
}

// Search finds packages matching a query
func (c *Catalog) Search(query string, filters map[string]interface{}) []models.Package {
	query = strings.ToLower(query)
	var results []models.Package

	for _, category := range c.data.Categories {
		for _, pkg := range category.Packages {
			if c.matchesSearch(pkg, query, filters) {
				results = append(results, pkg)
			}
		}
	}

	// Sort by relevance (popularity score)
	sort.Slice(results, func(i, j int) bool {
		return results[i].PopularityScore > results[j].PopularityScore
	})

	return results
}

// matchesSearch checks if package matches search criteria
func (c *Catalog) matchesSearch(pkg models.Package, query string, filters map[string]interface{}) bool {
	// Text matching
	if query != "" {
		matched := false
		if strings.Contains(strings.ToLower(pkg.Name), query) {
			matched = true
		}
		if strings.Contains(strings.ToLower(pkg.Description), query) {
			matched = true
		}
		for _, tag := range pkg.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Filter by maintained status
	if maintained, ok := filters["maintained"].(bool); ok {
		if pkg.Maintained != maintained {
			return false
		}
	}

	// Filter by minimum popularity score
	if minScore, ok := filters["min_popularity"].(float64); ok {
		if pkg.PopularityScore < int(minScore) {
			return false
		}
	}

	// Filter by tags
	if tags, ok := filters["tags"].([]string); ok && len(tags) > 0 {
		hasTag := false
		for _, filterTag := range tags {
			for _, pkgTag := range pkg.Tags {
				if strings.EqualFold(pkgTag, filterTag) {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	return true
}

// Recommend returns package recommendations based on use case
func (c *Catalog) Recommend(useCase string, limit int) []models.Package {
	useCase = strings.ToLower(useCase)
	var scored []struct {
		pkg   models.Package
		score int
	}

	for _, category := range c.data.Categories {
		for _, pkg := range category.Packages {
			score := c.calculateRelevanceScore(pkg, useCase)
			if score > 0 {
				scored = append(scored, struct {
					pkg   models.Package
					score int
				}{pkg, score})
			}
		}
	}

	// Sort by relevance score
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].pkg.PopularityScore > scored[j].pkg.PopularityScore
		}
		return scored[i].score > scored[j].score
	})

	// Extract packages
	results := make([]models.Package, 0, limit)
	for i, s := range scored {
		if i >= limit {
			break
		}
		results = append(results, s.pkg)
	}

	return results
}

// calculateRelevanceScore scores a package's relevance to a use case
func (c *Catalog) calculateRelevanceScore(pkg models.Package, useCase string) int {
	score := 0

	// Check use case field
	for _, uc := range pkg.UseCase {
		if strings.Contains(strings.ToLower(uc), useCase) {
			score += 50
		}
	}

	// Check description
	if strings.Contains(strings.ToLower(pkg.Description), useCase) {
		score += 20
	}

	// Check tags
	for _, tag := range pkg.Tags {
		if strings.Contains(strings.ToLower(tag), useCase) {
			score += 15
		}
	}

	// Boost score for maintained and popular packages
	if pkg.Maintained {
		score += 10
	}
	score += pkg.PopularityScore / 10

	return score
}

// GetPackage returns details for a specific package
func (c *Catalog) GetPackage(composerName string) (*models.Package, error) {
	composerName = strings.ToLower(composerName)

	for _, category := range c.data.Categories {
		for _, pkg := range category.Packages {
			if strings.ToLower(pkg.ComposerName) == composerName {
				return &pkg, nil
			}
		}
	}

	return nil, fmt.Errorf("package not found: %s", composerName)
}

// GetFeatures returns common implementation features for a package
func (c *Catalog) GetFeatures(packageName string) []string {
	// Get package
	pkg, err := c.GetPackage(packageName)
	if err != nil {
		return []string{}
	}

	// Features are derived from use cases and tags
	features := make([]string, 0)

	// Add use cases as features
	for _, uc := range pkg.UseCase {
		features = append(features, uc)
	}

	// Add relevant tags as features
	for _, tag := range pkg.Tags {
		features = append(features, tag)
	}

	return features
}
