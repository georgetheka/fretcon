# fretcon

Generates a fretted instrument fretboard in text format and draws custom data on it. 

## Using the Library

Instantiating a predefined instrument fretboard:
```
f := NewGuitar() 
f := NewShortGuitar() //half-length fretboard (12 frets)
f := NewBass4() //four-string bass 
f := NewBass5() //five-string bass
```

Instantiating a custom fretboard with a custom tuning:
```
startingFret := 0
endingFret := 22
stringList := []string{"d", "b", "a", "e"}
f, err := NewFretboard(startingFret, endingFret, startingList...)
```

Displaying data on the fretboard:
```
result, err := f.Draw(
    "1", "0", "0",
    "2", "1", "1",
    "3", "0", "0",
    "4", "2", "2",
    "5", "3", "3",
    "6", "0", "x",
)
```
result:

```
    0    1   2   3   4   5   6   7   8   9   10  11  12
e' _0_||___|___|___|___|___|___|___|___|___|___|___|___/
b  ___||_1_|___|___|___|___|___|___|___|___|___|___|___\
g  _0_||___|___|___|___|___|___|___|___|___|___|___|___/
d  ___||___|_2_|___|___|___|___|___|___|___|___|___|___\
a  ___||___|___|_3_|___|___|___|___|___|___|___|___|___/
e  _x_||___|___|___|___|___|___|___|___|___|___|___|___\
                 o       o       o       o           o
```

## Using the CLI Wrapper

The CLI wrapper simply exposes the GO library as a command line utility which could be useful when testing or experimenting.

Generating predefined fretboards:
```
./run.sh guitar 
./run.sh shortguitar 
./run.sh bass //same as bass4 
./run.sh bass4 
./run.sh bass5 
```

Generating fretboards of custom length, string number, and tuning: 
```
./run.sh "3 10 e a g d a d"
```
where the first two parameters within the quoted first argument specify the starting and ending fret numbers,
and the remaining sub-arguments specify the list of strings in the form of text labels.

Adding data to any fretboard - predefined or custom - is done by appending one or more argument triplets.
Each triplet represents the (string-number, fret-number, text-to-display):
```
./run.sh "0 12 e' b g d a e" \
1 0 0 \
2 1 1 \
3 0 0 \
4 2 2 \
5 3 3 \
6 0 x <Enter>
```

```
output:

    0    1   2   3   4   5   6   7   8   9   10  11  12
e' _0_||___|___|___|___|___|___|___|___|___|___|___|___/
b  ___||_1_|___|___|___|___|___|___|___|___|___|___|___\
g  _0_||___|___|___|___|___|___|___|___|___|___|___|___/
d  ___||___|_2_|___|___|___|___|___|___|___|___|___|___\
a  ___||___|___|_3_|___|___|___|___|___|___|___|___|___/
e  _x_||___|___|___|___|___|___|___|___|___|___|___|___\
                 o       o       o       o           o
```
