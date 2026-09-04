import 'dart:async';
import 'dart:isolate';
import 'dart:typed_data';

import 'package:go_git_dart/go_git_dart.dart';

class GitBindingsAsync {
  final String? _libPath;

  GitBindingsAsync([this._libPath]);

  Future<void> clone(
    String url,
    String directory,
    Uint8List pemBytes,
    String password,
  ) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextCloneRequestId++;
    var request =
        _CloneRequest(requestId, _libPath, url, directory, pemBytes, password);
    var completer = Completer<Exception?>();
    _cloneRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> fetch(
    String remote,
    String directory,
    Uint8List pemBytes,
    String password,
  ) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextFetchRequestId++;
    var request = _FetchRequest(
        requestId, _libPath, remote, directory, pemBytes, password);
    var completer = Completer<Exception?>();
    _fetchRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> pull(
    String remote,
    String directory,
    Uint8List pemBytes,
    String password,
  ) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextPullRequestId++;
    var request = _PullRequest(
        requestId, _libPath, remote, directory, pemBytes, password);
    var completer = Completer<Exception?>();
    _pullRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> push(
    String remote,
    String directory,
    Uint8List pemBytes,
    String password,
  ) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextPushRequestId++;
    var request = _PushRequest(
        requestId, _libPath, remote, directory, pemBytes, password);
    var completer = Completer<Exception?>();
    _pushRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<String> defaultBranch(
    String remoteUrl,
    Uint8List pemBytes,
    String password,
  ) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextDefaultBranchRequestsId++;
    var request = _DefaultBranchRequest(
        requestId, _libPath, remoteUrl, pemBytes, password);
    var completer = Completer<(String?, Exception?)>();
    _defaultBranchRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var result = await completer.future;
    if (result.$2 != null) throw Exception(result.$2!);
    return result.$1!;
  }

  Future<void> add(String directory, String path) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextAddRequestId++;
    var request = _AddRequest(requestId, _libPath, directory, path);
    var completer = Completer<Exception?>();
    _addRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<String> commit(String directory, String message) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextCommitRequestId++;
    var request = _CommitRequest(requestId, _libPath, directory, message);
    var completer = Completer<(String?, Exception?)>();
    _commitRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var result = await completer.future;
    if (result.$2 != null) throw Exception(result.$2!);
    return result.$1!;
  }

  Future<void> remove(String directory, String path) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextRemoveRequestId++;
    var request = _RemoveRequest(requestId, _libPath, directory, path);
    var completer = Completer<Exception?>();
    _removeRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> rm(String directory, String path) {
    return remove(directory, path);
  }

  Future<void> move(String directory, String fromPath, String toPath) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextMoveRequestId++;
    var request =
        _MoveRequest(requestId, _libPath, directory, fromPath, toPath);
    var completer = Completer<Exception?>();
    _moveRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> resetHard(String directory) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextResetHardRequestId++;
    var request = _ResetHardRequest(requestId, _libPath, directory);
    var completer = Completer<Exception?>();
    _resetHardRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> resetHardTo(String directory, String commitHash) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextResetHardToRequestId++;
    var request =
        _ResetHardToRequest(requestId, _libPath, directory, commitHash);
    var completer = Completer<Exception?>();
    _resetHardToRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> switchBranch(String directory, String branch) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextSwitchRequestId++;
    var request = _SwitchRequest(requestId, _libPath, directory, branch);
    var completer = Completer<Exception?>();
    _switchRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }

  Future<void> mergeCurrentBranch(String directory) async {
    var helperIsolateSendPort = await _helperIsolateSendPort;
    var requestId = _nextMergeCurrentBranchRequestId++;
    var request = _MergeCurrentBranchRequest(requestId, _libPath, directory);
    var completer = Completer<Exception?>();
    _mergeCurrentBranchRequests[requestId] = completer;
    helperIsolateSendPort.send(request);
    var ex = await completer.future;
    if (ex != null) throw Exception(ex);
  }
}

class _CloneRequest {
  final int id;
  final String? libPath;
  final String url;
  final String directory;
  final Uint8List pemBytes;
  final String password;

  const _CloneRequest(this.id, this.libPath, this.url, this.directory,
      this.pemBytes, this.password);
}

class _CloneResponse {
  final int id;
  final Exception? exception;

  const _CloneResponse(this.id, this.exception);
}

class _FetchRequest {
  final int id;
  final String? libPath;
  final String remote;
  final String directory;
  final Uint8List pemBytes;
  final String password;

  const _FetchRequest(this.id, this.libPath, this.remote, this.directory,
      this.pemBytes, this.password);
}

class _FetchResponse {
  final int id;
  final Exception? exception;

  const _FetchResponse(this.id, this.exception);
}

class _PullRequest {
  final int id;
  final String? libPath;
  final String remote;
  final String directory;
  final Uint8List pemBytes;
  final String password;

  const _PullRequest(this.id, this.libPath, this.remote, this.directory,
      this.pemBytes, this.password);
}

class _PullResponse {
  final int id;
  final Exception? exception;

