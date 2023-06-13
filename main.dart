import 'dart:ffi' as ffi;
import 'package:ffi/ffi.dart';
import 'dart:io';

typedef clone_func = ffi.Int32 Function(
  ffi.Pointer<Utf8> url,
  ffi.Pointer<Utf8> directory,
  ffi.Pointer<ffi.Uint8> privateKeyBytes,
  ffi.Int privateKeyLength,
  ffi.Pointer<Utf8> passphrase,
);
typedef Clone = int Function(
  ffi.Pointer<Utf8> url,
  ffi.Pointer<Utf8> directory,
  ffi.Pointer<ffi.Uint8> privateKeyBytes,
  int privateKeyLength,
  ffi.Pointer<Utf8> passphrase,
);

void main() {
  final pemBytes = File('/Users/vishesh/.ssh/id_ed25519').readAsBytesSync();
  final dylib = ffi.DynamicLibrary.open(_getCorrectLibrary());

  var bytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
  for (var i = 0; i < pemBytes.length; i++) {
    bytes[i] = pemBytes[i];
  }

  final Clone clone = dylib.lookupFunction<clone_func, Clone>('GitClone');
  clone(
    'git@github.com:vhanda/homelab'.toNativeUtf8(),
    'foo'.toNativeUtf8(),
    bytes,
    pemBytes.length,
    'pass'.toNativeUtf8(),
  );

  malloc.free(bytes);
}

String _getCorrectLibrary() {
  return '/Users/vishesh/src/gitjournal/go-git-dart/gitjournal.so';
}
