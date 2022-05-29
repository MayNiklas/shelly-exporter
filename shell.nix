{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  buildInputs = [ go_1_17 gcc ];

  shellHook = ''
    go run main.go
    exit
  '';
}
