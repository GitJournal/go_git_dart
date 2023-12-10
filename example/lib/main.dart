import 'dart:convert';

import 'package:flutter/material.dart';

import 'package:go_git_dart/go_git_dart.dart' as go_git_dart;
import 'package:path_provider/path_provider.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  final cloneUrlController = TextEditingController();
  final privateKeyController = TextEditingController();
  final privateKeyPasswordController = TextEditingController();

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    const spacerSmall = SizedBox(height: 10);
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Native Packages'),
        ),
        body: SingleChildScrollView(
          child: Container(
            padding: const EdgeInsets.all(10),
            child: Column(
              children: [
                TextField(
                  controller: cloneUrlController,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Clone URL',
                  ),
                ),
                spacerSmall,
                TextField(
                  controller: privateKeyController,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Private Key',
                  ),
                  maxLines: null,
                ),
                spacerSmall,
                TextField(
                  controller: privateKeyPasswordController,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Private Key Password',
                  ),
                ),
                spacerSmall,
                ElevatedButton(
                  onPressed: () async {
                    final cloneUrl = cloneUrlController.text;
                    final privateKey = privateKeyController.text;
                    final privateKeyPassword =
                        privateKeyPasswordController.text;

                    var gitBindings = go_git_dart.GitBindings();
                    final appDocumentsDir =
                        await getApplicationDocumentsDirectory();
                    gitBindings.clone(
                      cloneUrl,
                      appDocumentsDir.path + "/test",
                      utf8.encode(privateKey),
                      privateKeyPassword,
                    );
                  },
                  child: const Text('Clone'),
                ),
                spacerSmall,
              ],
            ),
          ),
        ),
      ),
    );
  }
}
