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
  final dylib = ffi.DynamicLibrary.open(_getCorrectLibrary());
  final Clone clone = dylib.lookupFunction<clone_func, Clone>('GitClone');

  final pemBytes = File('/Users/vishesh/.ssh/id_ed25519').readAsBytesSync();

  var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
  for (var i = 0; i < pemBytes.length; i++) {
    cPemBytes[i] = pemBytes[i];
  }

  var cloneUrl = 'git@github.com:vhanda/homelab'.toNativeUtf8();
  var cloneDir = 'foo'.toNativeUtf8();
  var pemPassphrase = 'vishu002'.toNativeUtf8();

  var retValue = clone(
    cloneUrl,
    cloneDir,
    cPemBytes,
    pemBytes.length,
    pemPassphrase,
  );
  print("ReturnValue: $retValue");

  malloc.free(cPemBytes);
  malloc.free(cloneUrl);
  malloc.free(cloneDir);
  malloc.free(pemPassphrase);
}

String _getCorrectLibrary() {
  return '/Users/vishesh/src/gitjournal/go-git-dart/gitjournal.so';
}
