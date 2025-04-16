package project

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Package yapısı
type Package struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

// Project yapısı
type Project struct {
	Name     string
	Path     string
	Packages []Package
}

// TemplateType yapısı
type TemplateType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateProject, yeni bir Flutter projesi oluşturur
func CreateProject(projectName string, types []string) error {
	// Mevcut dizini al
	currentDir := os.Getenv("PWD")
	if currentDir == "" {
		return fmt.Errorf("mevcut dizin alınamadı")
	}

	// Executable dizinini al
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}
	execDir := filepath.Dir(execPath)

	// Flutter projesi oluştur
	fmt.Printf("ℹ️ Flutter projesi oluşturuluyor: %s\n", projectName)

	// Önce mevcut dizine geç
	if err := os.Chdir(currentDir); err != nil {
		return fmt.Errorf("mevcut dizine geçilemedi: %v", err)
	}

	cmd := exec.Command("flutter", "create", projectName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Flutter projesi oluşturulamadı: %v", err)
	}
	fmt.Printf("✅ Flutter projesi başarıyla oluşturuldu: %s\n", projectName)

	// Proje dizinine geç
	projectPath := filepath.Join(currentDir, projectName)
	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("proje dizinine geçilemedi: %v", err)
	}

	// Paketleri al ve filtrele
	allPackages, err := GetPackages()
	if err != nil {
		return fmt.Errorf("paketler okunamadı: %v", err)
	}

	// Seçilen type'lara göre paketleri filtrele
	var filteredPackages []Package
	for _, pkg := range allPackages {
		// ALL type'ı olan paketleri her zaman ekle
		if contains(pkg.Types, "ALL") {
			filteredPackages = append(filteredPackages, pkg)
			continue
		}

		// Seçilen type'lardan herhangi biri paketin type'larında varsa ekle
		for _, selectedType := range types {
			if contains(pkg.Types, selectedType) {
				filteredPackages = append(filteredPackages, pkg)
				break
			}
		}
	}

	// Gerekli paketleri ekle
	fmt.Printf("ℹ️ Seçilen paketler ekleniyor...\n")
	for _, pkg := range filteredPackages {
		fmt.Printf("  📦 %s paketi ekleniyor...\n", pkg.Name)
		cmd := exec.Command("flutter", "pub", "add", pkg.Name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("paket eklenemedi %s: %v", pkg.Name, err)
		}
		fmt.Printf("  ✅ %s paketi başarıyla eklendi\n", pkg.Name)
	}

	// Template dosyalarını oku ve oluştur
	fmt.Printf("ℹ️ Template dosyaları oluşturuluyor...\n")
	templateDir := filepath.Join(execDir, "template_util", "templates")
	files, err := os.ReadDir(templateDir)
	if err != nil {
		return fmt.Errorf("template klasörü okunamadı: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			templatePath := filepath.Join(templateDir, file.Name())
			data, err := os.ReadFile(templatePath)
			if err != nil {
				return fmt.Errorf("template dosyası okunamadı %s: %v", file.Name(), err)
			}

			var template struct {
				Types []string `json:"types"`
			}
			if err := json.Unmarshal(data, &template); err != nil {
				return fmt.Errorf("template JSON parse hatası %s: %v", file.Name(), err)
			}

			// Template'in type'larından herhangi biri seçilen type'larda varsa işle
			shouldProcess := false
			for _, templateType := range template.Types {
				if contains(types, templateType) || templateType == "ALL" {
					shouldProcess = true
					break
				}
			}

			if shouldProcess {
				fmt.Printf("  📄 %s template dosyası işleniyor...\n", file.Name())
				if err := processTemplate(templatePath, projectName); err != nil {
					return fmt.Errorf("template işlenemedi %s: %v", file.Name(), err)
				}
				fmt.Printf("  ✅ %s template dosyası başarıyla oluşturuldu\n", file.Name())
			}
		}
	}

	// Mevcut dizine geri dön
	if err := os.Chdir(currentDir); err != nil {
		return fmt.Errorf("mevcut dizine dönülemedi: %v", err)
	}

	fmt.Printf("🎉 Proje başarıyla oluşturuldu ve yapılandırıldı!\n")
	fmt.Printf("📁 Proje dizini: %s\n", projectPath)
	return nil
}

