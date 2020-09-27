# Jumbled-word puzzler solver

A program to help solve [Jumbled word puzzles](https://www.google.com/search?q=jumble&tbm=isch)

The user interface is pathetic.
I am not a UI/UX specialist,
and I don't want to put the effort into learning how to do
decent user interfaces.

## Building

Has no dependencies outside of the Go standard library.

    $ go build runserver.go

## Executing

It's a webapp that runs correctly probably only under Linux,
maybe under a BSD.
You DO NOT want to expose this to the internet.
I made no effort to keep it safe, or to have the web app part
validate user input.
It is a security risk of unknown magnitude.

```sh
$ ./runserver -d /usr/share/dict/words
```

The file specified with "-d" is a text file,
one word per line.
You'll have to access it with a URL like `http://localhost:8012/jumble`


## Design and Algorithm

I used HTML and JavaScript to build the user interface because
I don't want to invest the time to figure out how to do a native UI.
This damages the UI, but I sincerely doubt anyone will care, as nobody will try it.

It works by sorting the letters of every word in a file (`/usr/share/dict/words`, for example)
that has one word per line.
The words "meat", "team" and "meta" all end up alphabetized as "aemt".
It makes a hashtable/dictionary/associative array using an alphabetized word
as a key, and an array of words that alphabetize to that key as the value.

The program can alphabetize the characters in a jumbled input word,
and use the resulting string to look up all words that potentially match.
A jumbled word can match more than one unjumbled alternate word.
That lets it "unjumble" words.

The same thing happens in the final solution.
Marked characters of the un-jumbled words get used to create a key
to look up possible matches in the dictionary.

One complication happens when a letter-alphabetized key matches more than one word.
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

### Example

1. Fire up the web app: `./runserver -d /usr/share/dict/words`
2. Use a web browser to access the HTML generated by `runserver`: http://localhost:8012/jumble
3. In the web browser, select the number of words from the pull-down list
4. Set the number of letters in each individual word
5. Check the letters of each solved word that will be carried forward
6. Enter the jumbled letters into the individual fields of each word
7. Click "Unjumble Words". The web app will find all the real words whose
letters can be moved around to match the jumbled words you enter.
8. Sometimes, the letters in the jumbled words don't match any words
in the dictionary you used. 
You can unjumble the word youself, and click the "use as-is" checkbox.
I found the Jumble authors using scrambled letters for "monkees",
the mid-60s synthetic rock band, in two separate puzzles.
9. Select the number of words in the solution,
and the size(s) of the words in the solution.
This is a crap interface:
you enter a number for the number of words,
and space-separated numbers for the number of characters in each of the words
of the solution.
10. Click the Solve button.
The webapp will use the letters in the unjumbled words to create
combinations of sub-words as specified in the space-separated numbers.
Some times more than 1 word in the dictionary will match the jumbled characters,
so some puzzles will have several lists of letters carried forward.
The webapp will create combinations of each of the several lists of letters.
The webapp will try to find words in the dictionary file that have the same
letters as the combinations of letters from the lists of letters carried forward.

Oddly there can be only one solution for a given puzzle,
or hundreds can exist.
