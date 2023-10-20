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

  void fetch(
    String remote,
    String directory,
    Uint8List pemBytes,
    String password,
  ) {
    var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
    for (var i = 0; i < pemBytes.length; i++) {
      cPemBytes[i] = pemBytes[i];
    }

    var remoteName = remote.toNativeUtf8();
    var cloneDir = directory.toNativeUtf8();
    var pemPassphrase = password.toNativeUtf8();

    var retValue = lib.GitFetch(
      remoteName.cast<Char>(),
      cloneDir.cast<Char>(),
      cPemBytes.cast<Char>(),
      pemBytes.length,
      pemPassphrase.cast<Char>(),
    );
    if (retValue != 0) {
      throw Exception("GitFetch failed with error code: $retValue");
    }

    malloc.free(cPemBytes);
    malloc.free(remoteName);
    malloc.free(cloneDir);
    malloc.free(pemPassphrase);
  }

  void push(
    String remote,
    String directory,
    Uint8List pemBytes,
    String password,
  ) {
    var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
    for (var i = 0; i < pemBytes.length; i++) {
      cPemBytes[i] = pemBytes[i];
    }

    var remoteName = remote.toNativeUtf8();
    var cloneDir = directory.toNativeUtf8();
    var pemPassphrase = password.toNativeUtf8();

    var retValue = lib.GitPush(
      remoteName.cast<Char>(),
      cloneDir.cast<Char>(),
      cPemBytes.cast<Char>(),
      pemBytes.length,
      pemPassphrase.cast<Char>(),
    );
    if (retValue != 0) {
      throw Exception("GitPush failed with error code: $retValue");
    }

    malloc.free(cPemBytes);
    malloc.free(remoteName);
    malloc.free(cloneDir);
    malloc.free(pemPassphrase);
  }

  void defaultBranch(
    String remote,
    String directory,
    Uint8List pemBytes,
    String password,
  ) {
    var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
    for (var i = 0; i < pemBytes.length; i++) {
      cPemBytes[i] = pemBytes[i];
    }

    var remoteName = remote.toNativeUtf8();
    var cloneDir = directory.toNativeUtf8();
    var pemPassphrase = password.toNativeUtf8();

    var retValue = lib.GitDefaultBranch(
      remoteName.cast<Char>(),
      cloneDir.cast<Char>(),
      cPemBytes.cast<Char>(),
      pemBytes.length,
      pemPassphrase.cast<Char>(),
    );
    if (retValue != 0) {
      throw Exception("GitDefaultBranch failed with error code: $retValue");
    }

    malloc.free(cPemBytes);
    malloc.free(remoteName);
    malloc.free(cloneDir);
    malloc.free(pemPassphrase);
  }
}

String _getCorrectLibrary() {
  return '/Users/vishesh/src/gitjournal/go-git-dart/gitjournal.so';
}

void main(List<String> arguments) {
  if (arguments.isEmpty) {
    print('Please provide a command: clone, fetch, push, defaultBranch');
    return;
  }

  final command = arguments[0];
  final bindings = GitBindings(_getCorrectLibrary());

  switch (command) {
    case 'clone':
      if (arguments.length != 5) {
        print('Usage: clone <url> <directory> <pemFile> <pemPassword>');
        return;
      }
      final pemBytes = File(arguments[3]).readAsBytesSync();
      bindings.clone(arguments[1], arguments[2], pemBytes, arguments[4]);
      break;
    case 'fetch':
    case 'push':
    case 'defaultBranch':
      if (arguments.length != 4) {
        print('Usage: $command <remote> <pemFile> <pemPassword>');
        return;
      }
      final directory = Directory.current.path;
      print(directory);
      final pemBytes = File(arguments[2]).readAsBytesSync();

      switch (command) {
        case 'fetch':
          bindings.fetch(arguments[1], directory, pemBytes, arguments[3]);
          break;
        case 'push':
          bindings.push(arguments[1], directory, pemBytes, arguments[3]);
          break;
        case 'defaultBranch':
          bindings.defaultBranch(
              arguments[1], directory, pemBytes, arguments[3]);
          break;
      }
      break;
    default:
      print('Unknown command: $command');
      break;
  }
}
