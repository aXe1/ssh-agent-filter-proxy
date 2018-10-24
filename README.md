# About

Where is authentication problem, when using git over ssh with popular services like Github and Gilab. Both of them require you to connect with user `git`, and authentication is done via public key. When you have more when one account at such service, and use seperate keys for them, then you will authenticated with the first key you supply, as all keys are valid. This tiny program is designed to filter out access to unnecessary keys in agent. It is inspired with [ssh-agent-filter](https://git.tiwe.de/ssh-agent-filter.git), but written in Go and works under Windows.

# Status

Currently, it's [MVP](https://en.wikipedia.org/wiki/Minimum_viable_product) only. Although, it's usable right now, there are too many required manual configuration steps, and it's tested in only one documented setup: [KeePass 2.x](https://keepass.info), [KeeAgent](https://lechnology.com/software/keeagent/) and [Git for Windows (ex. msysGit)](https://git-scm.com/download/win) under Windows 10.

# Usage

1. Follow [KeeAgent installation instructions](https://keeagent.readthedocs.io/en/stable/installation.html#windows) to setup it, and also [enable `Create msysGit compatible socket file` option](https://keeagent.readthedocs.io/en/stable/usage/tips-and-tricks.html#cygwin-and-msys).

2. Manualy create somewhere text file (it will be unix socket for any Cygwin-compatible program) with the following content, and set `system` attribute to it:

    ```
    !<socket ><tcp_port_num> s <random_guid>
    ```

    For example (`~/.ssh/keeagent-restricted.sock`):
  
    ```
    !<socket >52101 s A5520E1E-4D0DDFEF-C8F1089C-34EB0CB3
    ```

3. Then run `ssh-agent-filter-proxy`:

    ```shell
    SSH_AUTH_SOCK=<path_to_keeagent_msysgit_compatible_socket> go run ssh-agent-filter-proxy.go <tcp_port_num> <permitted_key_comment>
    ```
    
    Example:
    
    ```shell
    SSH_AUTH_SOCK=~/.ssh/keeagent-restricted.sock go run ssh-agent-filter-proxy.go 52101 some-user@example.com
    ```

4. Now you can use any program that uses authentication against openssh agent as follows, and it will use only the key with comment you supplied:

    ```shell
    SSH_AUTH_SOCK=<path_to_created_socket_file> ssh -T git@github.com
    ```
    
    Example:
    
    ```shell
    SSH_AUTH_SOCK=~/.ssh/keeagent-restricted.sock ssh -T git@github.com
    ```