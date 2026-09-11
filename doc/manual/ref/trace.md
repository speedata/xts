# Trace



Set debugging switches



##  Child elements

(none)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`boxmodel` (optional)
:   Draw a translucent overlay visualizing the CSS box model of block-level elements (margin, border, padding and content area, colored like the browser developer tools). Where margins of neighboring elements overlap, the colors add up and appear darker.



    `yes`
    :    Show the box model overlay.



    `no`
    :    Don't show the box model overlay (default).




`dests` (optional)
:   Draw PDF destinations with a black circle.



    `yes`
    :    Show dest.



    `no`
    :    Don't show dests (default).




`grid` (optional)
:   Draw the grid on the page.



    `yes`
    :    Show the grid.



    `no`
    :    Don't show the grid (default).




`gridallocation` (optional)
:   Draw allocated cells with yellow and conflicts with red markers.



    `yes`
    :    Show grid allocation.



    `no`
    :    Don't show the grid allocation (default).




`hboxes` (optional)
:   Visualize the vertical metrics of text lines: the area above the baseline (height) and the area below it (depth) get distinct translucent tints, the baseline and the outline of each line box are stroked.



    `yes`
    :    Show the line boxes.



    `no`
    :    Don't show the line boxes (default).




`hyperlinks` (optional)
:   Draw a box around hyperlinks.



    `yes`
    :    Show the hyperlinks.



    `no`
    :    Don't show the hyperlinks (default).




`hyphenation` (optional)
:   Draw little marks to show all hyphenation points.



    `yes`
    :    Show hyphenation points.



    `no`
    :    Don't show any hyphenation points (default).




`objects` (optional)
:   Draw boxes around objects.



    `yes`
    :    Show boxes.



    `no`
    :    Don't show the bounding boxes of objects.




## Example

```xml
<Trace grid="yes" />
```





