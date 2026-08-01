{
  description = "AXONE Protocol blockchain development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
  };

  outputs =
    { nixpkgs, ... }:
    let
      supportedSystems = [
        "aarch64-darwin"
        "x86_64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.bash-language-server
              pkgs.buf
              pkgs.git
              pkgs.gnumake
              pkgs.go_1_25
              pkgs.gofumpt
              pkgs.golangci-lint
              pkgs.gomplate
              pkgs.gopls
              pkgs.jq
              pkgs.markdownlint-cli2
              pkgs.marksman
              pkgs.nil
              pkgs.protobuf
              pkgs.protobuf-language-server
              pkgs.protoc-gen-go
              pkgs.protoc-gen-go-grpc
              pkgs.yaml-language-server
            ];

            shellHook = ''
              echo "AXONE development environment loaded"
              echo "Go: $(go version)"
            '';
          };
        }
      );
    };
}
