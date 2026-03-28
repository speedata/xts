# CopyOf



Copies the result of an XPath expression as-is, preserving node structure. Analogous to xsl:copy-of in XSLT.



##  Child elements

(none)

##  Parent elements

[A](../a), [B](../b), [I](../i), [Li](../li), [Paragraph](../paragraph), [SetVariable](../setvariable), [Span](../span), [U](../u)


## Attributes



`select` ([XPath expressions](/manual/data-processing/xpath))
:   XPath expression whose result is passed through unchanged.




## Remarks
Unlike Value, which always converts to text, CopyOf preserves the original node type (HTML nodes, XML elements, etc.).

Use this command to pass structured content from sd:decode-html() or XML data into Paragraph, HTML, or other elements.


## Example

```xml
<Paragraph>
  <CopyOf select="sd:decode-html(.)" />
</Paragraph>
```





