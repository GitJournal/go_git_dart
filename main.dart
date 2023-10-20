import 'dart:ffi' as ffi;
import 'dart:ffi';
import 'package:ffi/ffi.dart';
import 'dart:io';

import 'generated_bindings.dart';

void main() {
  var lib = NativeLibrary(ffi.DynamicLibrary.open(_getCorrectLibrary()));

  final pemBytes = File('/Users/vishesh/.ssh/id_ed25519').readAsBytesSync();

  var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
  for (var i = 0; i < pemBytes.length; i++) {
    cPemBytes[i] = pemBytes[i];
  }

  var cloneUrl = 'git@github.com:vhanda/homelab'.toNativeUtf8();
  var cloneDir = 'foo'.toNativeUtf8();
  var pemPassphrase = 'vishu002'.toNativeUtf8();

  var retValue = lib.GitClone(
    cloneUrl.cast<Char>(),
    cloneDir.cast<Char>(),
    cPemBytes.cast<Char>(),
    pemBytes.length,
    pemPassphrase.cast<Char>(),
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
