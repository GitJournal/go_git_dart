{ pkgs, ... }:

{
  packages = [ pkgs.git ];

  languages.go.enable = true;
  languages.dart.enable = true;
}
