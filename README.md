# Netgrate

> **Netgrate** – A dynamic Go-based web server that embeds executable Go scripts directly into HTML pages.  
> Created by **kawa** (aka `2kpz`) – now discontinued.

---

## 🚫 Discontinued
This project is **no longer maintained**. It was an experimental proof-of-concept to explore dynamic server-side Go script execution within HTML. While functional, it has known security and scalability limitations. Use at your own risk.

🔗 **GitHub Repository**: [https://github.com/2kpz/Netgrate](https://github.com/2kpz/Netgrate)

---

## 🌐 What is Netgrate?

**Netgrate** is a lightweight web server written in Go that allows you to:

- Serve static HTML pages from a `www/` directory.
- Dynamically **execute Go scripts** from within HTML using custom `<<go(...)>>` tags.
- Create interactive **form buttons** (`<gobutton>`) that run Go scripts on click.
- Build simple dynamic web content without templates or external dependencies.

It’s like PHP, but for Go — with immediate script execution on page load or form submission.

---

## ⚙️ Features

### ✅ Dynamic Go Script Execution
Embed Go code directly in HTML:
```html
<p>Your external IP address is: <<go(scripts/hello.go)>></p>
```
The server runs `hello.go` and injects its output into the page.

### ✅ Interactive Buttons with `<gobutton>`
Create buttons that trigger Go scripts via POST:
```html
<gobutton type="scripts/ip.go" class="btn">Get IP</gobutton>
```
Renders as a form that executes the script and displays the result.

### ✅ Simple Static File Serving
Serves `.html` files from the `www/` folder. Automatically appends `.html` to clean URLs.

### ✅ Custom Error Pages
Built-in 404, 500, and 429 error pages with fun GIFs.

### ✅ No Dependencies
Pure Go – only uses the standard library.

---

## 📁 Project Structure

```
Netgrate/
├── main.go               # Web server logic
├── www/                  # Static HTML files
│   └── index.html
├── scripts/              # Go scripts for dynamic content
│   └── hello.go
└── README.md
```

---

## 🧪 Example: `scripts/hello.go`

```go
package main

import "fmt"

func main() {
    fmt.Print("Hello from Go!")
}
```

Used in HTML:
```html
<p>Result: <<go(scripts/hello.go)>></p>
```

Output:
```
Result: Hello from Go!
```

---

## 🖱️ Using `<gobutton>`

```html
<gobutton type="scripts/ip.go" class="btn primary-btn">Show My IP</gobutton>
```

This becomes a POST form that runs `ip.go` and shows the output on the same page.

---

## 🛠️ Setup & Usage

### 1. Clone the repo
```bash
git clone https://github.com/2kpz/Netgrate.git
cd Netgrate
```

### 2. Add your scripts
Put `.go` files in the `scripts/` directory.

### 3. Add HTML pages
Place `.html` files in the `www/` folder.

### 4. Run the server
```bash
go run main.go
```

> ⚠️ Server runs on `http://localhost` (port 80, may require admin/sudo on Linux).

---

## ⚠️ Security Warning

- **Arbitrary Code Execution**: Any `.go` file in the project can be executed via URL.
- **No Input Sanitization**: Malicious scripts can be run if users can upload code.
- **Not for Production**: This is a proof-of-concept. **Do not use in production or public environments.**

---

## 🧱 Why Was It Discontinued?

- **Security Risks**: Running arbitrary Go code server-side is inherently dangerous.
- **Performance**: Spawning `go run` processes is slow and resource-heavy.
- **Better Alternatives**: Templates (Go `html/template`), WASM, or proper backends are safer and more scalable.

---

## 🙌 Credits

- **Creator**: [kawa](https://github.com/2kpz) (`@2kpz`)
- **Inspiration**: PHP-style dynamic scripting, simplicity-first design.
- **GIFs**: Hosted via [iili.io](https://iili.io/3RaXzk7.gif)

---

## 📄 License

MIT License – See [LICENSE](LICENSE) for details.

---

> 💡 **Made by kawa** | `@2kpz` | Discontinued – 2024  
> For educational purposes only.
