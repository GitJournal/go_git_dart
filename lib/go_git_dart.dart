import 'dart:ffi';
import 'dart:ffi' as ffi;

import 'dart:io';

import 'dart:typed_data';
import 'package:ffi/ffi.dart';

import 'go_git_dart_bindings_generated.dart';

const String _libName = 'go_git_dart';

class GitBindings {
  late final GoGitDartBindings lib;

  GitBindings([String? libPath]) {
    var dylib = () {
      if (libPath != null) {
        return DynamicLibrary.open(libPath);
      }
      if (Platform.isMacOS || Platform.isIOS) {
        return DynamicLibrary.open('$_libName.framework/$_libName');
      }
      if (Platform.isAndroid || Platform.isLinux) {
        return DynamicLibrary.open('$_libName.so');
      }
      if (Platform.isWindows) {
        return DynamicLibrary.open('$_libName.dll');
      }
      throw UnsupportedError('Unknown platform: ${Platform.operatingSystem}');
    }();

    lib = GoGitDartBindings(dylib);
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
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitClone failed with error: $err");
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
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitFetch failed with error code: $err");
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
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitPush failed with error code: $err");
    }

    malloc.free(cPemBytes);
    malloc.free(remoteName);
    malloc.free(cloneDir);
    malloc.free(pemPassphrase);
  }

  String defaultBranch(
    String remoteUrl,
    Uint8List pemBytes,
    String password,
  ) {
    var cPemBytes = malloc.allocate<ffi.Uint8>(pemBytes.length);
    for (var i = 0; i < pemBytes.length; i++) {
      cPemBytes[i] = pemBytes[i];
    }

    var remoteUrlN = remoteUrl.toNativeUtf8();
    var pemPassphrase = password.toNativeUtf8();
    var outputBranch = malloc.allocate<Pointer<Char>>(0);

    var retValue = lib.GitDefaultBranch(
      remoteUrlN.cast<Char>(),
      cPemBytes.cast<Char>(),
      pemBytes.length,
      pemPassphrase.cast<Char>(),
      outputBranch,
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitPush failed with error code: $err");
    }

    malloc.free(cPemBytes);
    malloc.free(remoteUrlN);
    malloc.free(pemPassphrase);

    var branch = outputBranch.value.cast<Utf8>().toDartString();
    lib.free(outputBranch.value.cast());
    malloc.free(outputBranch);

    return branch;
  }
}
