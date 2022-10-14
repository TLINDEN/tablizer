## Project Goals

The goal of this  project is to build a small tool  which helps in day
to day  work with tabular  output of various commandline  programs. It
should be small,  fast and easy to understand. The  idea is to replace
multiline  shell  pipes  using  awk,   sed  and  grep  with  just  one
binary.

There will be no GUI, no web interface, no public API of some sort, no
builtin interpreter.

The  programming  language  used  for  this  project  will  always  be
[GOLANG](https://go.dev/)  with  the  exception of  the  documentation
([Perl POD](https://perldoc.perl.org/perlpod)) and the Makefile.

# Contributing

You can contribute to this project in various ways:

## Open an issue

If you encounter  a problem or don't understand how  the program works
or if you think the documentation is unclear, please don't hesitate to
open an issue.

Please add as much information about the case as possible, such as:

- Your environment (operating system etc)
- tablizer version (`tablizer --version`)
- Input data. Please replace sensitive information with mock data!
- Actual program output.
- Expected program output.
- Error message - if any.

Be aware  that I  am working on  this (and some  other) project  in my
spare  time which  is scarce.   Therefore  please don't  expect me  to
respond to  your query within  hours or even  days. Be patient,  but I
WILL respond.

## Pull Requests

Code and documentation help is  always much appreciated! Please follow
thes guidelines to successfully contribute:

-  Every  pull   request  shall  be  based   on  latest  `development`
  branch. `main` is only used for releases.
  
- Execute the  unit tests before committing: `make  test`. There shall
  be no errors.
  
- Strive  to be  backwards compatible  so that  users who  are already
  using the program  don't have to change their habits  - unless it is
  really neccessary.

- Try to add a unit test for your addition.

- Don't ever change existing unit tests!

- Add a meaningful and comprehensive rationale about your contribution:
  - Why do you think it might be useful for others?
  - What did you actually change or add?
  - Is there an open issue which  this PR fixes and if so, please link
    to that issue.

- [Re-]format your code with `gofmt -s`.

- Avoid unneccesary dependencies, especially for very small functions.

- **If** a  new dependency is being added, it  must be compatible with
  our [license agreement](LICENSE).
  
- You  need to accept  that the  code or documentation  you contribute
  will be redistributed under the  terms of said license agreement. If
  your  contribution  is  considerably  large  or  if  you  contribute
  regularly, then  feel free  to add  your name and  if you  want your
  email    address     to    the    *AUTHORS*    section     of    the
  [manpage](tablizer.pod).

- Adhere to the above mentioned project goals.

-  If  you are  unsure if  your addition or  change will  be accepted,
  better ask before starting coding. Open an issue about your proposal
  and let's  discuss it! That way  we avoid doing unnessesary  work on
  both sides.
    
Each pull  request will be  carefully reviewed and  if it is  a useful
addition  it  will  be  accepted. However,  please  be  prepared  that
sometimes a  PR will be  rejected.  The reasons  may vary and  will be
documented.   Perhaps the  above guidelines  are not  matched, or  the
addition seems  to be not so  useful from my perspective,  maybe there
are  too  much  changes  or  there  might  be  changes  I  don't  even
understand.

But whatever happens: your contribution is always welcome!
