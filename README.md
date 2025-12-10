# Majotabi - Golang WhatsApp Bot

<p align="center">
  <img src="https://cdn.nefusoft.cloud/ICKeO.jpg" width="200" alt="Elaina Profile">
</p>

<p align="center">
  Sebuah bot WhatsApp canggih, cepat, dan modular yang dibangun menggunakan <strong>Go (Golang)</strong> dan library <strong>whatsmeow</strong>. Bot ini dirancang untuk stabilitas, kecepatan, dan kemudahan dalam penambahan fitur baru melalui sistem plugin.
</p>

## ✨ Fitur Utama

-   🤖 **AI Tsundere (Elaina)**: Ditenagai oleh Google Gemini, Elaina dapat diajak ngobrol, menjawab pertanyaan, dan bahkan "melihat" gambar. Persona tsundere-nya yang unik memberikan pengalaman interaksi yang menyenangkan.
-   🖼️ **Gambar Random**: Menghasilkan gambar anime acak dari berbagai kategori (`!neko`, `!waifu`, `!maid`) dengan logika *retry* untuk memastikan media selalu terkirim.
-   👑 **Mode Owner**: Bot dapat diatur ke mode `!self` (hanya merespon owner) atau `!public` (merespon semua orang).
-   🚀 **Super Cepat & Stabil**: Dibangun dengan Go, bot ini menawarkan performa tinggi, *startup* cepat, dan penggunaan memori yang efisien.
-   🧩 **Sistem Plugin Modular**: Menambah fitur baru sangat mudah. Cukup buat file `.go` baru di dalam folder `feature/` tanpa perlu menyentuh kode inti.
-   💅 **Tampilan & UX Modern**:
    -   Menu interaktif dengan *thumbnail*.
    -   Log terminal berwarna yang mudah dibaca.
    -   Reaksi emoji instan pada setiap perintah.
    -   ASCII art keren saat proses login.
    -   Anti pesan lama saat bot baru dinyalakan.

---

## ⚙️ Prasyarat

Sebelum memulai, pastikan sistem Anda telah terpasang:
1.  **Go (Golang)**: Versi 1.18 atau lebih baru.
2.  **Git**: Untuk melakukan kloning repository.

---

## 🚀 Instalasi & Menjalankan

Ikuti langkah-langkah berikut untuk menjalankan bot di server (VPS) atau komputer lokal Anda.

**1. Clone Repository**
```bash
git clone https://github.com/your-username/Majotabi-GolangWAbot.git
cd Majotabi-GolangWAbot
```
*(Jangan lupa ganti `your-username` dengan username GitHub Anda)*

**2. Konfigurasi Bot**

Buka file `config.json` dan isi semua nilainya sesuai data Anda.

```json
{
    "owner_number": "628xxxxxxxxxx",
    "owner_lid": "xxxxxxxxxxxxxxxxx",
    "bot_name": "Majotabi - Golang",
    "prefix": "!",
    "thumbnail": "https://cdn.nefusoft.cloud/ICKeO.jpg",
    "gemini_api_keys": [
        "YOUR_GEMINI_API_KEY_1",
        "YOUR_GEMINI_API_KEY_2"
    ]
}
```
-   `owner_number`: Nomor WhatsApp Anda dengan format `62...`.
-   `owner_lid`: LID Anda untuk deteksi owner di grup. (Bisa didapat dari log bot saat Anda mengirim pesan).
-   `gemini_api_keys`: Masukkan API Key Anda dari [Google AI Studio](https://aistudio.google.com/app/apikey).

**3. Install Dependencies**

Terminal akan otomatis mengunduh semua library yang dibutuhkan.
```bash
go mod tidy
```

**4. Jalankan Bot!**
```bash
go run .
```
-   Pada saat pertama kali menjalankan, Anda akan diminta memilih metode login. Pilih **2 (Pairing Code)**.
-   Masukkan nomor WhatsApp yang akan dijadikan bot (format `62...`).
-   Sebuah kode akan muncul di terminal. Buka WhatsApp di HP Anda → Perangkat Tertaut → Tautkan Perangkat → Masukkan kode tersebut.
-   Bot Anda sekarang sudah online!

---

## 📌 Contoh Penggunaan

-   `!menu` - Menampilkan daftar kategori fitur.
-   `!menu random` - Menampilkan semua fitur di kategori `random`.
-   `!allmenu` - Menampilkan semua perintah yang tersedia.
-   `!neko` - Mengirim gambar neko (kucing anime) acak.
-   `!elaina on` - Mengaktifkan mode AI.
-   `elaina, apa kabar?` - Berinteraksi dengan AI setelah diaktifkan.

---

## 📁 Struktur Proyek

```
.
├── feature/          # Folder utama untuk semua plugin/fitur
│   ├── ai/
│   ├── main/
│   ├── owner/
│   └── random/
├── lib/              # Library internal/pembantu
│   ├── config/
│   ├── database/
│   ├── helper/
│   └── registry/
├── config.json       # File konfigurasi utama
├── go.mod            # Manajemen dependency Go
├── main.go           # File inti untuk menjalankan bot
└── README.md         # Dokumentasi ini
```

