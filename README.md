# Yet Another Fastbin Client

```
	fastbin
		Yet Another CLI Client for https://github.com/jozsefsallai/fastbin-rewritten
		Usage:  ./fastbin  [Args] [file]
		Or You can also use pipe to send data to fastbin

           You should specify the server address with the environment variable FASTBIN_ADDR, default is http://localhost:3000

		e.g. echo "Hello World" |  ./fastbin
		      ./fastbin  -f file.txt
		      ./fastbin  -f file.txt
		      ./fastbin  -rm abcde
		      ./fastbin  -c abcde

	Args:
           -h, --help: Show this help message
           -d, --debug: Set log level as debug
           -f [file], Upload from a file.
           -rm [secret]: Delete the file with the given secret
           -c, Clone/Create a snippet with Text Editor, and then upload it to fastbin (vim/nano/...)
               If left empty, will create a new empty snippet
	Repo Addr: https://github.com/Linserin/yet-another-fastbin-cli

```
## Feature
- Use [vipe](https://joeyh.name/code/moreutils/) to clone a snippet from fastbin, edit it, and then upload it back to fastbin.

Please install via:
```shell
sudo apt install moreutils
```