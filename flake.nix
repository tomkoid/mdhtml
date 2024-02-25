{
  description = "mdhtml flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };

  outputs = { self, nixpkgs, flake-utils, rust-overlay }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import rust-overlay) ];
        pkgs = import nixpkgs { inherit system; inherit overlays; };
      in
      {
        defaultPackage = (import ./default.nix { inherit pkgs; });
        devShells.default = pkgs.mkShell
          {
            buildInputs = with pkgs; [
              rust-bin.stable.latest.default
            ];
          };
      }
    );
}
