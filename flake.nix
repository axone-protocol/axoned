{
  description = "AXONE Protocol blockchain development environment";

  inputs = {
    cosmovisor = {
      url = "github:cosmos/cosmos-sdk/tools/cosmovisor/v1.7.1";
      flake = false;
    };
    heighliner = {
      url = "github:strangelove-ventures/heighliner/v1.7.4";
      flake = false;
    };
    mock = {
      url = "github:uber/mock/v0.6.0";
      flake = false;
    };
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    nixpkgsGo124.url = "github:NixOS/nixpkgs/nixos-25.11";
  };

  outputs =
    {
      self,
      cosmovisor,
      heighliner,
      mock,
      nixpkgs,
      nixpkgsGo124,
      ...
    }:
    let
      supportedSystems = [
        "aarch64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          go124 = nixpkgsGo124.legacyPackages.${system}.go_1_24;
        in
        {
          cosmovisor = (pkgs.buildGoModule.override { go = go124; }) {
            pname = "cosmovisor";
            version = "1.7.1";
            src = cosmovisor;
            modRoot = "tools/cosmovisor";
            subPackages = [ "cmd/cosmovisor" ];
            env.CGO_ENABLED = 0;
            vendorHash = "sha256-DXgFvjm1fDHtDwPNfLIPi2vMdZTi3bYN7cgR1dgdgLk=";
          };
          heighliner = pkgs.buildGoModule {
            pname = "heighliner";
            version = "1.7.4";
            src = heighliner;
            vendorHash = "sha256-Sa/lCa7IFcnIGNKCPvFeQke6dNOycK6vjg5gmHOOdic=";
          };
          mockgen = pkgs.buildGoModule {
            pname = "mockgen";
            version = "0.6.0";
            src = mock;
            subPackages = [ "mockgen" ];
            vendorHash = "sha256-Cf7lKfMuPFT/I1apgChUNNCG2C7SrW7ncF8OusbUs+A=";
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              self.packages.${system}.cosmovisor
              self.packages.${system}.heighliner
              self.packages.${system}.mockgen
              pkgs.act
              pkgs.actionlint
              pkgs.bash-language-server
              pkgs.buf
              pkgs.deadnix
              pkgs.docker-client
              pkgs.git
              pkgs.gnumake
              pkgs.go_1_25
              pkgs.gofumpt
              pkgs.golangci-lint
              pkgs.gh
              pkgs.gomplate
              pkgs.gopls
              pkgs.jq
              pkgs.markdownlint-cli2
              pkgs.marksman
              pkgs.nil
              pkgs.nixfmt
              pkgs.nodejs_24
              pkgs.protobuf
              pkgs.protobuf-language-server
              pkgs.protoc-gen-go
              pkgs.protoc-gen-go-grpc
              pkgs.statix
              pkgs.uv
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
