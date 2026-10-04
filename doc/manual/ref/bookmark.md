# Bookmark



Create a bookmark for the PDF viewer (e.g. Adobe Reader). When the user clicks on a bookmark, the PDF viewer jumps to that place in the document.



##  Child elements

(none)

##  Parent elements

[Td](../td), [TextBlock](../textblock)


## Attributes



`level` (number)
:   1 is the top level, 2 is the next level, etc.




`open` (optional)
:   If yes, the child elements are shown. If no, the child elements are hidden.



    `yes`
    :    Show children.



    `no`
    :    Hide children.




`select` ([XPath expressions](/programming/xpath))
:   Title of the bookmark




## Example

```xml
<Bookmark level="1" select="$title" open="no" />
```

Create a bookmark on level 1 (top level) with the title stored in the variable `title`.








## See also

- Commands: [A](../a), [Mark](../mark)
- Manual: [PDF Options: Bookmarks](/manual/advanced/pdf-options#bookmarks)

