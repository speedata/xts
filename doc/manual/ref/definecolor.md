# DefineColor



Colors defined with DefineColors can be referenced later by their name.



##  Child elements

(none)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`b` (0 to 100 or 0 to 255, optional)
:   Blue part with rgb (0-100) or RGB (0-255).




`c` (0 up to 100, optional)
:   Cyan part with cmyk (0-100).




`colorname` (text, optional)
:   The name of the spot color if model is “spotcolor”. The name must match the required color name, such as “PANTONE 116 C”.




`g` (0 to 100 or 0 to 255, optional)
:   Green part with rgb (0-100) or RGB (0-255) / gray part when using the model gray (0-100).




`k` (0 up to 100, optional)
:   Black part with cmyk (0-100).




`m` (0 up to 100, optional)
:   Magenta part with cmyk (0-100)




`model` (optional)
:   Color model to be used for the color. Currently “rgb”, “cmyk”, “gray” and “spotcolor” are supported.



    `cmyk`
    :    CMYK (cyan, magenta, yellow, key/black), values between 0 and 100 (100 = full intensity)



    `rgb`
    :    rgb (red, green, blue), values between 0 and 100, 100 means full intensity



    `RGB`
    :    rgb (red, green, blue), values between 0 and 255, 255 means full intensity



    `gray`
    :    Gray (0=black, 100=white)



    `spotcolor`
    :    Use a PANTONE or HKS color.




`name` (text)
:   The name of the color to be defined.




`r` (0 to 100 or 0 to 255, optional)
:   Red part with rgb (0-100) or RGB (0-255).




`value` (text, optional)
:   Hex value of the color, such as `#FA5` or `#FFAA55` or `rgb(255,170,85)` or `rgba(255,170,85,1)`.




`y` (0 up to 100, optional)
:   Yellow part with cmyk (0-100).




## Example

```xml
<DefineColor name="black" model="cmyk" c="0" m="0" y="0" k="100"/>
<DefineColor name="white" model="rgb" r="100" g="100" b="100"/>
```





## Info

The CSS level 3 colors are predefined in RGB-space. See http://www.w3.org/TR/css3-color/ for the definitions. That means you can use common colors such as `red` or `goldenrod` without using [DefineColor](../definecolor).



The predefined colors are: `aliceblue`, `black`, `orange`, `rebeccapurple`, `antiquewhite`, `aqua`, `aquamarine`, `azure`, `beige`, `bisque`, `blanchedalmond`, `blue`, `blueviolet`, `brown`, `burlywood`, `cadetblue`, `chartreuse`, `chocolate`, `coral`, `cornflowerblue`, `cornsilk`, `crimson`, `darkblue`, `darkcyan`, `darkgoldenrod`, `darkgray`, `darkgreen`, `darkgrey`, `darkkhaki`, `darkmagenta`, `darkolivegreen`, `darkorange`, `darkorchid`, `darkred`, `darksalmon`, `darkseagreen`, `darkslateblue`, `darkslategray`, `darkslategrey`, `darkturquoise`, `darkviolet`, `deeppink`, `deepskyblue`, `dimgray`, `dimgrey`, `dodgerblue`, `firebrick`, `floralwhite`, `forestgreen`, `fuchsia`, `gainsboro`, `ghostwhite`, `gold`, `goldenrod`, `gray`, `green`, `greenyellow`, `grey`, `honeydew`, `hotpink`, `indianred`, `indigo`, `ivory`, `khaki`, `lavender`, `lavenderblush`, `lawngreen`, `lemonchiffon`, `lightblue`, `lightcoral`, `lightcyan`, `lightgoldenrodyellow`, `lightgray`, `lightgreen`, `lightgrey`, `lightpink`, `lightsalmon`, `lightseagreen`, `lightskyblue`, `lightslategray`, `lightslategrey`, `lightsteelblue`, `lightyellow`, `lime`, `limegreen`, `linen`, `maroon`, `mediumaquamarine`, `mediumblue`, `mediumorchid`, `mediumpurple`, `mediumseagreen`, `mediumslateblue`, `mediumspringgreen`, `mediumturquoise`, `mediumvioletred`, `midnightblue`, `mintcream`, `mistyrose`, `moccasin`, `navajowhite`, `navy`, `oldlace`, `olive`, `olivedrab`, `orangered`, `orchid`, `palegoldenrod`, `palegreen`, `paleturquoise`, `palevioletred`, `papayawhip`, `peachpuff`, `peru`, `pink`, `plum`, `powderblue`, `purple`, `red`, `rosybrown`, `royalblue`, `saddlebrown`, `salmon`, `sandybrown`, `seagreen`, `seashell`, `sienna`, `silver`, `skyblue`, `slateblue`, `slategray`, `slategrey`, `snow`, `springgreen`, `steelblue`, `tan`, `teal`, `thistle`, `tomato`, `turquoise`, `violet`, `wheat`, `white`, `whitesmoke`, `yellow` and `yellowgreen`




