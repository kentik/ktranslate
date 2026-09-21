# Tier B: a synthetic SNMP device farm driving the real Discover() code path
# (pkg/inputs/snmp/disco.go) against real respond/reject/drop network behavior --
# see docs/BENCHMARKING_PLAN.md #2.2 for why this needs real machines rather than a
# Go-level fake (silently-dropped vs. actively-rejected packets don't exist on
# loopback).
#
# Uses the standard `nodes = {...}` QEMU/KVM backend (not the newer systemd-nspawn
# `containers = {...}` backend -- tried first, but it unconditionally needs the
# `uid-range` system feature, which needs out-of-repo nix.conf changes; not worth it
# once the real blocker below turned out to be a non-issue).
#
# On Apple Silicon, this test's VMs run natively on the host Mac via apple-virt/HVF --
# NOT inside the `linux-builder` VM that Nix otherwise dispatches cross-system builds
# to. That matters because the linux-builder can't do nested virtualization at all
# (confirmed: nix-darwin's VZ-based builder has no KVM passthrough to its own guests),
# which looked like a hard blocker until realizing the NixOS test framework doesn't
# need it: since nixpkgs#282401 (2024), `runNixOSTest` runs its qemu process and Python
# driver directly on whichever host evaluates it. Evaluating this file's `pkgs` as an
# aarch64-darwin/x86_64-darwin package set (see flake.nix's `checks`) makes the test
# itself execute on that Darwin host's own accelerated qemu -- confirmed locally, see
# nix/tests/minimal-ping.nix -- while only the *build inputs* it depends on
# (`collectorBin`, a real Linux binary) still cross-build via the linux-builder as
# normal. `nix.settings.system-features` needs `nixos-test`/`apple-virt`, which recent
# Nix (2.19+) auto-detects with no nix-darwin config change at all.
{ pkgs, lib, collectorBin
, deviceCount ? 40
, respondFrac ? 0.7
, rejectFrac ? 0.2
}:
let
  vlanId = 1;
  subnetPrefix = "192.168.50";
  collectorIp = "${subnetPrefix}.1";

  respondCount = builtins.floor (deviceCount * respondFrac);
  rejectCount = builtins.floor (deviceCount * rejectFrac);
  # Remainder split evenly: half get an explicit drop container, half are simply
  # never assigned a container at all (an unclaimed address on the shared vlan).
  # Whether these two sub-cases actually behave differently is one of this test's
  # own open questions -- see docs/BENCHMARKING_PLAN.md.
  dropNodeCount = (deviceCount - respondCount - rejectCount) / 2;

  ids = lib.range 1 deviceCount;
  categoryOf = i:
    if i <= respondCount then "respond"
    else if i <= respondCount + rejectCount then "reject"
    else if i <= respondCount + rejectCount + dropNodeCount then "drop"
    else "unclaimed";

  respondIds = lib.filter (i: categoryOf i == "respond") ids;
  namedIds = lib.filter (i: categoryOf i != "unclaimed") ids;
  nameOf = i: "device${toString i}";

  # networking.firewall.rejectPackets (default false = silent DROP) gives the
  # respond/reject/drop distinction natively, no manual iptables needed:
  #   - drop:            leave rejectPackets at its default (false) -- unhandled
  #                       traffic is silently dropped, forcing the full timeout.
  #   - respond & reject: rejectPackets = true -- fast RST/ICMP-unreachable instead.
  #                       Respond nodes need this too, not just services.snmpd.enable
  #                       -- otherwise the pre-scan's TCP-dial-to-port-1 liveness
  #                       probe gets silently dropped on respond nodes as well,
  #                       corrupting the "respond nodes are fast" half of the
  #                       measurement.
  mkDeviceNode = i: cat: {
    # These VMs run one thing (snmpd, or nothing at all) -- the test framework's
    # ~1024 MiB default is wildly oversized and, multiplied across dozens of
    # concurrently-running device VMs, is what actually made a 40-node local run
    # thrash under memory pressure rather than complete. Confirmed by first hitting
    # that thrash at the framework default, then overcorrecting to 192 (too little --
    # a respond-category VM booting net-snmp crashed with a disconnected test-driver
    # shell rather than just running slow), then settling here.
    virtualisation.memorySize = 384;
    virtualisation.vlans = [ vlanId ];
    networking.useDHCP = false;
    networking.firewall.enable = true;
    networking.firewall.rejectPackets = cat == "respond" || cat == "reject";
    networking.interfaces.eth1.ipv4.addresses = [
      { address = "${subnetPrefix}.${toString (i + 10)}"; prefixLength = 24; }
    ];
  } // lib.optionalAttrs (cat == "respond") {
    environment.systemPackages = [ pkgs.net-snmp ];
    services.snmpd = {
      enable = true;
      openFirewall = true; # UDP/161 only -- TCP/1 stays behind rejectPackets above
      configText = "rocommunity public\n";
    };
  };

  deviceNodes = builtins.listToAttrs (map (i: {
    name = nameOf i;
    value = mkDeviceNode i (categoryOf i);
  }) namedIds);

  # Mirrors deployment/docker/snmp-base-nr.yaml -- the shipped example that sets
  # check_all_ips: true on purpose -- with discovery.cidrs pointed at this farm's
  # actual subnet. mib_profile_dir just needs to be a non-empty string
  # (disco.go:53-55); it doesn't need real profile files to avoid erroring, and
  # profile-matching cost is already covered by Tier A (mibs/profile_bench_test.go).
  snmpYaml = pkgs.writeText "snmp.yml" ''
    discovery:
      cidrs:
      - ${subnetPrefix}.0/24
      ports:
      - 161
      default_communities:
      - public
      default_v3: null
      add_devices: true
      add_mibs: true
      threads: 4
      check_all_ips: true
      replace_devices: true
    global:
      poll_time_sec: 300
      drop_if_outside_poll: false
      mib_profile_dir: /etc/ktranslate/profiles
      mibs_db: /etc/ktranslate/mibs.db
      mibs_enabled:
      - IF-MIB
      timeout_ms: 3000
      retries: 0
  '';

  # Built as a flat list of top-level Python statements joined with "\n", rather than
  # one `''...''` heredoc mixing indented literal text with column-0 interpolations --
  # the two disagree on baseline indentation, which broke the test driver's own
  # Python type-checker (dedent ambiguity, not a real Python bug).
  testScriptLines =
    [ "start_all()" "" ]
    ++ map (n: "${n}.wait_for_unit(\"multi-user.target\")")
         (["collector"] ++ map nameOf namedIds)
    ++ map (n: "${n}.wait_for_unit(\"snmpd.service\")") (map nameOf respondIds)
    ++ [
      ""
      "import time"
      ""
      "t0 = time.monotonic()"
      "collector.succeed("
      "    \"ktranslate -snmp=/etc/ktranslate/snmp.yml -snmp_discovery=true \""
      "    \"-snmp_out_file=/tmp/discovered.yml -log_level=info > /tmp/collector.log 2>&1\""
      ")"
      "elapsed = time.monotonic() - t0"
      # disco.go:134's unconditional tail sleep ("give logs time to get sent back") --
      # a fixed tax unrelated to farm size, subtracted so it doesn't skew small runs.
      "elapsed_minus_fixed_sleep = elapsed - 2.0"
      ""
      # addDevices() (disco.go:392+) writes discovered devices as "<name>__<ip>:" keys
      # (disco.go:421,429) into the -snmp_out_file output -- more robust than
      # screen-scraping the SnmpDiscoDeviceStat log line, which is only logged from
      # the timer-loop path (disco.go:222), not this one-shot -snmp_discovery run.
      "device_count = collector.succeed(\"grep -c '__' /tmp/discovered.yml || true\").strip()"
      ""
      # Correctness check, not just a measurement: catches a silent regression (e.g.
      # Discover() erroring out per-device and swallowing it, or a config change that
      # stops matching respond nodes) that wouldn't otherwise fail this test, since
      # collector.succeed() above only checks ktranslate's exit code, not what it
      # actually found.
      "device_count_int = int(device_count)"
      "assert device_count_int == ${toString respondCount}, ("
      "    f\"expected ${toString respondCount} devices (the respond-category node count), \""
      "    f\"discovered {device_count_int}\""
      ")"
      ""
      "result = ("
      "    '{\"elapsed_s\": %s, \"device_count\": \"%s\", \"node_count\": %d}'"
      "    % (elapsed_minus_fixed_sleep, device_count, ${toString deviceCount})"
      ")"
      "collector.succeed(\"echo '%s' > /tmp/result.json\" % result)"
      "collector.copy_from_machine(\"/tmp/result.json\", \"result.json\")"
    ];
in
pkgs.testers.runNixOSTest {
  name = "snmp-discovery-bench";

  nodes = {
    collector = {
      virtualisation.memorySize = 768; # runs the real ktranslate binary, not just snmpd
      virtualisation.vlans = [ vlanId ];
      networking.useDHCP = false;
      networking.firewall.enable = false; # collector must be free to dial out everywhere
      networking.interfaces.eth1.ipv4.addresses = [
        { address = collectorIp; prefixLength = 24; }
      ];
      environment.systemPackages = [ collectorBin pkgs.jq ];
      environment.etc."ktranslate/snmp.yml".source = snmpYaml;
      environment.etc."ktranslate/mibs.db".source = ../../config/mibs.db;
      systemd.tmpfiles.rules = [ "d /etc/ktranslate/profiles 0755 root root -" ];
    };
  } // deviceNodes;

  testScript = lib.concatStringsSep "\n" testScriptLines;
}
