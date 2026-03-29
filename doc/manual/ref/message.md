# Message



Writes a message onto the console and to the protocol file.



##  Child elements

[Element](../element), [Value](../value)

##  Parent elements

[AtPageCreation](../atpagecreation), [AtPageShipout](../atpageshipout), [Case](../case), [Contents](../contents), [ForAll](../forall), [Function](../function), [Layout](../layout), [Loop](../loop), [Otherwise](../otherwise), [Record](../record), [Section](../section), [Tr](../tr), [Until](../until), [While](../while)


## Attributes



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