  const _PullResponse(this.id, this.exception);
}

class _PushRequest {
  final int id;
  final String? libPath;
  final String remote;
  final String directory;
  final Uint8List pemBytes;
  final String password;

  const _PushRequest(this.id, this.libPath, this.remote, this.directory,
      this.pemBytes, this.password);
}

class _PushResponse {
  final int id;
  final Exception? exception;

  const _PushResponse(this.id, this.exception);
}

class _DefaultBranchRequest {
  final int id;
  final String? libPath;
  final String remoteUrl;
  final Uint8List pemBytes;
  final String password;

  const _DefaultBranchRequest(
      this.id, this.libPath, this.remoteUrl, this.pemBytes, this.password);
}

class _DefaultBranchResponse {
  final int id;
  final String? branch;
  final Exception? exception;

  const _DefaultBranchResponse(this.id, this.branch, this.exception);
}

class _AddRequest {
  final int id;
  final String? libPath;
  final String directory;
  final String path;

  const _AddRequest(this.id, this.libPath, this.directory, this.path);
}

class _AddResponse {
  final int id;
  final Exception? exception;

  const _AddResponse(this.id, this.exception);
}

class _CommitRequest {
  final int id;
  final String? libPath;
  final String directory;
  final String message;

  const _CommitRequest(this.id, this.libPath, this.directory, this.message);
}

class _CommitResponse {
  final int id;
  final String? hash;
  final Exception? exception;

  const _CommitResponse(this.id, this.hash, this.exception);
}

class _RemoveRequest {
  final int id;
  final String? libPath;
  final String directory;
  final String path;

  const _RemoveRequest(this.id, this.libPath, this.directory, this.path);
}

class _RemoveResponse {
  final int id;
  final Exception? exception;

  const _RemoveResponse(this.id, this.exception);
}

class _MoveRequest {
  final int id;
  final String? libPath;
  final String directory;
  final String fromPath;
  final String toPath;

  const _MoveRequest(
      this.id, this.libPath, this.directory, this.fromPath, this.toPath);
}

class _MoveResponse {
  final int id;
  final Exception? exception;

  const _MoveResponse(this.id, this.exception);
}

class _ResetHardRequest {
  final int id;
  final String? libPath;
  final String directory;

  const _ResetHardRequest(this.id, this.libPath, this.directory);
}

class _ResetHardResponse {
  final int id;
  final Exception? exception;

  const _ResetHardResponse(this.id, this.exception);
}

class _ResetHardToRequest {
  final int id;
  final String? libPath;
  final String directory;
  final String commitHash;

  const _ResetHardToRequest(
      this.id, this.libPath, this.directory, this.commitHash);
}

class _ResetHardToResponse {
  final int id;
  final Exception? exception;

  const _ResetHardToResponse(this.id, this.exception);
}

class _SwitchRequest {
  final int id;
  final String? libPath;
  final String directory;
  final String branch;

  const _SwitchRequest(this.id, this.libPath, this.directory, this.branch);
}

class _SwitchResponse {
  final int id;
  final Exception? exception;

  const _SwitchResponse(this.id, this.exception);
}

class _MergeCurrentBranchRequest {
  final int id;
  final String? libPath;
  final String directory;

  const _MergeCurrentBranchRequest(this.id, this.libPath, this.directory);
}

class _MergeCurrentBranchResponse {
  final int id;
  final Exception? exception;

  const _MergeCurrentBranchResponse(this.id, this.exception);
}

int _nextCloneRequestId = 0;
final _cloneRequests = <int, Completer<Exception?>>{};

int _nextFetchRequestId = 0;
final _fetchRequests = <int, Completer<Exception?>>{};

int _nextPullRequestId = 0;
final _pullRequests = <int, Completer<Exception?>>{};

int _nextPushRequestId = 0;
final _pushRequests = <int, Completer<Exception?>>{};

int _nextDefaultBranchRequestsId = 0;
final _defaultBranchRequests = <int, Completer<(String?, Exception?)>>{};

int _nextAddRequestId = 0;
final _addRequests = <int, Completer<Exception?>>{};

int _nextCommitRequestId = 0;
final _commitRequests = <int, Completer<(String?, Exception?)>>{};

int _nextRemoveRequestId = 0;
final _removeRequests = <int, Completer<Exception?>>{};

int _nextMoveRequestId = 0;
final _moveRequests = <int, Completer<Exception?>>{};

int _nextResetHardRequestId = 0;
final _resetHardRequests = <int, Completer<Exception?>>{};

int _nextResetHardToRequestId = 0;
final _resetHardToRequests = <int, Completer<Exception?>>{};

int _nextSwitchRequestId = 0;
final _switchRequests = <int, Completer<Exception?>>{};

int _nextMergeCurrentBranchRequestId = 0;
final _mergeCurrentBranchRequests = <int, Completer<Exception?>>{};

