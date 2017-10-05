set oEnv = CreateObject("WScript.Shell").Environment("Process")
chars = array (~ENVVARS~)
getCombinations "", chars 

Function Slice(arr, starting, ending)
    Dim out_array
    If Right(TypeName(arr), 2) = "()" Then
        out_array = Array()
        ReDim Preserve out_array(ending - starting)
        For index = starting To ending
            out_array(index - starting) = arr(index)
        Next
    Else
        Exit Function
    End If
    Slice = out_array
End Function
Function getCombinations(prefix, chars) 
    For index = 0 To ubound(chars)
        envCombos.add prefix + chars(index)
        getCombinations (prefix + chars(index)), Slice(chars, index + 1, ubound(chars))
    Next
End Function