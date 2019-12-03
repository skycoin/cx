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
[https://github.com/SkycoinProject/cx/issues] and register the bug.  But please try
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

### Documentation for CX Programmers

The main documententation for CX programmers is the CX textbook. The textbook
documents and teaches the language for people who already know programming
from other languages, but it may not be sufficient for people who are new to
programming.

Something that would be very valuable is a book about CX for beginner
programmers. This book would teach the basics of programming and CX at the
same time, giving lots of simple examples and teaching how to think as a
programmer.  There are a number of books about the Python language that can be
used as starting points on how this could be done.

There are many other types of documentation work that would both help CX 
programmers over the world:

* Translations of CX any book to other languages (Note, though, that the CX
  textbook is not yet finished, so this may be better done a bit later.)
* Articles, online and in magazines, about how to use CX in specific
  situations.  Especially valuable would be articles on how to use CX to
  create smart contracts for the blockchain.
* Create and maintain a user-oriented wiki about CX and Skycoin.
* Writing example programs that show important concepts and how to use them.
* etc

### Documentation for CX Developers

The codebase of CX needs more and more structured documentation than is
available today. This documentation should include important data structures,
overall data flow, control flow, etc.  The purpose should be to introduce the
codebase to new programmers rather than give details about API's to seasoned
CX developers.

The other type of documentation that is needed is detailed documentation of
functions and API's that is used by developers every day.  This type of
documentation is created with tools like doxygen, which many developers know,
and GoDoc (https://godoc.org/).

## Improving the Code / Programming

If you want to work on the CX codebase, there are a few ground rules that you
need to know about.  CX use fairly normal practices. Here is a summary of
them.

The sections below assumes that you are familiar with version control using
Git.  If you are not, then you can get a very good introduction in the free
book "Pro Git" ([https://git-scm.com/book/en/v2]). Chapters 1-3 is enough if
you only want to use it as a user. There are both PDF and EBook versions to
download.

You also need some familiarity with the Go language.  A good starting point is
the "Tour of Go" ([https://tour.golang.org/welcome/1]).

### GitHub Workflow

All public Skycoin code is on GitHub, so that's where you will find CX.  If
you are reading this guide, you probably already know this, but here is the
link anyway: [https://github.com/SkycoinProject/cx]

Skycoin uses the standard GitHub workflow where you fork the main repository,
work in your own copy and then create pull requests to get your code into the
main repository. This workflow described in detail on the GitHub help pages
here: [https://guides.github.com/activities/forking].

### Git Branching Model

Skycoin uses a well-tested git branching strategy that is described in detail
here: [https://nvie.com/posts/a-successful-git-branching-model/].

The main idea is that you have one branch called `develop`, where most of the
development is collected. When you start to develop a new feature or fix an
issue, you create a branch for this from the `develop` branch in your forked
repository.  When you are done with the development (or have reached an
important milestone where you want to merge it), you will create a pull
request from your branch and into the `develop` branch of the `skycoin/cx`
repository.

### Code Formatting

Before committing, you should always use the go formatter `gofmt`. The
simplest way to use it is to use the `-w` option like this:

```
gofmt -w path/to/the/file.go
```

That will format the Go source file in place.


### Running the Test Suite

The CX sources contain a sizable and growing test suite, containing regression
tests from the registered issues and a product test suite.  The simplest way
to run the test suite is running the following command in the top directory of
the repository:

```
cx tests/main.cx ++wdir=tests
```

When you do that, you will get a list of tests that fail, looking something
like this:

```
# 88 | FAILED  | 99ms | 'cx issue-141.cx' | expected cx.SUCCESS (0) | got cx.COMPILATION_ERROR (3) | Parser gets confused with `2 -2`
# 89 | FAILED  | 127ms | 'cx issue-153.cx' | expected cx.SUCCESS (0) | got cx.PANIC (2) | Panic in when assigning an empty initializer list to a []i32 variable
```

It should be easy to find an issue that you want to look at. If you think they
are to many or taking to long to run, you can disable some of the tests with
the `++disable-tests` option like this:

```
cx tests/main.cx ++wdir=tests ++disable-tests=issue
```

The available types of tests are `stable`, `issue` and `gui`.

If you are using the tests while debugging, you can add more information by
adding the `++log` option.

```
cx tests/main.cx ++wdir=tests ++log=fail,skip
```

The available types of extra logging are `success`, `stderr`, `fail` and
`skip`.  You can use them individually or combine them by separating them with
a comma.

And finally, you do not always have to run all of the test suite. You can also
run a single test like this:

```
cx tests/issue-141.cx
```

### Internal Documentation

Documentation of the inner workings of CX can be found in the `documentation/`
subdirectory.  At the time of writing this directory, does not contain a lot
of such documentation, but we are working on adding more.

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

# Other Types of Help

You can also help CX by doing other tasks than issue handling or programming.
Here are a few ways that you can help CX and its adoption.

## Marketing

If we want CX and Skycoin to be used everywhere, technical excellence is not
enough.  We also need to market both CX and Skycoin.

In some circuits marketing has a bad name. It is described as lying or
manipulation.  This is not true.  Marketing is simply to let people know about
what you have and let them make their own decisions.  Some marketeers do
indeed lie or use weasel words, but real good marketing does not have to do that.

### Articles

Just like for documentation, articles in various media can be used to spread
the message.  The difference is that when we do marketing, we are more
interested in showing what you can do with CX rather than give intricate
techical tips.  So keep the technical details out of the articles and instead
concentrate on how easy it is to use CX, how cheap and fast Skycoin
transactions are and how many big partners that we have.

For effective marketing, it is a bad idea to lie or exxagerate. CX and Skycoin
are good enough in reality that people will be interested in using them if you
just tell people about them.

## User Support

Do you know CX reasonably well?  In that case you can be a help just by
supporting other users.  No matter how easy a language is, there are always
ways to shoot yourself in the foot. In that case you can provide support to
the poor programmer and help them out of their predicament.

### Forums

The forum that is used by Skycoin is SkyWug (https://skywug.net/forum/). Lots
of new users of various aspects of Skycoin go there to get help.  You can
provide it.

### Chat rooms (Telegram)

Skycoin uses Telegram for chat.  The main channel is
https://t.me/Skycoin. There are many Skycoin channels, but few related to CX.
At the current time there is one CX channel releated to game programming (CX
Skycoin Game Development, or https://t.me/skycoin_game_dev), but we expect
more in the future.

## Training

You could give CX training to friends, students or colleagues.  There is no
formal courseware yet about CX, but we will be developing that.  You could
also contribute by creating such courseware or translating it to other
languages. 
