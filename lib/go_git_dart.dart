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

  void pull(
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

    var retValue = lib.GitPull(
      remoteName.cast<Char>(),
      cloneDir.cast<Char>(),
      cPemBytes.cast<Char>(),
      pemBytes.length,
      pemPassphrase.cast<Char>(),
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitPull failed with error code: $err");
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

  void add(String directory, String path) {
    var repoDir = directory.toNativeUtf8();
    var repoPath = path.toNativeUtf8();

    var retValue = lib.GitAdd(
      repoDir.cast<Char>(),
      repoPath.cast<Char>(),
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitAdd failed with error: $err");
    }

    malloc.free(repoDir);
    malloc.free(repoPath);
  }

  String commit(String directory, String message) {
    var repoDir = directory.toNativeUtf8();
    var commitMessage = message.toNativeUtf8();
    var outputHash = malloc.allocate<Pointer<Char>>(sizeOf<Pointer<Char>>());

    var retValue = lib.GitCommit(
      repoDir.cast<Char>(),
      commitMessage.cast<Char>(),
      outputHash,
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());
      malloc.free(repoDir);
      malloc.free(commitMessage);
      malloc.free(outputHash);

      throw Exception("GitCommit failed with error: $err");
    }

    malloc.free(repoDir);
    malloc.free(commitMessage);

    var hash = outputHash.value.cast<Utf8>().toDartString();
    lib.free(outputHash.value.cast());
    malloc.free(outputHash);

    return hash;
  }

  void remove(String directory, String path) {
    var repoDir = directory.toNativeUtf8();
    var repoPath = path.toNativeUtf8();

    var retValue = lib.GitRemove(
      repoDir.cast<Char>(),
      repoPath.cast<Char>(),
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitRemove failed with error: $err");
    }

    malloc.free(repoDir);
    malloc.free(repoPath);
  }

  void rm(String directory, String path) {
    remove(directory, path);
  }

  void move(String directory, String fromPath, String toPath) {
    var repoDir = directory.toNativeUtf8();
    var fromRepoPath = fromPath.toNativeUtf8();
    var toRepoPath = toPath.toNativeUtf8();

    var retValue = lib.GitMove(
      repoDir.cast<Char>(),
      fromRepoPath.cast<Char>(),
      toRepoPath.cast<Char>(),
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitMove failed with error: $err");
    }

    malloc.free(repoDir);
    malloc.free(fromRepoPath);
    malloc.free(toRepoPath);
  }

  void resetHard(String directory) {
    var repoDir = directory.toNativeUtf8();

    var retValue = lib.GitResetHard(repoDir.cast<Char>());
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitResetHard failed with error: $err");
    }

    malloc.free(repoDir);
  }

  void resetHardTo(String directory, String commitHash) {
    var repoDir = directory.toNativeUtf8();
    var hash = commitHash.toNativeUtf8();

    var retValue = lib.GitResetHardTo(
      repoDir.cast<Char>(),
      hash.cast<Char>(),
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitResetHardTo failed with error: $err");
    }

    malloc.free(repoDir);
    malloc.free(hash);
  }

  void switchBranch(String directory, String branch) {
    var repoDir = directory.toNativeUtf8();
    var branchName = branch.toNativeUtf8();

    var retValue = lib.GitSwitch(
      repoDir.cast<Char>(),
      branchName.cast<Char>(),
    );
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitSwitch failed with error: $err");
    }

    malloc.free(repoDir);
    malloc.free(branchName);
  }

  void mergeCurrentBranch(String directory) {
    var repoDir = directory.toNativeUtf8();

    var retValue = lib.GitMergeCurrentBranch(repoDir.cast<Char>());
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GitMergeCurrentBranch failed with error: $err");
    }

    malloc.free(repoDir);
  }

  (String, String) generateRsaKeys() {
    var outputPublicKey = malloc.allocate<Pointer<Char>>(0);
    var outputPrivateKey = malloc.allocate<Pointer<Char>>(0);

    var retValue = lib.GJGenerateRSAKeys(outputPublicKey, outputPrivateKey);
    if (retValue != nullptr) {
      var err = retValue.cast<Utf8>().toDartString();
      lib.free(retValue.cast());

      throw Exception("GenerateRsaKeys failed with error code: $err");
    }

    var publicKey = outputPublicKey.value.cast<Utf8>().toDartString();
    lib.free(outputPublicKey.value.cast());
    malloc.free(outputPublicKey);

    var privateKey = outputPrivateKey.value.cast<Utf8>().toDartString();
    lib.free(outputPrivateKey.value.cast());
    malloc.free(outputPrivateKey);

    return (publicKey, privateKey);
  }
}
