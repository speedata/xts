# Options



Set publisher specific options.



##  Child elements

(none)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`bleed` (length, optional)
:   The amount of bleed. Defaults to 0mm.




`cutmarks` (optional)
:   Cut marks / crop marks will be placed in the PDF. The distance of the marks from the imaginary center is determined by the attribute `trim`, but is at least 5mm. The length of the cut marks is 1cm. The default of this attribute is `no`, that means no cut marks will be displayed.



    `yes`
    :    Show crop marks.



    `no`
    :    Don't show crop marks (default).




`features` (text, optional)
:   A comma separated list of OpenType features such as +kern,-liga. Used as a default for all paragraphs.




`imagenotfound` (optional)
:   Controls behavior when an image file is not found. With "warning" (default), a placeholder image is inserted and processing continues. With "error", a missing image stops the layout run.



    `warning`
    :    Insert a placeholder image and continue (default).



    `error`
    :    Report an error and stop processing.




`mainlanguage` (optional)
:   The default language for text (hyphenation and rendering). You can also set the default language on the command line and locally set it at [Paragraph](../paragraph) and [TextBlock](../textblock).




`missingglyph` (optional)
:   Controls behavior when a character is not found in the font. With "warning" (default), a warning is logged and the character is skipped. With "error", the layout run is stopped. With "none", missing glyphs are silently ignored.



    `warning`
    :    Log a warning and skip the character (default).



    `error`
    :    Report an error and stop processing.



    `none`
    :    Silently ignore missing glyphs.




## Example

```xml
<Options
    cutmarks="yes"
    bleed="3mm"/>
```





## Info

The list of languages and the short code known to the system are:



`Ancient Greek` (`grc`), `Armenian` (`hy`), `Bahasa Indonesia` (`id`), `Basque` (`eu`), `Bulgarian` (`bg`), `Catalan` (`ca`), `Chinese` (`zh`), `Croatian` (`hr`), `Czech` (`cs`), `Danish` (`da`), `Dutch` (`nl`), `English` (`en_GB`), `English (Great Britain)` (`en_GB`), `English (USA)` (`en_US`), `Esperanto` (`eo`), `Estonian` (`et`), `Finnish` (`fi`), `French` (`fr`), `Galician` (`gl`), `German` (`de`), `Greek` (`el`), `Gujarati` (`gu`), `Hindi` (`hi`), `Hungarian` (`hu`), `Icelandic` (`is`), `Irish` (`ga`), `Italian` (`it`), `Kannada` (`kn`), `Kurmanji` (`ku`), `Latvian` (`lv`), `Lithuanian` (`lt`), `Malayalam` (`ml`), `Norwegian Bokmål` (`nb`), `Norwegian Nynorsk` (`nn`), `Other` (`--`), `Polish` (`pl`), `Portuguese` (`pt`), `Romanian` (`ro`), `Russian` (`ru`), `Serbian` (`sr`), `Serbian (cyrillic)` (`sc`), `Slovak` (`sk`), `Slovenian` (`sl`), `Spanish` (`es`), `Swedish` (`sv`), `Turkish` (`tr`), `Ukrainian` (`uk`), `Welsh` (`cy`)




