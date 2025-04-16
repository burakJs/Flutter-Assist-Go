package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Template yapısı
type Template struct {
	Path    string   `json:"path"`
	Content string   `json:"content"`
	Types   []string `json:"types"`
}

// CreateTemplate, yeni bir template oluşturur
func CreateTemplate(templatePath string, types []string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	// Proje ismini al
	fmt.Print("📝 Proje ismini girin: ")
	var projectName string
	fmt.Scanln(&projectName)

	// Template util klasörünü oluştur
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("template klasörü oluşturulamadı: %v", err)
	}

	// Dosya mı klasör mü kontrol et
	info, err := os.Stat(templatePath)
	if err != nil {
		return fmt.Errorf("template yolu okunamadı: %v", err)
	}

	if info.IsDir() {
		// Klasör ise içindeki tüm dosyaları işle
		return processDirectory(templatePath, types, templateDir, projectName)
	} else {
		// Dosya ise tek dosyayı işle
		return processFile(templatePath, types, templateDir, projectName)
	}
}

// GetTemplate, belirtilen template'i döndürür
func GetTemplate(templateName string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	// Template dosyasını oku
	templateFile := filepath.Join(templateDir, templateName+".json")
	data, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("template dosyası okunamadı: %v", err)
	}

	// Template yapısını parse et
	var template struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(data, &template); err != nil {
		return "", fmt.Errorf("template JSON parse hatası: %v", err)
	}

	return template.Content, nil
}

// DeleteTemplate, belirtilen template'i siler
func DeleteTemplate(templateName string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	// Template dosyasını sil
	templateFile := filepath.Join(templateDir, templateName+".json")
	if err := os.Remove(templateFile); err != nil {
		return fmt.Errorf("template dosyası silinemedi: %v", err)
	}

	return nil
}

// GetTemplateTypes, template'in type'larını döndürür
func GetTemplateTypes(templateName string) ([]string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	// Template dosyasını oku
	templateFile := filepath.Join(templateDir, templateName+".json")
	data, err := os.ReadFile(templateFile)
	if err != nil {
		return nil, fmt.Errorf("template dosyası okunamadı: %v", err)
	}

	// Template yapısını parse et
	var template struct {
		Types []string `json:"types"`
	}
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("template JSON parse hatası: %v", err)
	}

	return template.Types, nil
}

// UpdateTemplate, template'i günceller
func UpdateTemplate(templateName string, content string, types []string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}

	execDir := filepath.Dir(execPath)
	templateDir := filepath.Join(execDir, "template_util", "templates")

	// Template yapısını oluştur
	template := struct {
		Path    string   `json:"path"`
		Content string   `json:"content"`
		Types   []string `json:"types"`
	}{
		Path:    templateName,
		Content: content,
		Types:   types,
	}

	// JSON'a dönüştür
	jsonData, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	// Template dosyasını kaydet
	templateFile := filepath.Join(templateDir, templateName+".json")
	if err := os.WriteFile(templateFile, jsonData, 0644); err != nil {
		return fmt.Errorf("template dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// processFile, tek bir dosyayı template'e dönüştürür
func processFile(filePath string, types []string, templateDir string, projectName string) error {
	// Dosya içeriğini oku
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("dosya okunamadı: %v", err)
	}

	// İçerikteki proje ismini {FLUTTER_ASSIST} ile değiştir
	contentStr := strings.ReplaceAll(string(content), projectName, "{FLUTTER_ASSIST}")

	// Template oluştur
	template := struct {
		Path    string   `json:"path"`
		Content string   `json:"content"`
		Types   []string `json:"types"`
	}{
		Path:    strings.TrimPrefix(filePath, "/"),
		Content: contentStr,
		Types:   types,
	}

	// JSON'a dönüştür
	jsonData, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON dönüştürme hatası: %v", err)
	}

	// Template dosyasını kaydet
	fileName := filepath.Base(filePath) + ".json"
	outputPath := filepath.Join(templateDir, fileName)
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("template dosyası kaydedilemedi: %v", err)
	}

	return nil
}

// processDirectory, bir klasörü ve içindeki tüm dosyaları template'e dönüştürür
func processDirectory(dirPath string, types []string, templateDir string, projectName string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return processFile(path, types, templateDir, projectName)
		}
		return nil
	})
}
