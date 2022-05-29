{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  buildInputs = [ go gcc ];

  shellHook = ''
    # ...
  '';
}
