import 'dart:ffi' as ffi;
import 'package:ffi/ffi.dart';

typedef clone_func = ffi.Int32 Function(
  ffi.Pointer<Utf8> url,
  ffi.Pointer<Utf8> directory,
  ffi.Pointer<Utf8> privateKeyFile,
  ffi.Pointer<Utf8> passphrase,
);
typedef Clone = int Function(
  ffi.Pointer<Utf8> url,
  ffi.Pointer<Utf8> directory,
  ffi.Pointer<Utf8> privateKeyFile,
  ffi.Pointer<Utf8> passphrase,
);

void main() {
  final dylib = ffi.DynamicLibrary.open(_getCorrectLibrary());

  final Clone clone = dylib.lookupFunction<clone_func, Clone>('GitClone');
  clone(
    'git@github.com:vhanda/homelab'.toNativeUtf8(),
    'foo'.toNativeUtf8(),
    '/Users/vishesh/.ssh/id_ed25519'.toNativeUtf8(),
    'pass'.toNativeUtf8(),
  );
}

String _getCorrectLibrary() {
  return '/Users/vishesh/src/gitjournal/go-git-dart/gitjournal.so';
}
