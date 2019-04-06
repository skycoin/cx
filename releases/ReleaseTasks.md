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
     a new version with this note.

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

## Uploading

**FIXME**: Amaury, please update this section to include the actual commands that
you perform to build the releases.

 * Create a new release at [the Github release
   page](http://github.com/skycoin/cx/releases). Do this by pressing the
   "Draft a new release" button.
 * Build the binaries for Linux, Mac and Windows. Upload these to the release
   draft page.
 * FIXME: Check and tell how to create the source packages.