// readPackages, packages.json dosyasını okur
func readPackages() ([]Package, error) {
	data, err := os.ReadFile("template_util/packages.json")
	if err != nil {
		return nil, err
	}

	var packages []Package
	if err := json.Unmarshal(data, &packages); err != nil {
		return nil, err
	}

	return packages, nil
}

// processTemplate, bir template dosyasını işler
func processTemplate(templatePath string, projectName string) error {
	// Template dosyasını oku
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template dosyası okunamadı: %v", err)
	}

	var template struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("template JSON parse hatası: %v", err)
	}

	// {FLUTTER_ASSIST} yerine proje ismini koy
	content := strings.ReplaceAll(template.Content, "{FLUTTER_ASSIST}", projectName)

	fmt.Println("filePath WARNING:", template.Path)
	// Klasörü oluştur
	dir := filepath.Dir(template.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("klasör oluşturulamadı: %v", err)
	}

	// Dosyayı oluştur
	return os.WriteFile(template.Path, []byte(content), 0644)
}

func (p *Project) AddPackage(pkg Package) error {
	packages, err := p.readPackages()
	if err != nil {
		return err
	}

	// Check if package already exists
	for _, existingPkg := range packages {
		if existingPkg.Name == pkg.Name {
			return fmt.Errorf("package %s already exists", pkg.Name)
		}
	}

	packages = append(packages, pkg)
	return p.writePackages(packages)
}

func (p *Project) GetPackagesByTypes(types []string) ([]Package, error) {
	packages, err := p.readPackages()
	if err != nil {
		return nil, err
	}

	var filteredPackages []Package
	for _, pkg := range packages {
		// Check if package is for all types
		if contains(pkg.Types, "ALL") {
			filteredPackages = append(filteredPackages, pkg)
			continue
		}

		// Check if package matches any of the selected types
		for _, t := range types {
			if contains(pkg.Types, t) {
				filteredPackages = append(filteredPackages, pkg)
				break
			}
		}
	}

	return filteredPackages, nil
}

