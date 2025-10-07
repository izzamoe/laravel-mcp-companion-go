package packages

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/models"
)

// FormatCategoriesList formats categories for display
func FormatCategoriesList(categories []string) (string, error) {
	output := strings.Builder{}
	output.WriteString("📦 Laravel Package Categories\n\n")

	for i, cat := range categories {
		output.WriteString(fmt.Sprintf("%d. %s\n", i+1, cat))
	}

	return output.String(), nil
}

// FormatCategoryPackages formats packages in a category
func FormatCategoryPackages(category *models.PackageCategory) (string, error) {
	output := strings.Builder{}
	output.WriteString(fmt.Sprintf("📦 %s\n", category.Name))
	output.WriteString(fmt.Sprintf("%s\n\n", category.Description))

	for i, pkg := range category.Packages {
		output.WriteString(fmt.Sprintf("%d. **%s** (%s)\n", i+1, pkg.Name, pkg.ComposerName))
		output.WriteString(fmt.Sprintf("   %s\n", pkg.Description))
		output.WriteString(fmt.Sprintf("   📊 Score: %d/100 | ", pkg.PopularityScore))

		if pkg.Maintained {
			output.WriteString("✅ Maintained")
		} else {
			output.WriteString("⚠️ Not actively maintained")
		}
		output.WriteString("\n\n")
	}

	return output.String(), nil
}

// FormatPackageList formats a list of packages
func FormatPackageList(packages []models.Package, title string) (string, error) {
	if len(packages) == 0 {
		return "No packages found matching your criteria.", nil
	}

	output := strings.Builder{}
	output.WriteString(fmt.Sprintf("📦 %s\n\n", title))

	for i, pkg := range packages {
		output.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, pkg.Name))
		output.WriteString(fmt.Sprintf("   📦 Composer: `%s`\n", pkg.ComposerName))
		output.WriteString(fmt.Sprintf("   📝 %s\n", pkg.Description))
		output.WriteString(fmt.Sprintf("   📊 Popularity: %d/100\n", pkg.PopularityScore))
		output.WriteString(fmt.Sprintf("   🏷️  Tags: %s\n", strings.Join(pkg.Tags, ", ")))

		if len(pkg.UseCase) > 0 {
			output.WriteString(fmt.Sprintf("   💡 Use cases: %s\n", strings.Join(pkg.UseCase, ", ")))
		}

		if len(pkg.Alternatives) > 0 {
			output.WriteString(fmt.Sprintf("   🔄 Alternatives: %s\n", strings.Join(pkg.Alternatives, ", ")))
		}

		output.WriteString("\n")
	}

	return output.String(), nil
}

// FormatPackageDetails formats detailed package information
func FormatPackageDetails(pkg *models.Package) (string, error) {
	output := strings.Builder{}

	output.WriteString(fmt.Sprintf("# %s\n\n", pkg.Name))
	output.WriteString(fmt.Sprintf("**Composer Package:** `%s`\n\n", pkg.ComposerName))
	output.WriteString(fmt.Sprintf("## Description\n%s\n\n", pkg.Description))

	output.WriteString("## Statistics\n")
	output.WriteString(fmt.Sprintf("- **Popularity Score:** %d/100\n", pkg.PopularityScore))
	output.WriteString(fmt.Sprintf("- **Minimum Laravel Version:** %s\n", pkg.MinLaravelVersion))

	status := "✅ Actively Maintained"
	if !pkg.Maintained {
		status = "⚠️ Not Actively Maintained"
	}
	output.WriteString(fmt.Sprintf("- **Status:** %s\n\n", status))

	if len(pkg.Tags) > 0 {
		output.WriteString(fmt.Sprintf("## Tags\n%s\n\n", strings.Join(pkg.Tags, ", ")))
	}

	if len(pkg.UseCase) > 0 {
		output.WriteString("## Use Cases\n")
		for _, uc := range pkg.UseCase {
			output.WriteString(fmt.Sprintf("- %s\n", uc))
		}
		output.WriteString("\n")
	}

	if len(pkg.Alternatives) > 0 {
		output.WriteString("## Alternatives\n")
		for _, alt := range pkg.Alternatives {
			output.WriteString(fmt.Sprintf("- %s\n", alt))
		}
		output.WriteString("\n")
	}

	output.WriteString("## Installation\n")
	output.WriteString(fmt.Sprintf("```bash\ncomposer require %s\n```\n", pkg.ComposerName))

	return output.String(), nil
}

// FormatPackageJSON formats package as JSON
func FormatPackageJSON(pkg *models.Package) (string, error) {
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal package: %w", err)
	}
	return string(data), nil
}
