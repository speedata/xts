# PDFOptions



Set PDF options



##  Child elements

(none)

##  Parent elements

[Layout](../layout), [Section](../section)


## Attributes



`author` (text, optional)
:   Set the author of the document




`creator` (text, optional)
:   Set the creator application of the document




`displaymode` (optional)
:   Select the display mode when opening PDF document (mainly with Acrobat).



    `attachments`
    :    Display the attachment pane.



    `bookmarks`
    :    Display the bookmarks pane (only works if the PDF document contains at least one bookmark).



    `fullscreen`
    :    Open the document in fullscreen mode.



    `none`
    :    Do not display a special pane.



    `thumbnails`
    :    Display the thumbnail pane.




`duplex` (optional)
:   Set viewer preference to one or two page printing. Default: empty.



    `simplex`
    :    One page per sheet



    `duplexflipshortedge`
    :    Two pages per sheet and flip on short edge



    `duplexfliplongedge`
    :    Two pages per sheet and flip on long edge




`format` (optional)
:   Set the PDF output format.



    `PDF/A-3b`
    :    PDF/A-3b format (supports file attachments, required for ZUGFeRD)



    `PDF/X-3`
    :    PDF/X-3 format (for printing)



    `PDF/X-4`
    :    PDF/X-4 format (for printing)



    `PDF/UA`
    :    PDF/UA format (for accessibility)




`picktraybypdfsize` (optional)
:   Activate the check box in the PDF viewer for choosing the paper tray based on the page size.



    `yes`
    :    Activate checkbox



    `no`
    :    Deactivate checkbox




`printscaling` (optional)
:   Should the printer scale the pages?



    `appdefault`
    :    Use the default from the PDF viewer



    `none`
    :    No page scaling




`showhyperlinks` (yes or no, optional)
:   Show hyperlinks in Adobe Acrobat and perhaps other PDF viewers.




`subject` (text, optional)
:   Set the subject of the document




`title` (text, optional)
:   Set the title of the document




## Example

```xml
<PDFOptions author="A. U. Thor" subject="An interesting story" />
```





