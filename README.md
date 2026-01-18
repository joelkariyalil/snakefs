# snakefs
Trying to build a snake game to play on the terminal using Go



## Architectural Descisions

- we can keep track of the snake's whole body relative to each other, instead of having absolute values.

|   |    |  x (c)  |  x (b)   |   x (a)   |

- the head (a) would be a absolute value, however, (b) could be denoted as [-1,0]


also worthwhile to note that 2D is not required to denote this, but however, if I were to introduce a diagonal movement, then it would be something for me to add.

So, the above intuition that I had about the solution wasn't great. Basically I simply had to think of it as growing the snake from it's head, and removing the trailing tails if they weren't
supposed to be there, basically new heads everytime, and only remove the tails as the snake progresses.

That was actually a really interesting way of doing it.

The below method would be illustrated here below

Before : |   | x | x |   |

After :  |   |   | x | x |


Basically, adds a new head, and deletes the tail, unless food is concerned here.
