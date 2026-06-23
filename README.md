# 📄 CLI File Converter

A fast, secure, and lossless command-line tool built with Go (Golang) for converting various file formats (Images, Documents, Audio, Video).

## 🚀 Features

* **Fast & Efficient:** Built with Go for high performance and minimal resource usage.
* **CLI-First Design:** Easy to use directly from your terminal, perfect for scripting and automation.
* **Multi-Format Support:** Convert between popular formats effortlessly.
* **Lossless Conversion:** Maintains the highest possible quality during the conversion process.
* **Standalone Binary:** Compiles into a single executable file with no external runtime dependencies.

## 📂 Supported Formats

| Category | Source Formats | Target Formats |
| :--- | :--- | :--- |
| **Images** | PNG, JPG, JPEG, WEBP, BMP | PDF, PNG, JPG, WEBP |
| **Documents**| DOCX, PDF, TXT, HTML | PDF, DOCX, TXT |
| **Audio** | MP3, WAV, OGG, M4A | MP3, WAV |
| **Video** | MP4, AVI, MOV, MKV | MP4, MP3 (Audio Extraction) |

## 🛠️ Installation

### Prerequisites
* [Go 1.18+](https://go.dev/doc/install) installed on your system.
* Git
* Libreoffice headless
* FFmpeg

### Build from Source

1. **Clone the repository:**
  ```
   bash
   git clone https://github.com/bedirmirac/convertev.git
   cd convertev
   ```

2. Download dependencies:
  ```
  go mod tidy
  ```

3. Build the binary:
  ``` 
  go build -o convertev main.go 
  ```

4. (Optional) Move to a directory in your PATH for global use:
  ```
  sudo mv convertev /usr/local/bin/
  ```

### 💻 Usage
Run the tool directly from your terminal.

Basic syntax:
```
Usage: convertev -mode [img|media] -i [input] -o [output] || convertev -mode [doc] -i [input]
```

Examples:

Convert an image to PNG:
```
./convertev -mode img -i img.jpg -o img.png
```
View help menu and available commands:
```
./convertev
```

🛠️ Tech Stack
Language: Go (Golang)

Libraries:
+ github.com/disintegration/imaging
+ github.com/pdfcpu/pdfcpu
+ golang.org/x/image 

# 🤝 Contributing
- Contributions are always welcome! To contribute:

- Fork the repository.

- Create a new branch (git checkout -b feature/AmazingFeature).

- Commit your changes (git commit -m 'Add some AmazingFeature').

- Push to the branch (git push origin feature/AmazingFeature).

- Open a Pull Request.

# 📄 License
This project is licensed under the MIT License.

# ⚖️ Acknowledgments & Licenses
convertev is a wrapper tool that utilizes the following open-source software to perform conversions:

FFmpeg: Used for media conversions. FFmpeg is a trademark of Fabrice Bellard, originator of the FFmpeg project. Licensed under the LGPL/GPL.

LibreOffice: Used for document conversions. Licensed under the Mozilla Public License v2.0.

Note: convertev does not distribute these binaries. Users must install them separately.

# ✉️ Contact
Developer: Miraç Bedir

GitHub: @bedirmirac
