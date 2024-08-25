#!/bin/bash

set -eux

# Assumes Android Studio is installed on the standard path along with Android SDK/NDK
export ANDROID_HOME=~/Library/Android/sdk/
export ANDROID_NDK_HOME=~/Library/Android/sdk/ndk/$(ls $ANDROID_HOME/ndk | head -1)
export JAVA_HOME=/Applications/Android\ Studio.app/Contents/jre/Contents/Home

ROOTDIR=$(cd $(dirname $0); pwd -P)
DISTDIR=$ROOTDIR/../android/src/main
TMPDIR=$ROOTDIR/tmp/android

export PATH=$ANDROID_HOME/ndk/26.1.10909125/toolchains/llvm/prebuilt/darwin-x86_64/bin/:$ANDROID_HOME/platform-tools:$PATH

rm -rf $DISTDIR/jniLibs $TMPDIR
mkdir -p \
    $DISTDIR \
    $TMPDIR/armeabi-v7a \
    $TMPDIR/arm64-v8a \
    $TMPDIR/x86 \
    $TMPDIR/x86_64

ANDROID_SDK_VERSION=34
CLANG_SUFFIX=-linux-android$ANDROID_SDK_VERSION-clang

export GOOS=android
export CGO_ENABLED=1

GOARCH=arm CC=armv7a-linux-androideabi$ANDROID_SDK_VERSION-clang CXX=armv7a-linux-androideabi$ANDROID_SDK_VERSION-clang++ \
    go build -v -x -buildmode=c-shared -trimpath -o=$TMPDIR/armeabi-v7a/go_git_dart.so
GOARCH=arm64 CC=aarch64$CLANG_SUFFIX CXX=aarch64$CLANG_SUFFIX++ \
    go build -v -x -buildmode=c-shared -trimpath -o=$TMPDIR/arm64-v8a/go_git_dart.so
GOARCH=386 CC=i686$CLANG_SUFFIX CXX=i686$CLANG_SUFFIX++ \
    go build -v -x -buildmode=c-shared -trimpath -o=$TMPDIR/x86/go_git_dart.so
GOARCH=amd64 CC=x86_64$CLANG_SUFFIX CXX=x86_64$CLANG_SUFFIX++ \
    go build -v -x -buildmode=c-shared -trimpath -o=$TMPDIR/x86_64/go_git_dart.so

cp -rf $TMPDIR $DISTDIR/jniLibs
find $DISTDIR -name "*.h" -exec rm {} \;
rm -rf $DISTDIR/jniLibs/tmp

echo "Android build successful"
echo "Copying shared libraries to android project"

cp ./tmp/android/armeabi-v7a/go_git_dart.so ../android/src/main/jniLibs/armeabi-v7a/go_git_dart.so
cp ./tmp/android/arm64-v8a/go_git_dart.so ../android/src/main/jniLibs/arm64-v8a/go_git_dart.so
cp ./tmp/android/x86/go_git_dart.so ../android/src/main/jniLibs/x86/go_git_dart.so
cp ./tmp/android/x86_64/go_git_dart.so ../android/src/main/jniLibs/x86_64/go_git_dart.so

echo "Done"