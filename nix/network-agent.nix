# Builds the ktranslate binary via buildGoModule directly, not by shelling out to `make`
# (which may go away) -- everything `make all` does is already native here. No
# libpcap/pkg-config needed since the SYN scanner is pure Go (upstream #14).
#
# Version comes from the checked-in VERSION file, not flake.nix's self.rev: VERSION only
# changes on a real bump, so stamping it doesn't invalidate this derivation's Nix cache on
# every commit the way self.rev (or self.lastModifiedDate, for a date) would. `root` below
# is a literal relative path rather than `self` because `self` is an attrset, not a Nix
# `path`, and lib.fileset needs a real path.
#
# `src` is filtered to go.mod/go.sum/VERSION/*.go via lib.fileset, so unrelated changes
# (docs, workflows, etc.) don't force a rebuild. Also Tier B's VM test fixture -- see
# snmp-discovery-bench.nix.
{ pkgs }:

let
  root = ../.;
  fs = pkgs.lib.fileset;
  version = pkgs.lib.strings.trim (builtins.readFile (root + "/VERSION"));
  goSrc = fs.toSource {
    inherit root;
    fileset = fs.unions [
      (root + "/go.mod")
      (root + "/go.sum")
      (root + "/VERSION")
      (fs.fileFilter (file: file.hasExt "go") root)
    ];
  };
in

pkgs.buildGoModule {
  pname = "ktranslate";
  src = goSrc;
  inherit version;

  vendorHash = "sha256-ZQUnUlWTspAZMO90kEJ6+xukw3gX10+IwTegaCUtEo0=";
  subPackages = [ "cmd/ktranslate" ];
  env.CGO_ENABLED = "0";
  ldflags = [ "-X=github.com/kentik/ktranslate/pkg/version.versionStr=${version}" ];

  doCheck = false; # only needs to run, not pass go test
}
