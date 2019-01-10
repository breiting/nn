# nn - new note

**Note taking** is an activity I need every day for work. In the last years, I found myself always being disappointed
by existing note taking software tools. Starting with the most popular but bloated tools like Evernote or OneNote, I
also spent some time in evaluating other open source tools such as [Joplin](https://joplin.cozic.net/).

For me, as a heavy [vim](https://vim.org) user, I want to write my notes in vim, being most efficient in typing.
Additionally, for easy formatting, I am a fan of [Markdown](https://daringfireball.net/projects/markdown/) which allows
me to concentrate on content rather than on formatting.

## Why yet another note taking tool?

Well, I want to keep things simple. I already started out to create a certain directory structure, each of which storing
my notes as markdown (similar to [vimwiki](https://github.com/vimwiki/vimwiki)).

`nn` is a tool which helps me to work with my structure, simply create new notes based on templates, and also find notes
easily using existing tools such as [ack](https://beyondgrep.com/). In its core, all notes are written with my favorite
editor (vim), but can easily be switched to your favorite editor by setting the `EDITOR` environment variable.

## Roadmap

`nn` is at its beginning, and will grow over time. However, as being a friend of the [suckless community](https://suckless.org),
`nn` should stay lean and fast.

### Planned Features

All planned features can be found in [issues](https://github.com/breiting/nn/issues). Feel free to add new ideas, or to
help developing this tool.

## Why Go?

In my software development career, I passed by almost all bigger programming languages like C/C++, Java, C#. Recently, i
came across [Go](https://golang.org) and somehow like the language being precise and efficient.
