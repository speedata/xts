# Image



Includes an external Graphic. Allowed graphic formats are PDF (.pdf), PNG (.png) and JPEG (.jpg). See below for a limitation on the number of included PDF files.



##  Child elements

(none)

##  Parent elements

[PlaceObject](../placeobject), [Td](../td)


## Attributes



`height` (number or length, optional)
:   Image height. One of 'auto' (default, take image width), length (such as '3cm') or number (in grid cells).




`href` (text, optional)
:   Filename of the image. Can be a file in the search path, an absolute file name, a file-URI for absolute paths (e.g. `file:///path/to/image.pdf`) or a location on the web (http, https). Use `placeholder://WxH` (e.g. `placeholder://200x150`) to generate a placeholder image with the given dimensions in PDF points.




`maxheight` (number or length, optional)
:   The maximum height of the image. Only used when clip="no". Value is a number (grid cells) or a length.




`maxwidth` (number or length, optional)
:   The maximum width of the image. Only used when clip="no". Value is a number (grid cells), a length or the value “100%” for full width image.




`minheight` (number or length, optional)
:   The minimum height of the image. Only used when clip="no". Value is a number (grid cells) or a length.




`minwidth` (number or length, optional)
:   The minimum width of the image. Only used when clip="no". Value is a number (grid cells), a length or the value “100%” for full width image.




`page` (number, optional)
:   The page number from the PDF. Default is 1 (include the first page).




`stretch` (yes or no, optional)
:   Stretch image until one of maximum width and maximum height is reached. Useful if images should be as large as possible but should not use more than the given space.




`visiblebox` (optional)
:   The PDF box that represents the visible area of the included image. Default is “cropbox”.



    `artbox`
    :    Use the artbox as the visible area. The artbox is usually not contained in a PDF.



    `bleedbox`
    :    Use the bleedbox of the included PDF.



    `cropbox`
    :    Use the cropbox of the included PDF (default).



    `mediabox`
    :    Use the mediabox of the included PDF. This is the largest box.



    `trimbox`
    :    Use the trimbox of the includes PDF. The trimbox is the final paper size. For example, the trim box of an A4 PDF is 210mm x 297mm.




`width` (number or length, optional)
:   Image width. One of 'auto' (default, take image width), '100%' (whole area width), length (such as '3cm') or number (in grid cells).




## Example

```xml
<Record match="productdata">
  <PlaceObject column="{ $column }">
    <Image width="10" file="{ string(.) }"/>
  </PlaceObject>
</Record>

```

Takes the file name of the image from the contents of the current element in the data file (here: productdata). Sample data XML:


```xml
<productdata>image.pdf</productdata>
```





