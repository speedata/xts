# Feature parity with the speedata Publisher

XTS is the successor of the speedata Publisher, but it is not a port.
This list records which Publisher commands and layout functions exist in
XTS and which do not, so that the gaps are known when a milestone is
planned. It replaces the former checklist issues #21 and #15. A checked
item exists in XTS, possibly under a different name or with different
attributes; an unchecked item is a candidate, not a promise.

Update this file when a command or function is added. Whether a gap is
worth closing is decided per milestone, not here.

## Commands

- [x] AddSearchpath
- [x] AtPageShipout
- [x] AttachFile
- [ ] Clip
- [ ] DefineColorprofile
- [ ] DefineGraphic
- [ ] DefineMatter
- [ ] Frame
- [ ] Transformation
- [x] Groupcontents
- [x] HTML
- [ ] HSpace
- [ ] Hyphenation
- [ ] Include
- [ ] Initial
- [ ] InsertPages
- [ ] Makeindex
- [ ] NoBreak
- [ ] Output
- [ ] Overlay
- [ ] Position
- [ ] Rule
- [ ] SavePages
- [ ] SortSequence
- [ ] TableNewPage
- [x] Tablefoot (TableFoot)
- [x] Tablehead (TableHead)
- [x] Tablerule
- [ ] Text
- [ ] URL
- [ ] VSpace
- [x] Sub
- [x] DefineTextformat
- [x] Trace
- [x] Fontface
- [x] Color
- [x] Sup

## Layout functions

- [ ] sd:allocated()
- [ ] sd:alternating()
- [ ] sd:count-saved-pages()
- [ ] sd:current-column()
- [ ] sd:current-framenumber()
- [x] sd:dimexpr()
- [ ] sd:firstmark()
- [ ] sd:first-free-row()
- [ ] sd:format-string()
- [ ] sd:keep-alternating()
- [ ] sd:lastmark()
- [ ] sd:merge-pagenumbers()
- [ ] sd:pageheight()
- [ ] sd:pagewidth()
- [ ] sd:randomitem()
- [ ] sd:reset-alternating()
- [ ] sd:variable-exists()
- [ ] sd:visible-pagenumber()
- [x] sd:filecontents()
- [x] sd:decode-base64()
- [x] sd:decode-html()
- [x] sd:variable()
- [x] sd:attr()
- [x] sd:format-number()
- [x] sd:sha1()
- [x] sd:md5()
- [x] sd:mode()
- [x] sd:sha256()
- [x] sd:tounit()
- [x] sd:sha512()
