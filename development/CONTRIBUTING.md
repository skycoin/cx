# How to Contribute to CX Development

The easiest way to contribute to CX Development is to use CX for your own
programming projects and register any issues that you find. But there are many
other ways that you can help if you want to participate actively in the
development process.

Please find below a number of ways that you can help. The development team is
grateful for all the help we can get.

## Issues Handling

Proper handling of issues for CX is a very important aspect of the development
process.  Without it, there can be no quality output.

### Register Issues

If you find a bug in CX, the least you should do is to go to
[https://github.com/skycoin/cx/issues] and register the bug.  But please try
to make sure that you have found a bug that is not already registered by
browsing the headings of the already existing bugs.

If you have found a new bug, please write a descriptive heading and then
provide at least the following information:

* CX version (can be retrieved by the command `cx -v`)
* A more detailed description of the bug and how it appears.
* A minimal test case that triggers the bug.

After registering the issue, please also be prepared to answer follow-up
questions from the developers as they start working on it.  Please note,
though, that there may be many open issues and that it may take some time
before your issue is handled.

## Writing Tests

In the `tests/` subdirectory, you can find a set of test programs that test
various aspects of the CX implementation.  Currently there are two main groups
of tests:

* `issues-*.cx` which are tests which show errors that are registered as
  issues on github.com. 
* `test-*.cx` which are more thorough tests that test specific aspects of CX.

There is also a file `main.cx`, which is the starting point for running all
tests.

The first group is known as the _regression tests_. These tests are failing
for issues that are not yet fixed and they are passing for issues that have
been fixed.  It is important to test all fixes against all tests so that not
old issues spring up again.  Hence the name regression tests.  All new issues
should eventually get an automated test like the current ones.

The second group is the _product test suite_ for the CX language.  The
intention is that all aspects of the CX language should be as thoroughly
tested as is humanly possible. These tests should be broad-reaching and deep
at the same time.  They should at least cover the following aspects:

* the common case
* cases inside all limits (e.g. 1-254 for unsigned byte size values)
* cases on all limits (e.g. 0 and 255 for unsigned byte size values)
* cases just outside all limits (e.g. -1 and 256 for unsigned byte size values)
* far out cases (e.g. -100000 or +10000000 for unsigned byte size values)
* correct and wrong uses
* interesting combinations and type conversions, including wrong uses
* etc

Currently it is the product tests that need enhancements.  Not only are not
all aspects of CX covered by tests, each test set for a particular aspect of
the language are not covering all the cases mentioned above.

Writing a test suite that covers all aspects of something simple, let alone
something as complex as a programming language, is a very big and difficult
task.  The only way that it can be accomplished is by long and diligent work,
so all help here is very important.  Considering that CX will go on the block
chain and be used for creating smart contracts that will potentially handle
millions of dollars makes it even more important.

### Triaging Reported Issues

Sometimes end users register issues which are not actually errors.  Thus, it
is important to check every issue to find out whether it really points out a
real bug or if there are some mistaken assumptions behind the report.  A
registered issue may also be a duplicate of an already existing issue report,
in which case this should be marked as such.

It is inefficient if the developers with the deepest knowledge of the
internals of CX need to wade through many bug reports of which few point out
new errors.  Therefore it is very helpful if some trusted persons can provide
a first-line response to new issue reports and interact with the users that
reported the issues.  Such interaction should include:

* Asking the reporter for more details
* Investigating if the issue is actually a bug
* Investigating if the issue is already known, and mark it as a duplicate if so
* Assess the severity of the issue and alert developers if a critical issue appears
* Follow up with reporters when a bug is fixed by a developer to see if it is
  also fixed for the reporter

For a large project this task can be just as important for development
progress as actually developing.

## Writing or Improving Documentation

Documentation can take many shapes.  Two important target groups for the
documentation are CX users (programmers) and CX developers.

TODO

### Documentation for CX Programmers

The main documententation for CX programmers is the CX textbook.

TODO

### Documentation for CX Developers

TODO

## Marketing

TODO

### Articles

TODO

## User Support

TODO

### Forums

SkyWug

TODO

### Chat rooms (Telegram)

TODO

## Improving the Code 

TODO

### GitHub Workflow

Forking

TODO

### Git Branching Model

See https://nvie.com/posts/a-successful-git-branching-model/

TODO

### Code Formatting

gofmt

TODO

### Running the Test Suite

TODO

### Internal Documentation

Documentation of the inner workings of CX can be found...

TODO

### Register Issues

We already covered registering issues in the GitHub issue list, but as a
developer of the CX code itself you can do more than a normal CX user.  If you
find a bug in CX, please do the following:

1. Register an issue in the Github issue tracker as above. Make note of the
   issue number.
2. Write a new test in the `tests/` subdirectory that fails, showing the
   issue. Name the test file as issue-<xx>.cx, where <xx> is the actual issue
   number from step 1.   Don't forget to add the test to the list in
   `tests/main.cx`, and tag it as `TEST_ISSUE` to group it together with the
   other tests for Github issues.
3. Create a pull request for the new test and new call in `tests/main.cx`.

