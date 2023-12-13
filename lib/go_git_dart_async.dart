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

int _nextCloneRequestId = 0;
final _cloneRequests = <int, Completer<Exception?>>{};

int _nextFetchRequestId = 0;
final _fetchRequests = <int, Completer<Exception?>>{};

int _nextPushRequestId = 0;
final _pushRequests = <int, Completer<Exception?>>{};

int _nextDefaultBranchRequestsId = 0;
final _defaultBranchRequests = <int, Completer<(String?, Exception?)>>{};

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
        if (data is _PushRequest) {
          try {
            var repo = GitBindings(data.libPath);
            repo.fetch(
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
        throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
      });

    // Send the port to the main isolate on which we can receive requests.
    sendPort.send(helperReceivePort.sendPort);
  }, receivePort.sendPort);

  return completer.future;
}();
