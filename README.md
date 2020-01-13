# Jumbled-word puzzler solver

A program to help solve [Jumbled word puzzles](https://www.google.com/search?q=jumble&tbm=isch)

The user interface is pathetic.

## Building

Has no dependencies outside of the Go standard library.

    $ go build runserver.go

## Executing

    $ ./runserver -d /usr/share/dict/words

You'll have to access it with a URL like `http://localhost:8012/jumble`

## Design and Algorithm

It works by sorting the letters of every word in a file (`/usr/share/dict/words`, for example)
that has one word per line.
The words "meat", "team" and "meta" all end up alphabetized as "aemt".
It makes a hashtable/dictionary/associative array using an alphabetized word
as a key, and an array of words that alphabetize to that key as the value.

The program can alphabetize the characters in a jumbled input word,
and use the resulting string to look up all words that potentiall match.
That lets it "unjumble" words.
Potentially, a jumbled word can match more than one unjumbled alternate word.

The same thing happens in the final solution.
Marked characters of the un-jumbled words get used to create a key
to look up possible matches in the dictionary.
One complication happens when a  letter-alphabetized key matches more
than one word.
More than one marked character is possible.
The code creates keys by iterating through all the possible combinations.
That is,
more than one possible unjumbled word can cause the program to deal with
more than one key for the final solution.

Multi-word solutions also cause complications.
The code deals with it by creating all possible combinations of the
characters for the first word, then recursively treating the leftover characters
as the other word(s) of the solution, making combinations for them.
A single set of marked characters from the unjumbled words
can end up as many partitioned-into-words keys for the solution.
Each partitioned-into-words combination has to match a dictionary
entry (remember them from above?) so may of the paritioned-into-words
keys can be discarded.
Each of those dictionary answers can be a group of words that alphabetize
to the key.
