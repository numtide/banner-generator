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

    # Deploy
    pkgs.flyctl
  ]
  # SVG to PNG conversion uses headless Chrome. Chromium is Linux-only: on macOS, set CHROME_PATH.
  ++ pkgs.lib.optional pkgs.stdenv.hostPlatform.isLinux pkgs.chromium;

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
