# scte224

Utilities for modeling SCTE 224 data in Go

Details on SCTE 224 can be found at: https://www.scte.org/standards-development/library/standards-catalog/ansiscte-224-2018r1/

2015 schema: http://www.scte.org/schemas/224/SCTE224-20151115.xsd (link no longer working, refer to the spec located alongside structs in this repo)

2018 schema: http://www.scte.org/schemas/224/SCTE224-20180501.xsd (link no longer working, refer to the spec located alongside structs in this repo)

Latest schema (currently points to 2020): https://www.scte.org/standards-development/library/standards-catalog/ansiscte-224-2018r1/

**Some general notes on the representation of SCTE224 with Structs**
* For sub elements, use pointers to unmarshal the values into. This helps two things-
  * Avoids zero value sub elements, during marshaling 
  * Efficient when you have to pass these large structs around, since using pointers avoids copying over and over in memory
