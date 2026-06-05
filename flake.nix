{
  description = "porttcd dev environment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        python = pkgs.python312;
        pythonEnv = python.withPackages (ps: with ps; [
          pip
          virtualenv
        ]);
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go_1_26
            pythonEnv
            pkgs.curlFull
            pkgs.git
            pkgs.containerlab
            pkgs.qemu
            pkgs.docker-client
            pkgs.upx
            pkgs.goreleaser
            pkgs.nfpm
          ];
          shellHook = ''
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
            export PIP_PREFIX=$(pwd)/.pip
            export PYTHONPATH="$PIP_PREFIX/lib/python3.11/site-packages:$PYTHONPATH"
            export PATH="$PIP_PREFIX/bin:$PATH"
            if [ ! -d .venv ]; then
              echo "creating virtualenv..."
              virtualenv .venv
            fi
            source .venv/bin/activate
            echo "python: $(python --version)"
            echo "pip:    $(pip --version)"
            echo "go:     $(go version)"
            echo "upx:    $(upx --version | head -1)"
            echo "goreleaser: $(goreleaser --version | head -1)"
          '';
        };
      });
}