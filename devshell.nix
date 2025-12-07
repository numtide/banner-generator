{ pkgs, inputs }:
pkgs.mkShell {
  # Add build dependencies
  packages = [
    # Go compiler and standard tools
    pkgs.go
    pkgs.gopls
    pkgs.go-tools
    pkgs.golangci-lint
    pkgs.delve

    # Additional useful Go tools
    pkgs.gomodifytags
    pkgs.gotests
    pkgs.impl
    pkgs.go-outline

    # Development tools
    inputs.gorefresh.packages.${pkgs.system}.default

    # SVG to PNG conversion (CLI only, uses headless Chrome)
    pkgs.chromium

    # Deploy
    pkgs.flyctl
  ];

  # Add environment variables
  env = {
    # Ensure Go modules are enabled
    GO111MODULE = "on";
  };

  # Load custom bash code
  shellHook = ''
    echo "Go development environment loaded"
    echo "Go version: $(go version)"
  '';
}
