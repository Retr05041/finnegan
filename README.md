# finnegan
Finnegan is a Fill-Ins puzzle solver

### Ideas..
Best possible solution: Place a word horizontally on the first empty cell, then check every cell the word is on vertically for matches (accounting for placement of the cell... i.e. if it's the 3rd rune in a vertical candidate)
once all verticals have been checked for that horizontal, move to the next empty cell
