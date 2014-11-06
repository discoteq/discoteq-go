# Contributing to Discoteq

Want to contribute? Up-to-date pointers should be at:
<http://contributing.appspot.com/discoteq>

Got an idea? Something smell wrong? Cause you pain? Or lost seconds of
your life you'll never get back?

All contributions are welcome: ideas, patches, documentation, bug
reports, complaints, and even something you drew up on a napkin.

Programming is not a required skill. Whatever you've seen about open
source and maintainers or community members saying "send patches or
die" - you will not see that here.

It is more important to me that you are able to contribute.
If you haven't got time to do anything else, just email me and I'll try to help: <joseph@josephholsten.com>.

I promise to help guide this project with these principles:

-   Community: If a newbie has a bad time, it's a bug.
-   Software: Make it work, then make it right, then make it fast.
-   Technology: If it doesn't do a thing today, we can make it do it
    tomorrow.

(Some of the above was repurposed with \<3 from logstash)

Here are some ways you can be part of the community:

## Want to Learn?

Want to lurk about and see what others are doing with discoteq?

-    The irc channel (#discoteq on irc.freenode.org) is a good place for this
-    The mailing list is also great for learning from others.

## Something not working? Found a Bug?

Find something that doesn't feel quite right? Here are 5 steps to getting it fixed!

Check your version of discoteq
:    To make sure you're not wasting your time, you should be using the latest version of discoteq before you file your bug. First of all, you should download the latest nightly build to be sure you have the latest version. If you've done this and you still experience the bug, go ahead to the next step.

Search our [issues](https://github.com/discoteq/discoteq-go/issues)
:    Now that you have the latest discoteq version and still think you've found a bug, search through issues first to see if anyone else has already filed it. This step is very important! If you find that someone has filed your bug already, please go to the next step anyway, but instead of filing a new bug, comment on the one you've found. If you can't find your bug in issues, go to the next step.

Create a Github account https://github.com/join
:    You will need to create a Github account to be able to report bugs (and to comment on them). If you have registered, proceed to the next step.

File the bug!
:    Now you are ready to file a bug with discoteq. The [Writing a Good Bug Report](http://www.webkit.org/quality/bugwriting.html) document gives some tips about the most useful information to include in bug reports. The better your bug report, the higher the chance that your bug will be addressed (and possibly fixed) quickly!

What happens next?
:    Once your bug is filed, you will receive email when it is updated at each stage in the bug life cycle. After the bug is considered fixed, you may be asked to download the latest nightly and confirm that the fix works for you.

(This section lovingly adapted from the [Webkit project](http://www.webkit.org/quality/reporting.html))


## Submitting patches

* use a feature branch

* rebase into single patch

* single patch for each feature

* say _why_ the changes were made, we can look at the diff to see _how_ they were made. 

* any new codepaths (not just features) should have new tests

* any new features should have new integration tests


## Setting up a local dev environment


For those of you who do want to contribute with code, we've tried to
make it easy to get started. You can install all dependencies and tools
with:

    script/bootstrap

Then you can start support services (like a local chef server) with:

    forego start

Plenty of example data already exists in `stubs/`, though it probably
deserves more explanation.

Good luck!

## Style guide

At the moment, we're just using the [style standards of the golang team](https://code.google.com/p/go-wiki/wiki/CodeReviewComments). Many of these will automatically be verified by running `make lint`, which runs 

* [flint](https://github.com/pengwynn/flint)
* [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
* [golint](https://github.com/golang/lint)
* [go vet](https://godoc.org/golang.org/x/tools/cmd/vet)
