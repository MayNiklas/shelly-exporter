{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  buildInputs = with pkgs.python39Packages; [ prometheus-client ];

  shellHook = ''
    # ...
  '';
}
