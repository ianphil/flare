# Introduction

First off, thank you for considering contributing to Flare. It's people like you that make it such a great tool.

Following these guidelines helps to communicate that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalize your pull requests.

Flare is an open source project and we love to receive contributions from our community — you! There are many ways to contribute, from writing tutorials or blog posts, improving the documentation, submitting bug reports and feature requests or writing code which can be incorporated into Flare itself.

Please feel free to use the issue tracker for support questions.

# Ground Rules

Responsibilities
 * Ensure cross-platform compatibility for every change that's accepted. Windows, Mac, Debian & Fedora Linux.
 * Ensure that code that goes into core meets all requirements in this checklist: https://gist.github.com/audreyr/4feef90445b9680475f2
 * Create issues for any major changes and enhancements that you wish to make. Discuss things transparently and get community feedback.
 * Don't add any classes to the codebase unless absolutely needed. Err on the side of using functions.
 * Keep feature versions as small as possible, preferably one new feature per version.
 * Be welcoming to newcomers and encourage diverse new contributors from all backgrounds.

# Your First Contribution

Unsure where to begin contributing to Atom? You can start by looking through these beginner and help-wanted issues:
- Beginner issues - issues which should only require a few lines of code, and a test or two.
- Help wanted issues - issues which should be a bit more involved than beginner issues.
- Both issue lists are sorted by total number of comments. While not perfect, number of comments is a reasonable proxy for impact a given change will have.

Here are a couple of friendly tutorials you can include: http://makeapullrequest.com/ and http://www.firsttimersonly.com/

Working on your first Pull Request? You can learn how from this *free* series, [How to Contribute to an Open Source Project on GitHub](https://egghead.io/series/how-to-contribute-to-an-open-source-project-on-github).

At this point, you're ready to make your changes! Feel free to ask for help; everyone is a beginner at first :smile_cat:

If a maintainer asks you to "rebase" your PR, they're saying that a lot of code has changed, and that you need to update your branch so it's easier to merge.

# Getting started

For something that is bigger than a one or two line fix:

>1. Create your own fork of the code
>2. Do the changes in your fork
>3. If you like the change and think the project could use it:
    * Be sure you have followed the code style for the project.

# How to report a bug

If you find a security vulnerability, do NOT open an issue. Email ian.philpot@gmail.com instead.

In order to determine whether you are dealing with a security issue, ask yourself these two questions:
 * Can I access something that's not mine, or something I shouldn't have access to?
 * Can I disable something for other people?

If the answer to either of those two questions are "yes", then you're probably dealing with a security issue. Note that even if you answer "no" to both questions, you may still be dealing with a security issue, so if you're unsure, just email me.

When filing an issue, make sure to answer these five questions:

 1. What version of Go are you using (go version)?
 2. What operating system and processor architecture are you using?
 3. What did you do?
 4. What did you expect to see?
 5. What did you see instead?

General questions should go to the golang-nuts mailing list instead of the issue tracker. The gophers there will answer or ask you to file an issue if you've tripped over a bug.

# How to suggest a feature or enhancement

The Flare philosophy is to provide small, robust tooling for testing HTTP servers.

If you find yourself wishing for a feature that doesn't exist in Flare, you are probably not alone. There are bound to be others out there with similar needs. Many of the features that Flare has today have been added because our users saw the need. Open an issue on our issues list on GitHub which describes the feature you would like to see, why you need it, and how it should work.

# Code review process

The core team looks at Pull Requests on a regular weekly basis. All colaboration will take part on the PR we are reviewing.

After feedback has been given we expect responses within two weeks. After two weeks we may close the pull request if it isn't showing any activity.

Borrowed from: [CONTRIBUTING-template.md](https://github.com/nayafia/contributing-template/blob/master/CONTRIBUTING-template.md)
