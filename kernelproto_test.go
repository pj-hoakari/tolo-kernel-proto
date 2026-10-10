package kernelproto_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	kernelproto "github.com/pj-hoakari/tolo-kernel-proto"
	kernelv1 "github.com/pj-hoakari/tolo-kernel-proto/gen/tolo/kernel/v1"
)

func TestWriteProtoWritesRegisteredPath(t *testing.T) {
	root := t.TempDir()

	written, err := kernelproto.WriteProto(root)
	if err != nil {
		t.Fatal(err)
	}

	registered := kernelv1.File_tolo_kernel_v1_kernel_proto.Path()
	want := filepath.Join(root, filepath.FromSlash(registered))

	if len(written) != 1 || written[0] != want {
		t.Fatalf("written = %v, want [%s]", written, want)
	}

	got, err := os.ReadFile(want)
	if err != nil {
		t.Fatal(err)
	}

	source, err := os.ReadFile(filepath.Join("proto", filepath.FromSlash(registered)))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, source) {
		t.Fatalf("%s differs from proto/%s", want, registered)
	}
}

func TestWriteProtoOverwritesStaleFile(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, "tolo", "kernel", "v1", "kernel.proto")

	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filename, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := kernelproto.WriteProto(root); err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(filename)
	want, _ := fs.ReadFile(kernelproto.Files(), "tolo/kernel/v1/kernel.proto")

	if !bytes.Equal(got, want) {
		t.Fatal("stale file was not overwritten")
	}
}
