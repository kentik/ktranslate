# Fast (~15-20s) sanity check, kept permanently as a regression guard: confirms a
# NixOS VM test can run at all on the current host -- natively via apple-virt/HVF on
# Darwin, kvm on Linux -- before debugging failures in the much larger, slower real
# test below. See nix/tests/snmp-discovery-bench.nix's header comment and
# docs/BENCHMARKING_PLAN.md for why this distinction (native execution vs. dispatching
# to a remote Linux builder, which can't do nested virtualization on Apple Silicon)
# matters and was not obvious going in.
{
  name = "minimal-ping";

  nodes = {
    machine1 = { ... }: { };
    machine2 = { ... }: { };
  };

  testScript = ''
    start_all()
    machine1.systemctl("start network-online.target")
    machine2.systemctl("start network-online.target")
    machine1.wait_for_unit("network-online.target")
    machine2.wait_for_unit("network-online.target")
    machine1.succeed("ping -c 1 machine2")
    machine2.succeed("ping -c 1 machine1")
  '';
}
