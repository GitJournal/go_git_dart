// ignore_for_file: avoid_print

import 'dart:io';
import 'package:go_git_dart/go_git_dart.dart';

String _getCorrectLibrary() {
  return '/Users/vishesh/src/gitjournal/go-git-dart/src/gitjournal.so';
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
          var branch =
              bindings.defaultBranch(arguments[1], pemBytes, arguments[3]);
          print("DefaultBranch: $branch");
          break;
      }
      break;
    case 'keygen':
      var (publicKey, privateKey) = bindings.generateRsaKeys();
      print('Public Key: $publicKey');
      print('Private Key: $privateKey');

      break;
    default:
      print('Unknown command: $command');
      break;
  }
}
