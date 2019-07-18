# Tasks to be performed when doing a release

These are the tasks that need to be performed for every release.  In the
instructions below, everywhere where we write `<x.y.z>`, you should replace
that string with the actual version number, e.g. `0.6.1`.

## Documentation

 * Possibly update the installation instructions in the `README.md` file.
 * Create a new release announcement in the `releases/` subdirectory, named
   RELEASE-<x.y.z>.md.  Use a previous one as template.

## Book keeping in the source code

 * Rebuild CX and run all the tests, not just the default ones.  If they don't
   pass as expected, stop here and fix the issues.
 * Change the version number in `cxgo/cxgo.nex` to the current version. This
   is stored in the `VERSION` constant.
 * Rebuild CX and rerun the standard tests.

At this point, everything should be working and the code is ready to actually
be released.

 * Change the CHANGELOG.md file:
   - remove the `(NOT YET RELEASED)` note from the current version, and create
     a new version with the same note.

 * Merge the `develop` branch in git to the `master` branch.  This is done by
   using these commands:
   ```
   git checkout master
   git merge develop
   ```

 * Tag the release in git using the following command:
   `git tag v<x.y.z>`. Note that `x`, `y` and `z` must be replaced, but the `v`
   should be an actual 'v'.  This command should also be performed in the
   `master` branch.

## Uploading the binary packages

 * Create a new release at [the Github release
   page](http://github.com/skycoin/cx/releases). Do this by pressing the
   "Draft a new release" button.
 * Build the binaries for Linux, Mac and Windows. Upload these to the release draft page.

The binary releases are generated for only three platforms at the moment: Linux, MacOS and Windows, all of them in their 64-bit variations. The binaries need to be created using either virtual machines or host computers with the targeted operating system installed on it, as Golang's crosscompiler currently does not work for generating CX binaries for other platforms (there is a conflict with CX's OpenGL libraries) .

The current CX binaries are being generated in a setup that consists of a Windows 10 host machine and two virtual machines running on VirtualBox: one with Debian 9.8 installed on it and another one with MacOS 10.12.

The steps followed in this setup are as follow.

#### Linux Virtual Machine:

* Navigate to CX's git repository:
```bash
cd $GOPATH/src/github.com/skycoin/cx
```
* Fetch the origin remote, set the git repo to be on the `develop` branch and reset the branch to point to `origin/develop`.
```bash
git fetch && git checkout develop && git reset --hard origin/develop
```
* Build the new version of CX:
```bash
make build-full
```
* Check that the version of the built CX binary corresponds to the one that the release is being built for:
```
cx --version
```
* The resulting binary, which is stored in `$GOPATH/bin` needs to be archived in a zip file. In the following command, replace variable `$VERSION` with an appropriate value.
```bash
zip $GOPATH/bin/cx-$VERSION-bin-linux-x64.zip .
```
* Download the zip file to your host machine. The following command is an example, please change according to your setup. Also, replace variable `$VERSION` with an appropriate value.
```bash
scp -P 2022 amherag@127.0.1.1:/home/amherag/go/bin/cx-$VERSION-bin-linux-x64.zip .
```

#### MacOS Virtual Machine

The steps are similar to the steps for the Linux virtual machine. The exceptions are in the commands where the name of the zip file needs to be set for `macos` instead of `linux` and in the command for downloading the zip file. The steps are presented below anyway, for clarification purposes:

* Navigate to CX's git repository:
```bash
cd $GOPATH/src/github.com/skycoin/cx
```
* Fetch the origin remote, set the git repo to be on the `develop` branch and reset the branch to point to `origin/develop`.
```bash
git fetch && git checkout develop && git reset --hard origin/develop
```
* Build the new version of CX:
```bash
make build-full
```
* Check that the version of the built CX binary corresponds to the one that the release is being built for:
```
cx --version
```
* The resulting binary, which is stored in `$GOPATH/bin` needs to be archived in a zip file. In the following command, replace variable `$VERSION` with an appropriate value.
```bash
zip $GOPATH/bin/cx-$VERSION-bin-macos-x64.zip .
```
* Download the zip file to your host machine. The following command is an example, please change according to your setup. Also, replace variable `$VERSION` with an appropriate value.
```bash
scp -P 3022 amherag@127.0.1.1:/Users/amherag/go/bin/cx-$VERSION-bin-macos-x64.zip .
```

#### Windows Host Machine

* Navigate to CX's git repository:
```bash
cd %GOPATH%\src\github.com\skycoin\cx
```
* Fetch the origin remote, set the git repo to be on the `develop` branch and reset the branch to point to `origin/develop`.
```bash
git fetch && git checkout develop && git reset --hard origin/develop
```
* Build the new version of CX. The Makefile does not work for Windows systems under common circumstances, and the Windows build script is often broken, so we need to run the commands that build CX ony by one:
```bash
%GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.nex
%GOPATH%\bin\goyacc -v '' -o %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.go %GOPATH%\src\github.com\skycoin\cx\cxgo\cxgo0\cxgo0.y
%GOPATH%\bin\nex -e %GOPATH%\src\github.com\skycoin\cx\cxgo\parser\cxgo.nex
%GOPATH%\bin\goyacc -v '' -o %GOPATH%\src\github.com\skycoin\cx\cxgo\parser\cxgo.go %GOPATH%\src\github.com\skycoin\cx\cxgo\parser\cxgo.y
go build -tags full -i -o %GOPATH%\bin\cx.exe github.com\skycoin\cx\cxgo\
```
* Check that the version of the built CX binary corresponds to the one that the release is being built for:
```
%GOPATH%\bin\cx --version
```
* The resulting binary, which is stored in `$GOPATH/bin` needs to be archived in a zip file. In the following command, replace variable `%VERSION%` with an appropriate value. In this particular case, we're using 7zip to create the zip archive:
```bash
"C:\Program Files\7-Zip\7z.exe" a -tzip %GOPATH%\bin\cx-%VERSION%-bin-windows-x64.zip %GOPATH%\bin\cx.exe
```

## Upload the online documentation

 * Take the file `documentation/BLOCKCHAIN.md`which contains the user-targeted
   documentation on how to run CX programs on the blockchain.  This should be
   uploaded to https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

   NOTE: Before uploading, remove the comment in the beginning about copying to the wiki.

