+++
title = "Markdown Formatting Test"
tags = [ "markdown", "md", "formatting", "html", "golang", "website", "static" ]
desc = "This file shows the markdown syntax options available in Burlough."
+++

Burlough also supports Github flavoured markdown.

# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6


## Horizontal Rule

--------------------------------------------------------------------------------

## Inline Stuff

Paragraph

Another Paragraph

**Bold Text**

*Italic Text*

~~Strikethrough Text~~

`Inline Code`

[Link](https://github.com/aghorui/burlough)


## Block Level Stuff

```
Code Block
```

```python
# Code Block with Highlighting

def fib(n):
	if n < 1:
		return 0
	elif n == 1:
		return 1
	else:
		return fib(n - 1) + fib(n - 2)
```

>
> Blockquote
>
> >
> > Nested Blockquote
> >
>

## List

### Unordered

* This is a point.
  * This is a nested point.
* This is another point.
* This is yet another point.

### Ordered

1. This is a numbered list.
2. This is another point.
   * A subpoint.
3. This is yet another point.
   1. A numbered subpoint
   2. Another subpoint

### Setting Offset

34. Numberings can start from any number.
1. Even Here.
1. Here too.


## Image:

![Image](https://github.com/aghorui/burlough/raw/master/doc/logo.svg)


## HTML Tags

<details>
<summary>
Expand Me
</summary>
<p> Some more text </p>
</details>

## Tables

|Thing   |Value  |
|--------|-------|
|Apple   |$1     |
|Ball    |$2     |
|Camera  |$300   |
|Dog     |$40000 |
