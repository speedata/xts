# Message



Writes a message onto the console and to the protocol file.



##  Child elements

[Element](../element), [Value](../value)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Layout](../layout), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Section](../section), [Tr](../tr), [Until](../until), [While](../while)


## Attributes



`errorcode` (number, optional)
:   If an error is raised, use this code on exit. Defaults to 1. Negative values are reserved for system purpose.




`exit` (optional)
:   Tells the software to exit immediately.



    `no`
    :    The speedata Publisher continues with the PDF creation.



    `yes`
    :    The speedata Publisher exits without finishing the PDF file.




`select` ([XPath expressions](/manual/data-processing/xpath), optional)
:   Contents of the message. You can alternatively specify the message by the child elements [Value](../value).




`type` (optional)
:   The type of the message.



    `error`
    :    Report an error.



    `warning`
    :    Report a warning.



    `debug`
    :    Report a debug message.



    `notice`
    :    Report an information (default).




## Example






