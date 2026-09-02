# Column



Set the properties of a column in the table.



##  Child elements

(none)

##  Parent elements

[Case](../case), [Columns](../columns), [ForAll](../forall), [Function](../function), [Loop](../loop), [Otherwise](../otherwise), [SetVariable](../setvariable), [Until](../until), [While](../while)


## Attributes



`align` (optional)
:   Horizontal alignment of the contents of the cells in this column. A cell can override the alignment with its own align attribute or with CSS.



    `left`
    :    The contents is aligned at the left.



    `center`
    :    The contents is centered.



    `right`
    :    The contents is aligned at the right.



    `justify`
    :    The contents is justified.




`background-color` (text, optional, "CSS property": background-color)
:   Background color of the cells in this column. A cell with its own background color takes precedence.




`valign` (optional)
:   Vertical alignment of the contents of the cells in this column. A cell can override the alignment with CSS (`vertical-align`).



    `top`
    :    The contents is aligned at the top.



    `middle`
    :    The contents is centered vertically.



    `bottom`
    :    The contents is aligned at the bottom.




`width` (Number, length or *-numbers, optional)
:   Width of the column. Argument can be a number (in grid cells) a length (e.g. 2cm) or a \*-number (e.g. 4*).




## Example


See the example at [Columns](../columns).







