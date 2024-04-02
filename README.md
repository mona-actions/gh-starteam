# gh-starteam

`gh-starteam` is a GitHub Action that allows you to migrate a StarTeam repository to GitHub. As we are still developing this tool, current it just _**simulates**_ the migration progress by creating text files in place of the StarTeam files. Eventually this tool will make use of the StarTeam CLI tool to fetch the appropriate files and commit them to the GitHub repository.

## Install

```bash
gh extension install mona-actions/gh-starteam
```

## Usage

Ths tool takes as input the output of the `st list` command. The command should be run in te project you wish to convert and saved to a file.

```bash
Usage:
  starteam [flags]

Flags:
  -h, --help                  help for starteam
  -f, --history-file string   Path to the history file
  -r, --repo-path string      Path to the repository
```

### Example

```bash
gh starteam -f history.txt -r /path/to/repo
```

## StarTeam to Git Mapping

The following table shows how StarTeam concepts are mapped to Git concepts. The StarTeam concepts are pulled from the output of the `st list` command.

| StarTeam | Git         |
| -------- | ----------- |
| Folder   | Directory   |
| File     | File        |
| Revision | Commit      |
| Date     | Commit Date |
| Author   | Committer   |
| Comment  | Commit Msg  |

### Example Parsing of `st list` Output

Given the following output from the `st list` command:

```
StarTeam 10.4 Command Line Interface, Build 10.4.8.36
Copyright (c) 2003-2008 Borland Software Corporation. All rights reserved.
Using option file: /home/cmbld/.starteam-client/platform.ini
Folder: src  (working dir: C:/WorkingFolder/MY_PROJECT)
Folder: lib  (working dir: C:/WorkingFolder/MY_PROJECT/src/lib)
History for: example.java
Description: Example Library
Locked by:
Status: Missing
----------------------------
Revision: 1 View: src Branch Revision: 1.0
Author: Ron Bugundy Date: 11/19/09, 4:26:51 PM EST
Adding example.java

=============================================================================

History for: more.java
Description: More Example Library
Locked by:
Status: Missing
----------------------------
Revision: 4 View: src Branch Revision: 1.3
Author: Brick Tamland Date: 12/28/09, 2:56:12 PM EST
Issue #566: Fixed bug in more.java

----------------------------
Revision: 3 View: src Branch Revision: 1.2
Author: Brick Tamland Date: 12/13/09, 10:35:48 AM EST
Made some changes to more.java
Forgot to add a line

----------------------------
Revision: 2 View: src Branch Revision: 1.1
Author: Brick Tamland Date: 11/30/09, 12:17:26 PM EST
Added some functionality based on feedback

----------------------------
Revision: 1 View: src Branch Revision: 1.0
Author: Brick Tamland Date: 11/23/09, 6:29:02 PM EST
Ron yelled at me to add more.java

=============================================================================

Folder: src  (working dir: C:/WorkingFolder/MY_PROJECT/src)
History for: main.java
Description: Java entry point
Locked by:
Status: Missing
----------------------------
Revision: 2 View: src Branch Revision: 1.0.3.0
Author: Veronica Corningstone Date: 12/13/09, 10:35:48 AM EST
Another amazing update

----------------------------
Revision: 1 View: src Branch Revision: 1.0
Author: Veronica Corningstone Date: 11/30/09, 12:17:26 PM EST
This project is a mess, no entry point

=============================================================================

Folder: sports  (working dir: C:/WorkingFolder/MY_PROJECT/src/sports)
History for: whammy.java
Description: Whammy class
Locked by:
Status: Missing
----------------------------
Revision: 2 View: src Branch Revision: 1.1
Author: Champ Kind Date: 12/28/09, 2:56:12 PM EST
WHAMMY!

----------------------------
Revision: 1 View: src Branch Revision: 1.0
Author: Champ Kind Date: 12/13/09, 10:35:48 AM EST
WHAMMY!
=============================================================================

Folder: stuff  (working dir: C:/WorkingFolder/MY_PROJECT/stuff)
Folder: other_stuff  (working dir: C:/WorkingFolder/MY_PROJECT/stuff/other_stuff)
Folder: some_other_stuff  (working dir: C:/WorkingFolder/MY_PROJECT/stuff/some_other_stuff)
Folder: look_at_all_the_stuff  (working dir: C:/WorkingFolder/MY_PROJECT/stuff/look_at_all_the_stuff)

```

