{ lib, buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "snaprd";
  version = "0.1.0";

  src = fetchFromGitHub {
    owner = "glenn-m";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-1hbkcbj769fmq1vzgb2588jb8dqzzlss7yldwn18jq0n47ah5y5a";
  };

  vendorSha256 = lib.fakeHash;

  meta = with lib; {
    description =
      "Daemon that runs snaprd on a schedule and surfaces Prometheus metrics";
    homepage = "https://github.com/glenn-m/snaprd";
    platforms = platforms.linux;
    license = licenses.mit;
  };
}
