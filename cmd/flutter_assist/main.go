package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/burak/flutter_assist/internal/project"
	"github.com/burak/flutter_assist/internal/template"
)

// getExecutableDir, çalıştırılabilir dosyanın bulunduğu dizini döndürür
func getExecutableDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("çalıştırılabilir dosya yolu alınamadı: %v", err)
	}
	return filepath.Dir(execPath), nil
}

func main() {
	// Komut satırı argümanlarını tanımla
	templateFlag := flag.String("t", "", "Template oluşturmak için dosya veya klasör yolu")
	templateForFlag := flag.Bool("tf", false, "Template for işlemleri için")
	packageFlag := flag.Bool("p", false, "Paket işlemleri için")
	templateDeleteFlag := flag.Bool("tdelete", false, "Template silme işlemi için")
	templateForDeleteFlag := flag.Bool("tfdelete", false, "Template for silme işlemi için")
	packageDeleteFlag := flag.Bool("pdelete", false, "Paket silme işlemi için")
	flag.Parse()

	// Emoji tanımlamaları
	const (
		infoEmoji    = "ℹ️"
		successEmoji = "✅"
		errorEmoji   = "❌"
	)

	// Template silme işlemi
	if *templateDeleteFlag {
		execDir, err := getExecutableDir()
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Çalışma dizinini executable dizinine değiştir
		if err := os.Chdir(execDir); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		templates, err := project.GetTemplates()
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		if len(templates) == 0 {
			fmt.Printf("%s Silinecek template bulunamadı\n", infoEmoji)
			return
		}

		// Template'leri listele
		fmt.Println("📋 Template'ler:")
		for i, template := range templates {
			fmt.Printf("%d. %s\n", i+1, template)
		}

		// Kullanıcıdan seçim iste
		fmt.Printf("\nℹ️ Silinecek template'lerin numaralarını boşlukla ayırarak girin (örn: 1 3): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Seçilen numaraları parse et
		selected := make(map[string]bool)
		numbers := strings.Fields(input)
		for _, numStr := range numbers {
			num, err := strconv.Atoi(numStr)
			if err != nil || num < 1 || num > len(templates) {
				continue
			}
			selected[templates[num-1]] = true
		}

		// Seçilen template'leri sil
		var templatesToDelete []string
		for template, isSelected := range selected {
			if isSelected {
				templatesToDelete = append(templatesToDelete, template)
			}
		}

		if err := project.DeleteTemplates(templatesToDelete); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		fmt.Printf("%s Template'ler başarıyla silindi!\n", successEmoji)
		return
	}

	// Template for silme işlemi
	if *templateForDeleteFlag {
		fmt.Println("ℹ️ Template for silme modu başlatılıyor...")
		templateFors, err := project.GetTemplateFors()
		if err != nil {
			fmt.Printf("❌ Hata: %v\n", err)
			return
		}

		if len(templateFors) == 0 {
			fmt.Println("ℹ️ Silinecek template for bulunamadı")
			return
		}

		selected, err := selectFromTemplateTypes(templateFors, "Silinecek template for'u seçin:")
		if err != nil {
			fmt.Printf("❌ Hata: %v\n", err)
			return
		}

		if err := project.DeleteTemplateFor(selected); err != nil {
			fmt.Printf("❌ Hata: %v\n", err)
			return
		}

		fmt.Printf("✅ Template for başarıyla silindi: %s\n", selected)
		return
	}

	// Paket silme işlemi
	if *packageDeleteFlag {
		fmt.Println("ℹ️ Paket silme modu başlatılıyor...")
		packages, err := project.GetPackages()
		if err != nil {
			fmt.Printf("❌ Hata: %v\n", err)
			return
		}

		if len(packages) == 0 {
			fmt.Println("ℹ️ Silinecek paket bulunamadı")
			return
		}

		selected, err := selectFromPackages(packages, "Silinecek paketi seçin:")
		if err != nil {
			fmt.Printf("❌ Hata: %v\n", err)
			return
		}

		if err := project.DeletePackage(selected); err != nil {
			fmt.Printf("❌ Hata: %v\n", err)
			return
		}

		fmt.Printf("✅ Paket başarıyla silindi: %s\n", selected)
		return
	}

	// Template for ekleme modu
	if *templateForFlag {
		fmt.Printf("%s Template for ekleme modu başlatılıyor...\n", infoEmoji)

		// Komut satırı argümanlarını kontrol et
		if len(flag.Args()) == 0 {
			fmt.Println("❌ Template for ismi belirtilmedi")
			return
		}

		// Template for ismini al
		templateForName := flag.Args()[0]

		execDir, err := getExecutableDir()
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Çalışma dizinini executable dizinine değiştir
		if err := os.Chdir(execDir); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Template for'u ekle
		if err := project.AddTemplateFor(templateForName); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		fmt.Printf("%s Template for başarıyla eklendi!\n", successEmoji)
		return
	}

	// Type seçimi sadece gerekli komutlar için
	var selectedTypes []string
	if *templateFlag != "" || *packageFlag || len(flag.Args()) > 0 {
		selectedTypes = selectTypes()
		if len(selectedTypes) == 0 {
			fmt.Println("❌ En az bir type seçilmelidir")
			return
		}
	}

	// Template oluşturma modu
	if *templateFlag != "" {
		fmt.Printf("%s Template oluşturma modu başlatılıyor...\n", infoEmoji)

		// Template oluştur
		err := template.CreateTemplate(*templateFlag, selectedTypes)
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		fmt.Printf("%s Template başarıyla oluşturuldu!\n", successEmoji)
		return
	}

	// Paket ekleme modu
	if *packageFlag {
		fmt.Printf("%s Paket ekleme modu başlatılıyor...\n", infoEmoji)

		// Komut satırı argümanlarını kontrol et
		if len(flag.Args()) == 0 {
			fmt.Println("❌ Paket ismi belirtilmedi")
			return
		}

		// Paket ismini al
		packageName := flag.Args()[0]

		execDir, err := getExecutableDir()
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Çalışma dizinini executable dizinine değiştir
		if err := os.Chdir(execDir); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Paketi ekle
		if err := project.AddPackage(packageName, selectedTypes); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Paketin eklendiğini kontrol et
		packages, err := project.GetPackages()
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Paket eklendi mi kontrol et
		found := false
		for _, pkg := range packages {
			if pkg.Name == packageName {
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("%s Paket eklenemedi!\n", errorEmoji)
			os.Exit(1)
		}

		fmt.Printf("%s Paket başarıyla eklendi!\n", successEmoji)
		return
	}

	// Proje oluşturma modu
	if len(flag.Args()) > 0 {
		projectName := flag.Args()[0]
		fmt.Printf("%s Proje oluşturma modu başlatılıyor...\n", infoEmoji)

		execDir, err := getExecutableDir()
		if err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Çalışma dizinini executable dizinine değiştir
		if err := os.Chdir(execDir); err != nil {
			fmt.Printf("%s Hata: %v\n", errorEmoji, err)
			os.Exit(1)
		}

		// Projeyi oluştur
		if err := project.CreateProject(projectName, selectedTypes); err != nil {
			fmt.Printf("❌ Proje oluşturulamadı: %v\n", err)
			return
		}

		fmt.Printf("%s Proje başarıyla oluşturuldu!\n", successEmoji)
		return
	}

	// Yardım mesajı
	fmt.Printf("%s Kullanım:\n", infoEmoji)
	fmt.Println("  flutter_assist <proje_ismi>          - Proje oluştur")
	fmt.Println("  flutter_assist -t <template_ismi>    - Template oluştur")
	fmt.Println("  flutter_assist -p                    - Paket ekle")
	fmt.Println("  flutter_assist -tf                   - Template for ekle")
	fmt.Println("  flutter_assist -tdelete              - Template'leri sil")
	fmt.Println("  flutter_assist -pdelete              - Paketleri sil")
	fmt.Println("  flutter_assist -tfdelete             - Template for'ları sil")
}

func selectTypes() []string {
	types, err := project.GetTemplateTypes()
	if err != nil {
		fmt.Printf("❌ Type'lar alınamadı: %v\n", err)
		return nil
	}

	// Type'ları listele
	fmt.Println("📋 Type'lar:")
	for i, t := range types {
		fmt.Printf("%d. %s\n", i+1, t.Name)
	}

	// Kullanıcıdan seçim iste
	fmt.Printf("\nℹ️ Seçilecek type'ların numaralarını boşlukla ayırarak girin (örn: 1 3): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Seçilen numaraları parse et
	selected := make(map[string]bool)
	numbers := strings.Fields(input)
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 1 || num > len(types) {
			continue
		}
		selected[types[num-1].Name] = true
	}

	// Seçilen type'ları döndür
	var result []string
	for t, isSelected := range selected {
		if isSelected {
			result = append(result, t)
		}
	}
	return result
}

func selectFromList(items []string, prompt string) (string, error) {
	// Öğeleri listele
	fmt.Println(prompt)
	for i, item := range items {
		fmt.Printf("%d. %s\n", i+1, item)
	}

	// Kullanıcıdan seçim iste
	fmt.Printf("\nℹ️ Seçilen öğelerin numaralarını boşlukla ayırarak girin (örn: 1 3): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Seçilen numaraları parse et
	selected := make(map[string]bool)
	numbers := strings.Fields(input)
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 1 || num > len(items) {
			continue
		}
		selected[items[num-1]] = true
	}

	// Seçilen öğeleri döndür
	var result string
	for item, isSelected := range selected {
		if isSelected {
			result = item
			break
		}
	}
	return result, nil
}

func selectFromTemplateTypes(items []project.TemplateType, prompt string) (string, error) {
	var itemNames []string
	for _, item := range items {
		itemNames = append(itemNames, item.Name)
	}
	return selectFromList(itemNames, prompt)
}

func selectFromPackages(items []project.Package, prompt string) (string, error) {
	var itemNames []string
	for _, item := range items {
		itemNames = append(itemNames, item.Name)
	}
	return selectFromList(itemNames, prompt)
}
