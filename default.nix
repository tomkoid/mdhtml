{ pkgs, ... }:

pkgs.buildGoModule {
  pname = "mdhtml";
  version = "1.0";
  src = ./.;

  vendorHash = null;

  meta = with pkgs.lib; {
    description = "Really simple CLI Markdown to HTML converter with styling support";
    homepage = "https://codeberg.org/Tomkoid/mdhtml";
    license = licenses.mit;
    maintainers = with maintainers; [ tomkoid ];
  };
}
