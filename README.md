### Input

![input.png](input.png)

### Output (Expand to see)

<details>

<summary>(before you click, I have to warn you against the flashing lights)</summary>

![output.gif](output.gif)

</details>


### Usage

```bash
go install github.com/egeozcan/gifDisco@latest
# within the directory that contains the input.png file
gifDisco
```

No, you can't specify the input file. 
It's always `input.png` and the output is always `{timestamp}_disco.gif`. I'm lazy.

You can also pipe a file to the program, like so:

```bash
cat input.png | gifDisco
```

No, the piped file doesn't have to be named `input.png`.

### TODO

- [X] Nothing. It's perfect.

### Adopt a Panda

If you like this project (for some inexplicable reason),
you can adopt a panda from [WWF](https://gifts.worldwildlife.org/gift-center/gifts/Species-Adoptions/Panda),
or you can become a panda, which has the added benefit of being able to adopt yourself.
