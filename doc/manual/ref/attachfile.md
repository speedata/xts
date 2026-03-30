# AttachFile



Attach a file to the PDF document. This is used for example to embed ZUGFeRD/Factur-X invoices.



##  Child elements

(none)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Contents](../contents), [Layout](../layout), [Record](../record), [Section](../section)


## Attributes



`description` (text, optional)
:   Description of the attachment.




`href` (text)
:   Path to the file to attach.




`name` (text, optional)
:   Display name in PDF viewer. Defaults to the file name.




`type` (text, optional)
:   The MIME type of the attachment (e.g. "application/pdf"). If set to "ZUGFeRD invoice" or "facturx", the attached XML is treated as a Factur-X / ZUGFeRD invoice: the conformance profile is detected automatically from the XML data, the required XMP metadata for ZUGFeRD compliance is added, and the PDF format is automatically set to PDF/A-3b. If not set, the MIME type is detected from the file extension.




## Example

```xml
<AttachFile href="zugferd.xml" name="factur-x.xml"
    description="ZUGFeRD invoice data" type="ZUGFeRD invoice" />
```





