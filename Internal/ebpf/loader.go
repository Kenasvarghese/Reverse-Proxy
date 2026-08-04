//go:build linux

package ebpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang XDP packet_droper.c -- -I.
