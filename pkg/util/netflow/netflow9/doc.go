/*
Package netflow9 implements NetFlow version 9 as specified in RFC 3954.

About

NetFlow version 9 is the IETF standard mechanism for information export.

Structure

The basic output of NetFlow is a flow record. Several different formats for
flow records have evolved as NetFlow has matured. The most recent evolution of
the NetFlow flow-record format is known as NetFlow version 9. The
distinguishing feature of the NetFlow Version 9 format, which is the basis for
an IETF standard, is that it is template-based.

Templates provide an extensible design to the record format, a feature that
should allow future enhancements to NetFlow services without requiring
concurrent changes to the basic flow-record format.
*/
package netflow9
