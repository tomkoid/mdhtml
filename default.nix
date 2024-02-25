{ pkgs, ... }:

pkgs.rustPlatform.buildRustPackage rec {
  pname = "mdhtml";
  version = "2.0";
  src = pkgs.lib.cleanSource ./.;

  cargoLock.lockFile = ./Cargo.lock;
  cargoHash = "sha256-KatRpdro7sMMK96NCyJPcieSIlJ+7edBfvJcEBX5qt8=";

  meta = with pkgs.lib; {
    description = "Really simple CLI Markdown to HTML converter with styling support";
    homepage = "https://codeberg.org/Tomkoid/mdhtml";
    license = licenses.mit;
    maintainers = with maintainers; [ tomkoid ];
  };
}
