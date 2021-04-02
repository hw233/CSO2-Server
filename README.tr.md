## Counter-Strike Online 2 Sunucusu 

[![Build status](https://ci.appveyor.com/api/projects/status/a4pj1il9li5s08k5?svg=true)](https://ci.appveyor.com/project/KouKouChan/cso2-server)
[![](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![](https://img.shields.io/badge/version-v0.6.0-blue)](https://github.com/KouKouChan/CSO2-Server/releases)

#Birileri Projeyi Kötüye Kullandığından Github'daki açık kaynak kodu şimdilik yenilenmeyecektir.
#Taobao'daki 80 Game House, Qianyun Technology, Qianyun Games, Yinyue Entertainment Network ve diğer korsanlar gibi vicdansız tüccarlara güvenmeyin !
#Geliştirici diyorki bazı insanlar kötüye kullanım gerçekleştirdiği için yukarda yazılı Taobao kullanıcılarına güvenmeyin ve aldırış etmeyiniz
#Bu Sunucu KouKouChan'a Aittir farklı insanlar tarafından kötüye kullanılmaya devam edilirse sunucu desteği kesilecektir...
#Sadece KouKouChan Tarafından bilinen sunuculara destek verilmektedir. herşey için teşekkürler 'KouKouChan'

### 0x01 Oyun Açıklaması

Counter-Strike Online 2 Sunucusu

***L-leite tarafından [cso2-master-server](https://github.com/L-Leite/cso2-master-server) temel alınmıştır.***

*** Sunucu Dilini 'Location' Dosyasından Ayarlayınız Eğer Bilmiyorsanız Sayfa 3'e Göz Atınız.***

### 0x02 Oyunun Özellikleri

    1. Temel Oynanış [Tamamlandı] √
    2. Eksik Fonksiyonları Ekle Ve Geliştir..

### 0x03 Lokasyonunuzu Ayarlayın

```
1. 'CSO2-Server\locales\' klasöründe en-us.ini gibi bir .ini dosyası oluşturun [global.ini] (Kayıt Sayfası Dili)
2. Dil Dosyanızı Kaydetiğinizi ve Ayarladığınızı Varsayıyorum global.ini
3  'CSO2Server\CSO2-Server\configure' Dizininde server.conf 'Notepad++' İle Açın.
3. Server.conf dosyasını düzenleyin ve LocaleFile'ı dosya adınıza ayarlayın
4. LocaleFile= adlı metni aratıyoruz ve orada yazanı silip bunu yazalım LocaleFile=global.ini
5. Kayıt Sayfamız İçin E-Posta Sistemimizi Hazırlayalım
===========================
#E-Posta Sunucusunun Kullanıcı Adı
REGEmail=username@gmail.com
===========================
#E-Posta Sunucusunun Şifresi
REGPassWord=1547927439752
===========================
#'smtp.gmail.com' gibi bir e-posta sunucunuzu ayarlayın
REGSMTPaddr=smtp.gmail.com
===========================
6.Şimdi Sayfamız İçin Yazı Kodunu Ayarlayalım
#Dilinize göre sistem yazıları için, ZH-CN='gbk' , ZH-TW='big5' , GLOBAL='utf-8'
CodePage=utf-8
7.Şimdi Ana Oyun Yöneticimizin Kullanıcı Adı Ve Şifresini Ayarlayalım
===========================
#Yöneticinin Bağlantı Noktası (Bağlantı Portu)
GMport=1315
===========================
#Yöneticinin Kullanıcı Adı
GMusername=admin
===========================
#Yöneticinin Şifresi (Bunu Değiştirmenizi Öneririm Aksi Takdirde Yönetici Hesabınız Çalınabilir)
GMpassword=cso2server123
===========================
8.Şimdi Market Sistemimizi Etkinleştirelim
EnableShop=1
===========================
9.Şimdi Sunucumuzun Web Sunucu Kontrolünü Sağlayalım
Windows 10'da Denetim Masasını açın ve Sistem ve Güvenlik seçeneğini seçin.
Windows Güvenlik Duvarı yazısına tıklayın ve sol menüde çıkacak olan Windows Güvenlik Duvarı'nı etkinleştir veya devre dışı bırak seçeneğine tıklayın.
Özel ağ ve Ortak ağ seçenekleri altında Windows Güvenlik Duvarı'nı kapat (önerilmez) seçeneğini seçin. Ve Tamam butonuna tıklayın. 
Bu işlemle Windows 10 güvenlik duvarı kapatılacaktır. Daha sonra tekrardan aynı yolu izleyerek Firewall, güvenlik duvarını açabilirsiniz.
===========================
10. Firewall Kapatmadanda Yapabilirsiniz Portu Açarak (Fakat Bağlantı Kesilmesinin Önüne Geçebilmek İçin Kapatabilirsiniz)

```

### 0x04 Nasıl Bağlantılar Sağlanır?

    1. Oyununuzun Kore İstemcisi Olması Gereklidir (2017 Sürümüde Kullanabilirsiniz)
    2. L-leite'in github sayfasından bir başlatıcı indirin.
    3. En son oyun sunucusu dosyasını ( https://github.com/KouKouChan/CSO2-Server/releases ) adresinden indirin
    4. Oyun sunucusunu başlatın ve oyununuzu başlatmak için Başlatma dosyasını kullanın.
    5. İyi Eğlenceler

**Bildiri**!

Hala Web Sunucunuza Erişemiyorsanız Adımları Tekrar Kontrol Ediniz Hala Aynı Sorun Devam Ediyorsa Güvenlik Duvarı Ve Anti-Virüs Programını Kapatınız

### 0x05 Nasıl Yeni Sunucu Dosyası Oluşturulur

    1. Herhangi Bir Dizinde Klasör Açın Örnek 'C:\CSO2-Server'
    2. Sunucu Dosyalarını 'C:\CSO2-Server' Dizinine Atınız
    2. Sunucu Dosyalarınızı Dönüştürmek İçin 'go build' Komutunu Kullanabilirsiniz
    3. Çalıştırın (İşlem Tamamlandı)

### 0x06 Aşağıdaki Gereksinimler

    Go 1.15.6 (Zorunlu
    Bağlantı Noktaları:30001-TCP、30002-UDP (Zorunlu

***Bir LAN veya İnternet Sunucusu kurmak istiyorsanız, lütfen güvenlik duvarı bağlantı noktasını açın.***

### 0x07 Ekran Görüntüleri

![Image](https://i.hizliresim.com/fJgBch.png)

![Image](https://i.hizliresim.com/kVxUIG.png)

![Image](https://i.hizliresim.com/8DXgUk.png)

![Image](https://i.hizliresim.com/nSdDPk.png)

![Image](https://i.hizliresim.com/ysROOO.png)

![Image](https://i.hizliresim.com/JjJynK.png)

![Image](https://i.resmim.net/A4XDT.png)