Future<SendPort> _helperIsolateSendPort = () async {
  final Completer<SendPort> completer = Completer<SendPort>();

  final ReceivePort receivePort = ReceivePort()
    ..listen((dynamic data) {
      if (data is SendPort) {
        completer.complete(data);
        return;
      }
      if (data is _CloneResponse) {
        final completer = _cloneRequests[data.id]!;
        _cloneRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _FetchResponse) {
        final completer = _fetchRequests[data.id]!;
        _fetchRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _PullResponse) {
        final completer = _pullRequests[data.id]!;
        _pullRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _PushResponse) {
        final completer = _pushRequests[data.id]!;
        _pushRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _DefaultBranchResponse) {
        final completer = _defaultBranchRequests[data.id]!;
        _defaultBranchRequests.remove(data.id);
        completer.complete((data.branch, data.exception));
        return;
      }
      if (data is _AddResponse) {
        final completer = _addRequests[data.id]!;
        _addRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _CommitResponse) {
        final completer = _commitRequests[data.id]!;
        _commitRequests.remove(data.id);
        completer.complete((data.hash, data.exception));
        return;
      }
      if (data is _RemoveResponse) {
        final completer = _removeRequests[data.id]!;
        _removeRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _MoveResponse) {
        final completer = _moveRequests[data.id]!;
        _moveRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _ResetHardResponse) {
        final completer = _resetHardRequests[data.id]!;
        _resetHardRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _ResetHardToResponse) {
        final completer = _resetHardToRequests[data.id]!;
        _resetHardToRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _SwitchResponse) {
        final completer = _switchRequests[data.id]!;
        _switchRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      if (data is _MergeCurrentBranchResponse) {
        final completer = _mergeCurrentBranchRequests[data.id]!;
        _mergeCurrentBranchRequests.remove(data.id);
        completer.complete(data.exception);
        return;
      }
      throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
    });

  // Start the helper isolate.
  await Isolate.spawn((SendPort sendPort) async {
    final ReceivePort helperReceivePort = ReceivePort()
      ..listen((dynamic data) {
        if (data is _CloneRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.clone(data.url, data.directory, data.pemBytes, data.password);
            sendPort.send(_CloneResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_CloneResponse(data.id, e));
          }
          return;
        }
        if (data is _FetchRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.fetch(
                data.remote, data.directory, data.pemBytes, data.password);
            sendPort.send(_FetchResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_FetchResponse(data.id, e));
          }
          return;
        }
        if (data is _PullRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.pull(
                data.remote, data.directory, data.pemBytes, data.password);
            sendPort.send(_PullResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_PullResponse(data.id, e));
          }
          return;
        }
        if (data is _PushRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.push(
                data.remote, data.directory, data.pemBytes, data.password);
            sendPort.send(_PushResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_PushResponse(data.id, e));
          }
          return;
        }
        if (data is _DefaultBranchRequest) {
          try {
            var repo = GitBindings(data.libPath);
            var branch = repo.defaultBranch(
                data.remoteUrl, data.pemBytes, data.password);
            sendPort.send(_DefaultBranchResponse(data.id, branch, null));
          } on Exception catch (e) {
            sendPort.send(_DefaultBranchResponse(data.id, null, e));
          }
          return;
        }
        if (data is _AddRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.add(data.directory, data.path);
            sendPort.send(_AddResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_AddResponse(data.id, e));
          }
          return;
        }
        if (data is _CommitRequest) {
          try {
            var repo = GitBindings(data.libPath);
            var hash = repo.commit(data.directory, data.message);
            sendPort.send(_CommitResponse(data.id, hash, null));
          } on Exception catch (e) {
            sendPort.send(_CommitResponse(data.id, null, e));
          }
          return;
        }
        if (data is _RemoveRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.remove(data.directory, data.path);
            sendPort.send(_RemoveResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_RemoveResponse(data.id, e));
          }
          return;
        }
        if (data is _MoveRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.move(data.directory, data.fromPath, data.toPath);
            sendPort.send(_MoveResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_MoveResponse(data.id, e));
          }
          return;
        }
        if (data is _ResetHardRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.resetHard(data.directory);
            sendPort.send(_ResetHardResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_ResetHardResponse(data.id, e));
          }
          return;
        }
        if (data is _ResetHardToRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.resetHardTo(data.directory, data.commitHash);
            sendPort.send(_ResetHardToResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_ResetHardToResponse(data.id, e));
          }
          return;
        }
        if (data is _SwitchRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.switchBranch(data.directory, data.branch);
            sendPort.send(_SwitchResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_SwitchResponse(data.id, e));
          }
          return;
        }
        if (data is _MergeCurrentBranchRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.mergeCurrentBranch(data.directory);
            sendPort.send(_MergeCurrentBranchResponse(data.id, null));
          } on Exception catch (e) {
            sendPort.send(_MergeCurrentBranchResponse(data.id, e));
          }
          return;
        }
        throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
      });

    // Send the port to the main isolate on which we can receive requests.
    sendPort.send(helperReceivePort.sendPort);
  }, receivePort.sendPort);

  return completer.future;
}();
