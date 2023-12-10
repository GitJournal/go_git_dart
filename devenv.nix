{ pkgs, ... }:

{
  packages = [ pkgs.git pkgs.cmake ];

  languages.go.enable = true;
}
