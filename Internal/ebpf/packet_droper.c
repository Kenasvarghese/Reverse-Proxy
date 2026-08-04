//go:build ignore
// xdp_block_ip.bpf.c

#include "vmlinux.h"

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define ETH_P_IP 0x0800

char LICENSE[] SEC("license") = "GPL";

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __uint(max_entries, 1024);
  __type(key, __u32);  // IPv4 address in network byte order
  __type(value, __u8); // dummy value
} blocked_ips SEC(".maps");

SEC("xdp")
int xdp_block_ip(struct xdp_md *ctx) {
  void *data = (void *)(long)ctx->data;
  void *data_end = (void *)(long)ctx->data_end;

  struct ethhdr *eth = data;
  if ((void *)(eth + 1) > data_end)
    return XDP_PASS;

  if (bpf_ntohs(eth->h_proto) != ETH_P_IP)
    return XDP_PASS;

  struct iphdr *ip = (void *)(eth + 1);
  if ((void *)(ip + 1) > data_end)
    return XDP_PASS;

  __u32 src = ip->saddr;

  if (bpf_map_lookup_elem(&blocked_ips, &src))
    return XDP_DROP;

  return XDP_PASS;
}
