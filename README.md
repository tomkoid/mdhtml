<img width="300" src="assets/logo.jpg">

Really simple CLI Markdown to HTML converter with styling support

## 🌟 Usage

To convert a Markdown file to HTML, you can run the following command:

```bash
mdhtml convert <input file>
```

This will create a file called `output.html` in the current directory. You can also specify the output file name:

```bash
mdhtml convert <input file> -o <output file>
```

You can also specify a CSS file to style the HTML file:

```bash
mdhtml convert <input file> -o <output file> -s <css file>
```

To see all available options, you can run:

```bash
mdhtml --help
```

## 💻 Installation

[![Packaging status](https://repology.org/badge/vertical-allrepos/mdhtml.svg)](https://repology.org/project/mdhtml/versions)

To install mdhtml, you can either download the binary from the [releases page](https://codeberg.org/Tomkoid/mdhtml/releases) and install it to the system or build it from source.

### Installing the binary

To install the binary, you can download it from the [releases page](https://codeberg.org/Tomkoid/mdhtml/releases) and install it to your system (Linux):

```bash
sudo mv mdhtml /usr/local/bin
```

### Building from source

To build mdhtml from source, you need to have [Go](https://golang.org/) installed on your system. Then, you can run the following command to build the binary:

```bash
go build -o mdhtml main.go
```

Then you can install the binary to your system (Linux):

```bash
sudo mv mdhtml /usr/local/bin
```

### Installing from AUR

If you are using Arch Linux, you can install mdhtml from the [AUR](https://aur.archlinux.org/packages/mdhtml/).
