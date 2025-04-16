# 🚀 Flutter Assist CLI

Flutter projelerinizi saniyeler içinde oluşturmanızı sağlayan, modern ve güçlü bir komut satırı aracı. Projenizi baştan sona yapılandırır, gerekli paketleri ekler ve özelleştirilebilir mimari şablonları uygular.

## ✨ Öne Çıkan Özellikler

- ⚡ **Hızlı Proje Oluşturma**: Tek komutla Flutter projenizi oluşturun
- 🏗️ **Akıllı Mimari**: Type bazlı proje yapılandırması ile istediğiniz mimariyi seçin
- 📦 **Otomatik Paket Yönetimi**: Gerekli paketleri otomatik olarak ekler ve yapılandırır
- 🎨 **Özelleştirilebilir Şablonlar**: Kendi template'lerinizi oluşturun ve yönetin
- 🔄 **Kolay Genişletilebilirlik**: Yeni type'lar ve şablonlar ekleyerek sistemi genişletin
- 💫 **Modern CLI Deneyimi**: Emojili bildirimler ve kullanıcı dostu arayüz

## 🛠️ Kurulum

### Go ile Build Alma
```bash
go build -o flutter_assist cmd/flutter_assist/main.go
```

### Go Install ile Kurulum
```bash
go install github.com/burakJs/flutter_assist@latest
```

## 🚀 Kullanım

### Proje Oluşturma
```bash
flutter_assist <proje_ismi>
```
- Type seçimi için interaktif menü
- Seçilen type'lara göre proje yapılandırması
- Otomatik paket ekleme
- Özelleştirilmiş dosya yapısı

### Template Yönetimi
```bash
# Template oluşturma
flutter_assist -t <dosya_veya_klasor_yolu>

# Template silme
flutter_assist -tdelete

# Template for ekleme
flutter_assist -tf

# Template for silme
flutter_assist -tfdelete
```

### Paket Yönetimi
```bash
# Paket ekleme
flutter_assist -p

# Paket silme
flutter_assist -pdelete
```

## 🏗️ Proje Yapısı

CLI, aşağıdaki yapılandırma dosyalarını kullanır:

- 📁 `template_util/templates/`: Özelleştirilebilir dosya şablonları
- 📦 `template_util/packages.json`: Paket yapılandırmaları
- 🏷️ `template_util/template_for.json`: Kullanılabilir type'lar

## 🔄 İş Akışı

1. **Proje Oluşturma**:
   - Type seçimi
   - Proje oluşturma
   - Paket ekleme
   - Dosya yapısı oluşturma

2. **Template Oluşturma**:
   - Dosya/klasör seçimi
   - Type atama
   - Otomatik proje ismi değiştirme
   - JSON formatında kaydetme

3. **Paket Yönetimi**:
   - Type bazlı paket ekleme
   - Paket silme
   - Yapılandırma güncelleme

## 🎯 Örnek Kullanım

```bash
# Yeni bir Flutter projesi oluştur
flutter_assist my_awesome_app

# Template oluştur
flutter_assist -t lib/core/context/app_provider.dart

# Paket ekle
flutter_assist -p
```

## 🤝 Katkıda Bulunma

1. Bu repo'yu fork edin
2. Yeni bir branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add some amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Bir Pull Request oluşturun

## 📝 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.

## 👨‍💻 Geliştirici

Bu proje [Burak Sekili](https://github.com/buraksekili) tarafından geliştirilmiştir. 
