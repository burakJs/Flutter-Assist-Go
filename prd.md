# Flutter Proje Başlatıcı CLI - Proje Gereksinim Dokümanı (PRD)

## 1. Amaç
Kullanıcıdan proje ismi ve type seçimi alarak, otomatik şekilde Flutter projesi oluşturan, gerekli dosya ve paketleri ekleyen, esnek ve genişletilebilir bir komut satırı aracı (CLI) geliştirmek.

## 2. Hedef Kitle
Flutter ile hızlıca yeni projeler başlatmak isteyen geliştiriciler.

## 3. Temel Özellikler
- Kullanıcıdan proje ismi alma
- Type seçimi için çoklu seçim yapabilme
- `flutter create <proje-ismi>` ile yeni Flutter projesi oluşturma
- Seçilen type'lara göre gerekli dosya ve klasörleri oluşturma
- Gerekli Flutter paketlerini otomatik ekleme (`flutter pub add <paket>`)
- Dosya içeriklerinde sabit bir anahtar kelime olan `{FLUTTER_ASSIST}` ifadesini proje ismiyle değiştirme
- Tüm yapılandırmaların bir klasör altında yer alan JSON dosyaları ile yönetilmesi
- Her adımda emojili bilgilendirme mesajları gösterme
- Genişletilebilir ve sürdürülebilir kod yapısı (Go dili ile)
- Template ve paket yönetimi için kapsamlı komut seti
- Klasör ve dosya bazlı template oluşturma desteği

## 4. CLI Kullanım Şekilleri

### Build Alma
```bash
go build -o flutter_assist cmd/flutter_assist/main.go
```

### Proje Oluşturma
```bash
flutter_assist <proje_ismi>
```
- Type seçimi için çoklu seçim ekranı gösterir
- Seçilen type'lara göre projeyi oluşturur

### Template Oluşturma
```bash
flutter_assist -t <dosya_veya_klasor_yolu>
```
- Type seçimi için çoklu seçim ekranı gösterir
- Verilen dosya veya klasörü seçilen type'lar için template olarak kaydeder
- Klasör verildiğinde içindeki tüm dosyaları template olarak kaydeder
- Proje ismi sorma ve `{FLUTTER_ASSIST}` ile değiştirme

### Paket Ekleme
```bash
flutter_assist -p
```
- Type seçimi için çoklu seçim ekranı gösterir
- Paket ismini sorar
- Seçilen type'lar için paketi ekler

### Template For Ekleme
```bash
flutter_assist -tf
```
- Yeni bir template for ismi ekler

### Template Silme
```bash
flutter_assist -tdelete
```
- Mevcut template'leri listeler
- Çoklu seçim ile silinecek template'leri seçtirir

### Paket Silme
```bash
flutter_assist -pdelete
```
- Mevcut paketleri listeler
- Çoklu seçim ile silinecek paketleri seçtirir

### Template For Silme
```bash
flutter_assist -tfdelete
```
- Mevcut template for'ları listeler
- Çoklu seçim ile silinecek template for'ları seçtirir

## 5. Proje Oluşturma İş Akışı
1. Kullanıcıdan proje ismi alınır
2. Type seçimi için çoklu seçim ekranı gösterilir
3. Proje kökünde bir klasör (ör: `template_util/templates/`) altında yer alan tüm JSON dosyaları okunur:
   - Her bir JSON dosyası, oluşturulacak bir dosyanın path'ini ve içeriğini içerir
   - Dosya içeriklerinde `{FLUTTER_ASSIST}` anahtar kelimesi bulunur ve bu ifade proje ismiyle değiştirilir
   - Tüm bu JSON dosyaları bir araya getirilerek mimari oluşturulur
4. Paketler için ayrı bir JSON dosyası (ör: `template_util/packages.json`) okunur:
   - Eklenecek tüm paketler ve hangi type'larda ekleneceği bu dosyada tutulur
5. `flutter create <proje-ismi>` komutu çalıştırılır
6. Gerekli paketler eklenir
7. Dosyalar ilgili yerlere oluşturulur ve içerikleri anahtar kelimeye göre düzenlenir
8. Her adımda kullanıcıya emojili bilgilendirme mesajı gösterilir
9. İşlem sonunda başarı mesajı gösterilir

## 6. Template Oluşturma İş Akışı
1. Kullanıcı `flutter_assist -t dosya_veya_klasor_pathi` komutunu çalıştırır
2. Type seçimi için çoklu seçim ekranı gösterilir
3. Kullanıcıdan proje ismi istenir
4. Seçilen dosya(lar)ın içeriğinde proje ismi geçen tüm yerler tespit edilir ve bunlar `{FLUTTER_ASSIST}` anahtar kelimesiyle değiştirilir
5. Her dosya için uygun JSON formatı oluşturulur:
   - Dosya path'i, başındaki `/` işareti kaldırılarak JSON'a yazılır
   - `{FLUTTER_ASSIST}` ile değiştirilmiş dosya içeriği
6. Oluşan JSON dosyası, dosya ismiyle `template_util/templates` klasörüne kaydedilir
7. Kullanıcıya emojili bilgilendirme ve başarı mesajı gösterilir

## 7. JSON Yapısı
- **template_util/templates/**: İçerisinde her biri oluşturulacak bir dosyayı tanımlayan JSON dosyaları bulunur
- **template_util/packages.json:** Eklenecek tüm paketlerin listesi ve hangi type'larda ekleneceği bu dosyada tutulur
- **template_util/template_for.json:** Kullanılabilir type'ların listesi bu dosyada tutulur

## 8. Genişletilebilirlik
- İleride yeni type'lar kolayca eklenebilecek
- Anahtar kelimeye ek olarak, PascalCase, camelCase gibi dönüşümler ileride eklenebilir
- Yeni dosya şablonları eklemek için sadece ilgili klasöre yeni bir JSON dosyası eklemek yeterli olacaktır
- Template oluşturma özelliği ile kullanıcılar kendi klasör veya dosyalarını kolayca template JSON'a dönüştürüp sisteme ekleyebilir

## 9. Kullanıcı Deneyimi
- Her adımda ne yapıldığına dair açıklayıcı ve emojili mesajlar
- Hata durumunda kullanıcıya açıklayıcı hata mesajı
- İşlem sonunda özet ve başarı mesajı
- Çoklu seçim için kullanıcı dostu arayüz

## 10. Teknik Gereksinimler
- Go dili ile yazılacak
- Flutter ve Go sistemde kurulu olmalı
- JSON dosyaları ile yapılandırma
- Tüm mimari şablonlar ve paketler ilgili klasörlerde JSON dosyaları olarak tutulacak
- Template ve paket yönetimi için kapsamlı komut seti 