In order to parse the output the tool takes the following steps:

1. Remove the first 3 lines as they are not needed as it only contains the StarTeam command line version and the option file.
2. Split the output by the `=============================================================================` separator.
3. The last section is removed as it only contains folder information.
4. The remaining sections are split by the `----------------------------` separator and the data is flattened into a list of revisions.
5. Each revision is grouped by date and combined into a single commit. The commit message is the concatenation of the comments.

The final output in the resulting repository would look like this (not a single branch is currently only supported):

```
commit asldfhoiherhdv8eh52435gg (HEAD -> master)
Author: Ron Bugundy <test@test.com>
Date: 11/19/09, 4:26:51 PM EST

    Number: 1
    Author: Ron Bugundy
    Date: 11/19/09, 4:26:51 PM EST
    File: example.java
    Folder: MY_PROJECT/src/lib
    Memo:
    Adding example.java

    ----------------------------

commit 9s8dfh9s8dfh9s8dfh9s8dfh
Author: Brick Tamland <test@test.com>
Date: 11/23/09, 6:29:02 PM EST

    Number: 1
    Author: Brick Tamland
    Date: 11/23/09, 6:29:02 PM EST
    File: more.java
    Folder: MY_PROJECT/src/lib
    Memo:
    Ron yelled at me to add more.java

    ----------------------------

commit arewio092r70gj0gj0gj0g3
Author: Brick Tamland <test@test.com>
Date: 11/30/09, 12:17:26 PM EST

    Number: 2
    Author: Brick Tamland
    Date: 11/30/09, 12:17:26 PM EST
    File: more.java
    Folder: MY_PROJECT/src/lib
    Memo:
    Added some functionality based on feedback

    ----------------------------

    Number: 1
    Author: Veronica Corningstone
    Date: 11/30/09, 12:17:26 PM EST
    File: main.java
    Folder: MY_PROJECT/src
    Memo:
    This project is a mess, no entry point

    ----------------------------

commit u745oi458g2985445o2348g
Author: Brick Tamland <test@test.com>
Date: 12/13/09, 10:35:48 AM EST

    Number: 3
    Author: Brick Tamland
    Date: 12/13/09, 10:35:48 AM EST
    File: more.java
    Folder: MY_PROJECT/src/lib
    Memo:
    Made some changes to more.java
    Forgot to add a line

    ----------------------------

    Number: 2
    Author: Veronica Corningstone
    Date: 12/13/09, 10:35:48 AM EST
    File: main.java
    Folder: MY_PROJECT/src
    Memo:
    Another amazing update

    ----------------------------

    Number: 1
    Author: Champ Kind
    Date: 12/13/09, 10:35:48 AM EST
    File: whammy.java
    Folder: MY_PROJECT/src/sports
    Memo:
    WHAMMY!

    ----------------------------

commit 9s8df489g87gfgd987ag98g
Author: Brick Tamland <test@test.com>
Date: 12/28/09, 2:56:12 PM EST

    Number: 4
    Author: Brick Tamland
    Date: 12/28/09, 2:56:12 PM EST
    File: more.java
    Folder: MY_PROJECT/src/lib
    Memo:
    Issue #566: Fixed bug in more.java

    ----------------------------

    Number: 2
    Author: Champ Kind
    Date: 12/28/09, 2:56:12 PM EST
    File: whammy.java
    Folder: MY_PROJECT/src/sports
    Memo:
    WHAMMY!

    ----------------------------
```

## Roadmap

There are number of questions that need to be answered before this tool can be considered complete:

- How to handle branches?
  - How do we identify them in the history output?
- How to handle labels?
  - How do we identify them in the history output?
  - Should they map to tags/releases in Git?
- How to handle email addresses (perhaps a mapping file)?
- Does the StarTeam command line tool provide the best way to get the history?
  - Can a direct database connection be used instead?

## License

- [MIT](./license) (c) [Mona-Actions](https://github.com/mona-actions)
- [Contributing](./.github/contributing.md)
