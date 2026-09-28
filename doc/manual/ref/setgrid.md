# SetGrid



Set size of the grid cells. All objects are placed in the grid.



##  Child elements

(none)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`dx` (optional)
:   Gap between two grid cells (horizontal)




`dy` (optional)
:   Gap between two grid cells (vertical)




`height` (length, optional)
:   The height of a grid cell. Use either height or ny, but not both.




`nx` (number, optional)
:   Specify the number of grid cells in horizontal direction. Use either nx or width, not both.




`ny` (number, optional)
:   Set the number of grid cells in vertical direction. Give ny or height, but not both.




`rounding` (optional)
:   How the height of an object is turned into a number of grid rows. Without the attribute the setting is kept.



    `up`
    :    An object takes the next whole number of rows (default).



    `nearest`
    :    An object takes the nearest number of rows. What that gains or loses is carried into the next object placed on the row below it, so a run of objects on a fine grid keeps its exact height instead of growing by up to a row per object.




`width` (length, optional)
:   The width of a grid cell. Use either width or nx, not both.




## Example

```xml
<SetGrid width="4mm" height="14pt"/>
```





