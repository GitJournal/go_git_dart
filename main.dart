import 'dart:ffi' as ffi;
import 'dart:ffi';
import 'dart:typed_data';
import 'package:ffi/ffi.dart';
import 'dart:io';

import 'generated_bindings.dart';

class GitBindings {
  late final NativeLibrary lib;

  GitBindings(String libraryPath) {
    lib = NativeLibrary(ffi.DynamicLibrary.open(libraryPath));
  }

  void clone(
    String url,
    String directory,
    Uint8List pemBytes,
    String password,
  ) {
    var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
    for (var i = 0; i < pemBytes.length; i++) {
      cPemBytes[i] = pemBytes[i];
    }

    var cloneUrl = url.toNativeUtf8();
    var cloneDir = directory.toNativeUtf8();
    var pemPassphrase = password.toNativeUtf8();

    var retValue = lib.GitClone(
      cloneUrl.cast<Char>(),
      cloneDir.cast<Char>(),
      cPemBytes.cast<Char>(),
      pemBytes.length,
      pemPassphrase.cast<Char>(),
    );
    if (retValue != 0) {
      throw Exception("GitClone failed with error code: $retValue");
    }

    malloc.free(cPemBytes);
    malloc.free(cloneUrl);
    malloc.free(cloneDir);
    malloc.free(pemPassphrase);
  }
}

void main() {
  final pemBytes = File('/Users/vishesh/.ssh/id_ed25519').readAsBytesSync();
  var cloneUrl = 'git@github.com:vhanda/homelab';
  var cloneDir = 'foo';
  var pemPassphrase = 'vishu002';

  var bindings = GitBindings(_getCorrectLibrary());
  bindings.clone(cloneUrl, cloneDir, pemBytes, pemPassphrase);
}

String _getCorrectLibrary() {
  return '/Users/vishesh/src/gitjournal/go-git-dart/gitjournal.so';
}
