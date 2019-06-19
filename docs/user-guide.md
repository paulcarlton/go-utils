(c) Copyright 2018 Hewlett Packard Enterprise Development LP

# User Guide

This repository contains testing and development utilities as well as library code for use with golang and python.

## Golang Utilities

The 'goutils' directory contains golang utility functions.

### CoreError

This defines a error which comprises an error code and basic error information as well as options to provide
further details and recommended actions to address the issue. It also includes a nested error field so the
underlying error details can be included.  This enables the reporting of a stack of errors.  Finally, there
is an 'ID' field which can be be used to record the user defined subject of the error or processing routine
that emitted the error.

To report an application error use `MakeError()` to generate a new 'CoreError' without a nested error.
If a called function returns an error, use `RaiseError` to generate a new CoreError with the error
from the called function as the nested error. If the nested error is not a CoreError, RaiseError will
create a new CoreError using the supplied error.

When reporting a CoreError use the `CoreError.FullInfo()` method to report all details and nested errors or
`CoreError.String()` to just report the 'id', error code and message.

CoreError implements the 'Error()' method so it can be used anywhere an 'error' is expected.

An `ErrorText()` function is also provided which will will generate text from whatever it is passed.
If it is passed a CoreError it will return the output from CoreError.FullInfo().  If the interface past
to it is not a CoreError but implements 'Error()', it will return the output from that method.  Otherwise
it will create a new 'error' from the input and return the output from the Error() method of that new error.
This can be used when generating log output to report an error without having to be concerned with its type.

### Callers and GetCaller

The `GetCaller()` function provides details of the Caller's Function Name and source file/line number. This
can be used to include the function name in log output. The `Callers()` function will return details of the
caller's caller and their caller as far up the stack as is requested and available. This can be used in
debugging output.

