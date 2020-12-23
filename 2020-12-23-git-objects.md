# Git Objects

Git Objects are Blob Object, Tree Object, Commit Object.

All objects are stores at `.git/objects` with unique key (hash generated).

**Blob Object** is file content.

**Tree Object** is filename or directory.

![Simple version of the Git data model](https://git-scm.com/book/en/v2/images/data-model-1.png)

**Commit Object** contains top-level tree (snapshot at the time), parent commits, author and commit message.

Branch is pointer, which point to commit we want.

## References

- [Git Internals - Git Objects](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)

- [Commits are snapshots, not diffs](https://github.blog/2020-12-17-commits-are-snapshots-not-diffs/)