func (p *Project) readPackages() ([]Package, error) {
	filePath := filepath.Join(p.Path, "packages.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packages []Package
	if err := json.Unmarshal(data, &packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func (p *Project) writePackages(packages []Package) error {
	filePath := filepath.Join(p.Path, "packages.json")
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// contains checks if a string exists in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// hasAnyType, verilen paket tiplerinden herhangi biri seçili tiplerde varsa true döner
func hasAnyType(pkgTypes []string, selectedTypes []string) bool {
	for _, pkgType := range pkgTypes {
		for _, selectedType := range selectedTypes {
			if pkgType == selectedType {
				return true
			}
		}
	}
	return false
}

// GetTemplateTypes, template_util/types.json dosyasındaki type'ları döndürür
func GetTemplateTypes() ([]TemplateType, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	typesPath := filepath.Join(execDir, "template_util", "template_for.json")

	// Dosya yoksa boş liste döndür
	if _, err := os.Stat(typesPath); os.IsNotExist(err) {
		return []TemplateType{}, nil
	}

	data, err := os.ReadFile(typesPath)
	if err != nil {
		return nil, fmt.Errorf("type'lar dosyası okunamadı: %v", err)
	}

	// JSON yapısını kontrol et
	var result struct {
		Types []TemplateType `json:"types"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("type'lar JSON parse hatası: %v", err)
	}

	return result.Types, nil
}

// GetAllPackages, tüm paketleri döndürür
func GetAllPackages() ([]Package, error) {
	return readPackages()
}

// GetAllTemplates, tüm template dosyalarını döndürür
func GetAllTemplates() ([]string, error) {
	files, err := os.ReadDir("template_util/templates")
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			templates = append(templates, file.Name())
		}
	}

	return templates, nil
}

// DeleteProject, bir projeyi siler
func DeleteProject(projectName string) error {
	return os.RemoveAll(projectName)
}

// DeleteTemplate, bir template'i siler
func DeleteTemplate(templateName string) error {
	return os.Remove(fmt.Sprintf("template_util/templates/%s.json", templateName))
}

// GetAllProjects, tüm projeleri döndürür
func GetAllProjects() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var projects []string
	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			projects = append(projects, file.Name())
		}
	}

	return projects, nil
}

// DeleteType, bir type'ı siler
func DeleteType(typeName string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	typesPath := filepath.Join(execDir, "template_util", "template_for.json")

	types, err := GetTemplateTypes()
	if err != nil {
		return err
	}

	// Type'ı bul ve sil
	var newTypes []TemplateType
	for _, t := range types {
		if t.Name != typeName {
			newTypes = append(newTypes, t)
		}
	}

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(newTypes, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	if err := os.WriteFile(typesPath, data, 0644); err != nil {
		return fmt.Errorf("type'lar dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// CreateTemplate, yeni bir template oluşturur
func CreateTemplate(templateName string, types []string) error {
	// Template dosyasını oluştur
	template := struct {
		Types []string `json:"types"`
		Files []string `json:"files"`
	}{
		Types: types,
		Files: []string{},
	}

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("template_util/templates/%s.json", templateName), data, 0644)
}

// DeletePackage, belirtilen paketi siler
func DeletePackage(name string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	packagesPath := filepath.Join(execDir, "template_util", "packages.json")

	// Mevcut paketleri al
	packages, err := GetPackages()
	if err != nil {
		return err
	}

	// Paketi bul ve sil
	var newPackages []Package
	for _, pkg := range packages {
		if pkg.Name != name {
			newPackages = append(newPackages, pkg)
		}
	}

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(newPackages, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	if err := os.WriteFile(packagesPath, data, 0644); err != nil {
		return fmt.Errorf("paketler dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// AddPackage, yeni bir paket ekler
func AddPackage(packageName string, types []string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	packagesPath := filepath.Join(execDir, "template_util", "packages.json")

	// Mevcut paketleri al
	packages, err := GetPackages()
	if err != nil {
		return err
	}

	// Paket zaten var mı kontrol et
	for _, pkg := range packages {
		if pkg.Name == packageName {
			return fmt.Errorf("paket zaten mevcut: %s", packageName)
		}
	}

	// Yeni paketi ekle
	packages = append(packages, Package{
		Name:  packageName,
		Types: types,
	})

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	// Template util klasörünü oluştur
	templateUtilDir := filepath.Join(execDir, "template_util")
	if err := os.MkdirAll(templateUtilDir, 0755); err != nil {
		return fmt.Errorf("template util klasörü oluşturulamadı: %v", err)
	}

	if err := os.WriteFile(packagesPath, data, 0644); err != nil {
		return fmt.Errorf("paketler dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// AddType, yeni bir type ekler
func AddType(typeName string, description string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	typesPath := filepath.Join(execDir, "template_util", "template_for.json")

	types, err := GetTemplateTypes()
	if err != nil {
		return err
	}

	// Type zaten var mı kontrol et
	for _, t := range types {
		if t.Name == typeName {
			return fmt.Errorf("type zaten mevcut: %s", typeName)
		}
	}

	// Yeni type'ı ekle
	types = append(types, TemplateType{
		Name:        typeName,
		Description: description,
	})

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(types, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	if err := os.WriteFile(typesPath, data, 0644); err != nil {
		return fmt.Errorf("type'lar dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// AddTemplateFor, yeni bir template for ekler
func AddTemplateFor(name string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateForPath := filepath.Join(execDir, "template_util", "template_for.json")

	// Mevcut template for'ları al
	templateTypes, err := GetTemplateFors()
	if err != nil {
		return err
	}

	// Template for zaten var mı kontrol et
	for _, tt := range templateTypes {
		if tt.Name == name {
			return fmt.Errorf("template for zaten mevcut: %s", name)
		}
	}

	// Yeni template for'u ekle
	templateTypes = append(templateTypes, TemplateType{
		Name:        name,
		Description: fmt.Sprintf("Template for %s", name),
	})

	// JSON yapısını oluştur
	result := struct {
		Types []TemplateType `json:"types"`
	}{
		Types: templateTypes,
	}

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	// Template util klasörünü oluştur
	templateUtilDir := filepath.Join(execDir, "template_util")
	if err := os.MkdirAll(templateUtilDir, 0755); err != nil {
		return fmt.Errorf("template util klasörü oluşturulamadı: %v", err)
	}

	if err := os.WriteFile(templateForPath, data, 0644); err != nil {
		return fmt.Errorf("template for dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// GetTemplates, template_util/templates klasöründeki tüm template'leri döndürür
func GetTemplates() ([]string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	files, err := os.ReadDir(templateDir)
	if err != nil {
		return nil, fmt.Errorf("template klasörü okunamadı: %v", err)
	}

	var templates []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			templates = append(templates, file.Name())
		}
	}

	return templates, nil
}

// DeleteTemplates, seçilen template'leri siler
func DeleteTemplates(templates []string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	for _, template := range templates {
		filePath := filepath.Join(templateDir, template)
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("template silinemedi: %v", err)
		}
	}
	return nil
}

// GetPackages, mevcut paketleri listeler
func GetPackages() ([]Package, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	packagesPath := filepath.Join(execDir, "template_util", "packages.json")

	// Dosya yoksa boş liste döndür
	if _, err := os.Stat(packagesPath); os.IsNotExist(err) {
		return []Package{}, nil
	}

	data, err := os.ReadFile(packagesPath)
	if err != nil {
		return nil, fmt.Errorf("paketler dosyası okunamadı: %v", err)
	}

	var packages []Package
	if err := json.Unmarshal(data, &packages); err != nil {
		return nil, fmt.Errorf("paketler JSON parse hatası: %v", err)
	}

	return packages, nil
}

// DeletePackages, seçilen paketleri siler
func DeletePackages(packages []string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	packagesPath := filepath.Join(execDir, "template_util", "packages.json")

	existing, err := GetPackages()
	if err != nil {
		return err
	}

	// Silinecek paketleri listeden çıkar
	var newList []Package
	for _, pkg := range existing {
		shouldDelete := false
		for _, toDelete := range packages {
			if pkg.Name == toDelete {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			newList = append(newList, pkg)
		}
	}

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(newList, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	if err := os.WriteFile(packagesPath, data, 0644); err != nil {
		return fmt.Errorf("paketler dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// GetTemplateFors, mevcut template for'ları listeler
func GetTemplateFors() ([]TemplateType, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateForPath := filepath.Join(execDir, "template_util", "template_for.json")

	// Dosya yoksa boş liste döndür
	if _, err := os.Stat(templateForPath); os.IsNotExist(err) {
		return []TemplateType{}, nil
	}

	data, err := os.ReadFile(templateForPath)
	if err != nil {
		return nil, fmt.Errorf("template for dosyası okunamadı: %v", err)
	}

	var result struct {
		Types []TemplateType `json:"types"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("template for JSON parse hatası: %v", err)
	}

	return result.Types, nil
}

// DeleteTemplateFor, belirtilen template for'u siler
func DeleteTemplateFor(name string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateForPath := filepath.Join(execDir, "template_util", "template_for.json")

	// Mevcut template for'ları al
	templateTypes, err := GetTemplateFors()
	if err != nil {
		return err
	}

	// Template for'u bul ve sil
	var newTypes []TemplateType
	for _, tt := range templateTypes {
		if tt.Name != name {
			newTypes = append(newTypes, tt)
		}
	}

	// JSON yapısını oluştur
	result := struct {
		Types []TemplateType `json:"types"`
	}{
		Types: newTypes,
	}

	// JSON'a dönüştür ve kaydet
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	if err := os.WriteFile(templateForPath, data, 0644); err != nil {
		return fmt.Errorf("template for dosyası kaydedilemedi: %v", err)
	}

	return nil
}
