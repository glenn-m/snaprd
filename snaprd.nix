{ lib, buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "snaprd";
  version = "0.10.0";

  src = fetchFromGitHub {
    owner = "glenn-m";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-yoio22D8k4rO8lRLoGSJGl8raMVO9fOGHFobAZngcxw=";
  };

  vendorSha256 = "sha256-RSAT9VtsdXvWDhIZlOjwCF9nhONPXCSEaxVlgW14IKA=";

  meta = with lib; {
    description =
      "Daemon that runs snaprd on a schedule and surfaces Prometheus metrics";
    homepage = "https://github.com/glenn-m/snaprd";
    platforms = platforms.linux;
    license = licenses.mit;
  };
}
