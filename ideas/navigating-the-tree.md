
# Navigating the OPML Tree

2016-04-06, R. S. Doiel <rsdoiel@gmail.com>

It seem convient to have an API method on the OPML structure that supports
selecting a specific node through a path of index values which returns
the a specific outline element. This element could then have simple operations such as Insert (before this node), Append (after this node), Update (this node), Delete (this node). These could be used to create
more complex functions like AppendChild, InsertChild, UpdateChild, DeleteChild . 

SelectElement would take an address and work from the root of the OPML file.The address would be an array of integeters indicating the traversal of the tree structure. An empty array would would indicate the root node. A value of 1 would indicate the first sibling. A two cell array would indicate a nested list such that a[0] indicates the siblings at the root and a[1] would indicate the immediate children. Grandchildren would be a three element array, great grand children four elements and so on. Thus the array of integers would indicate the path to the targeted element.

A variation of this could be used to select groups of elements from the tree by allowing ranges and lists of individual nodes. A friendly string representation could thus be expressed as path. "2/3" would be come [2,3] which would mean the third sibling from the root (arrays count from zero) and the fourth child of that sibling. Groupings of individual noted in a branch would be separated by commas and ranges by a dash. "1,2/3-4" would give us the fourth and fifth children of the roots' second and third siblings.




