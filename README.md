<img width="300" src="assets/logo.jpg">

Really simple CLI Markdown to HTML converter with styling support

## Showcase

[![asciicast](https://asciinema.org/a/624645.svg)](https://asciinema.org/a/624645)

## 🌟 Usage

To convert a Markdown file to HTML, you can run the following command:

```bash
mdhtml convert <input file>
```

This will create a file called `output.html` in the current directory. You can also specify the output file name:

```bash
mdhtml convert <input file> --output <output file>
```

You can also specify a CSS file to style the HTML file:

```bash
mdhtml convert <input file> --output <output file> --stylesheet <css file>
```

MDHTML also supports watching a file for changes and automatically updating the output file:

```bash
mdhtml convert <input file> --watch
```

If you want, you can also enable an HTTP server with live reloading:

```bash
mdhtml convert <input file> --watch --httpserver
```



To see all available options, you can run:

```bash
mdhtml --help
```

## 💻 Installation

mdhtml is currently packaged only for Arch Linux and nixpkgs. You can find the packages here:

[![Packaging status](https://repology.org/badge/vertical-allrepos/mdhtml.svg)](https://repology.org/project/mdhtml/versions)

**NOTE:** The version of mdhtml in the nixpkgs repository is on 0.3.0, which is an old version. This version has some bugs and problem with high CPU usage on watch mode. You probably want to use the latest version. 

For NixOS users, I recommend using the nixpkgs-unstable channel, which should have the latest version of mdhtml soon. Or you can install mdhtml through a Nix flake.

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

### Installing from a Nix flake

If you are using Nix, you can install mdhtml from a Nix flake. To do this, you need to have [Nix](https://nixos.org/download.html) installed on your system. Then you can add the mdhtml input to your configuration and install it:

```nix
# flake.nix

{
  inputs.mdhtml.url = "git+https://codeberg.org/Tomkoid/mdhtml";
  # ...

  outputs = {nixpkgs, ...} @ inputs: {
    nixosConfigurations.HOSTNAME = nixpkgs.lib.nixosSystem {
      specialArgs = { inherit inputs; }; # this is the important part
      modules = [
        ./configuration.nix
      ];
    };
  } 
}

# configuration.nix

{inputs, pkgs, ...}: {
  # ...
  environment.systemPackages = [
    inputs.mdhtml.defaultPackage.${system}
  ];
  # ...
}
``` 
