{
  description = "go-random development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      forEachSystem = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forEachSystem (pkgs: {
        default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.golangci-lint
            pkgs.govulncheck
            pkgs.just
            pkgs.lefthook
            pkgs.reuse
          ];
          shellHook = ''
            lefthook install >/dev/null
          '';
        };
      });

      formatter = forEachSystem (pkgs: pkgs.nixfmt);
    };
}
