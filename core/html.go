package core

// func (xd *xtsDocument) parseHTML(elt *goxml.Element) (frontend.FormatToVList, error) {
// 	ftv := func(wd bag.ScaledPoint) (*node.VList, error) {
// 		str := elt.ToXML()
// 		d := document.NewWithFrontend(xd.document, xd.layoutcss)
// 		te, err := d.HTMLToText(str)
// 		if err != nil {
// 			return nil, err
// 		}
// 		vl, err := d.CreateVlist(te, wd)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return vl, err
// 	}
// 	return ftv, nil
// }
