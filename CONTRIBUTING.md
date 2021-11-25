# Contributing

## Procedure

Merge requests are welcome. Just follow this basis:

1. Create an issue before getting into code. There we will discuss what you would like to change if need.
2. Configure your git client. Only signed commits will be merged.
3. Follow [Conventional Commits](https://www.conventionalcommits.org/)
4. Attach tests to each change you have performed
5. Run local verifications before push
6. Explain your exact intentions and ensure you will address all your change's motivation criteria
7. Wait for review and feedback

## Code of conduct
This project follows the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

Instances of abusive, harassing, or otherwise unacceptable behavior will be a reason of permanent ban.

## Good Practicing

### KISS, YAGNI, MVP, etc.

Sometimes we need to remind each other of core tenets of software design - Keep It Simple, You Aren't Gonna Need It, Minimum Viable Product, and so on.
Adding a feature "because we might need it later" is antithetical to software that ships.
Add the things you need NOW and (ideally) leave room for things you might need
later - but don't implement them now.

### Smaller Is Better: Small Commits, Small Pull Requests

Small commits and small pull requests get reviewed faster and are more likely to be correct than big ones.

Attention is a scarce resource.
If your pull request takes 60 minutes to review, the reviewer's eye for detail is not as keen in the last 30 minutes as it was in the first.
It might not get reviewed at all if it requires a large continuous block of time from the reviewer.

**Breaking up commits**

Break up your pull request into multiple commits, at logical break points.

Making a series of discrete commits is a powerful way to express the evolution of an idea or the different ideas that make up a single feature.
Strive to group logically distinct ideas into separate commits.

For example, if you found that Feature-X needed some prefactoring to fit in, make a commit that JUST does that prefactoring.
Then make a new commit for Feature-X.

Strike a balance with the number of commits.
A pull request with 25 commits is still very cumbersome to review, so use your best judgment.

**Breaking up Pull Requests**

Or, going back to our prefactoring example, you could also fork a new branch, do the prefactoring there and send a pull request for that.
If you can extract whole ideas from your pull request and send those as pull requests of their own, you can avoid the painful problem of continually rebasing.

Multiple small pull requests are often better than multiple commits.
Don't worry about flooding us with pull requests. We'd rather have 100 small,obvious pull requests than 10 unreviewable monoliths.

We want every pull request to be useful on its own, so use your best judgment on what should be a pull request vs. a commit.

As a rule of thumb, if your pull request is directly related to Feature-X and nothing else, it should probably be part of the Feature-X pull request.
If you can explain why you are doing seemingly no-op work ("it makes the Feature-X change easier, I promise") we'll probably be OK with it.
If you can imagine someone finding value independently of Feature-X, try it as a pull request.
Instead, reference other pull requests via the pull request your commit is in.)

### Open a Different Pull Request for Fixes and Generic Features

**Put changes that are unrelated to your feature into a different pull request.**

Often, as you are implementing Feature-X, you will find bad comments, poorly named functions, bad structure, weak type-safety, etc.

You absolutely should fix those things (or at least file issues, please) - but not in the same pull request as your feature. Otherwise, your diff will have way too many changes, and your reviewer won't see the forest for the trees.

**Look for opportunities to pull out generic features.**

For example, if you find yourself touching a lot of modules, think about the dependencies you are introducing between packages.
Can some of what you're doing be made more generic and moved up and out of the Feature-X package?
Do you need to use a function or type from an otherwise unrelated package?
If so, promote!
We have places for hosting more generic code.

Likewise, if Feature-X is similar in form to Feature-W which was checked in last month, and you're duplicating some tricky stuff from Feature-W, consider prefactoring the core logic out and using it in both Feature-W and
Feature-X.
(Do that in its own commit or pull request, please.)

### Comments Matter

In your code, if someone might not understand why you did something (or you won't remember why later), comment it. Many code-review comments are about this exact issue.

If you think there's something pretty obvious that we could follow up on, add a TODO.

Read up on [GoDoc](https://blog.golang.org/godoc-documenting-go-code) - follow those general rules for comments.

### Test

Nothing is more frustrating than starting a review, only to find that the tests are inadequate or absent.
Very few pull requests can touch the code and NOT touch tests.

If you don't know how to test Feature-X, please ask!
We'll be happy to help you design things for easy testing or to suggest appropriate test cases.

### Additional Resources

- [How to Write a Git Commit Message - Chris Beams](https://chris.beams.io/posts/git-commit/)
- [Distributed Git - Contributing to a Project (Commit Guidelines)](https://git-scm.com/book/en/v2/Distributed-Git-Contributing-to-a-Project)
- [What’s with the 50/72 rule? - Preslav Rachev](https://preslav.me/2015/02/21/what-s-with-the-50-72-rule/)
- [A Note About Git Commit Messages - Tim Pope](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)


## Thanks

Many thanks in advance to everyone who contributes their time and effort to improve this application.