# Columns



Set the widths and other properties of the columns in a table.



##  Child elements

[Column](../column), [ForAll](../forall), [Loop](../loop), [Switch](../switch), [Value](../value)

##  Parent elements

[Function](../function), [Table](../table)


## Attributes
(none)

## Remarks
The `\*` widths in the command “Column” allow dynamic cell widths. For that the total width of the table must be set and the attribute (on [Table](../table)) `stretch` must be set to `yes`.
        The widths of the columns are calculated as follows: first the absolute widths are taken into account. After that, the `*` columns are distributed across the remaining space. The
        numbers before the `*` denote the fraction of the space. In the example below the third column gets 1/6 of the remaining width, the fourth column get 5/6.


## Example

```xml
<Table>
  <Columns>
    <Column width="14mm" />
    <Column width="2" />
    <Column width="1*" align="right" valign="top" />
    <Column width="5*" />
    <Column width="5mm" backgroundcolor="gray" />
    <Tr>
       ....
    </Tr>
  </Columns>
</Table>

```





